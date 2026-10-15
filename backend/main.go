
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type ServiceItem struct {
	ID          int       `json:"id"`
	ServiceName string    `json:"service_name"`
	OwnerTeam   string    `json:"owner_team"`
	StorageType string    `json:"storage_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type HealthResponse struct {
	Status           string `json:"status"`
	StudentNIM       string `json:"student_nim"`
	StudentName      string `json:"student_name"`
	ContainerOrHostID string `json:"container_or_host_id"`
	StorageMode      string `json:"storage_mode"`
}

var (
	memStore []ServiceItem
	memMutex sync.Mutex
	db       *sql.DB
)

var (
	studentNIM  = "42430052"
	studentName = "Gede Troy Wiswama Bhareswara"
)

func main() {
	storageMode := getEnv("STORAGE_MODE", "memory")

	if storageMode == "postgres" {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getEnv("DB_HOST", "db"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASSWORD", "postgres"),
			getEnv("DB_NAME", "praktikum"),
		)

		var err error
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal(err)
		}

		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}
		log.Println("Terhubung ke PostgreSQL")
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status:            "ok",
			StudentNIM:        studentNIM,
			StudentName:       studentName,
			ContainerOrHostID: hostname,
			StorageMode:       storageMode,
		})
	})

	http.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			if storageMode == "postgres" {
				rows, err := db.Query(
					"SELECT id, service_name, owner_team, created_at FROM services ORDER BY id",
				)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer rows.Close()

				var services []ServiceItem
				for rows.Next() {
					var s ServiceItem
					if err := rows.Scan(
						&s.ID, &s.ServiceName, &s.OwnerTeam, &s.CreatedAt,
					); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					s.StorageType = "postgresql-persistent"
					services = append(services, s)
				}
				json.NewEncoder(w).Encode(services)
				return
			}

			memMutex.Lock()
			result := append([]ServiceItem(nil), memStore...)
			memMutex.Unlock()
			json.NewEncoder(w).Encode(result)

		case http.MethodPost:
			var input struct {
				ServiceName string `json:"service_name"`
				OwnerTeam   string `json:"owner_team"`
			}

			if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
				http.Error(w, "JSON tidak valid", http.StatusBadRequest)
				return
			}

			if input.ServiceName == "" || input.OwnerTeam == "" {
				http.Error(w, "service_name dan owner_team wajib diisi", http.StatusBadRequest)
				return
			}

			if storageMode == "postgres" {
				var s ServiceItem
				err := db.QueryRow(
					"INSERT INTO services (service_name, owner_team) VALUES ($1, $2) RETURNING id, service_name, owner_team, created_at",
					input.ServiceName, input.OwnerTeam,
				).Scan(&s.ID, &s.ServiceName, &s.OwnerTeam, &s.CreatedAt)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				s.StorageType = "postgresql-persistent"
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(s)
				return
			}

			memMutex.Lock()
			newItem := ServiceItem{
				ID:          len(memStore) + 1,
				ServiceName: input.ServiceName,
				OwnerTeam:   input.OwnerTeam,
				StorageType: "in-memory",
				CreatedAt:   time.Now(),
			}
			memStore = append(memStore, newItem)
			memMutex.Unlock()

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newItem)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server berjalan di port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
