# Laporan Praktikum Modul 1 (Pertemuan 4): Kontainerisasi Docker & Orkestrasi Multi-Tier (Lab 1 & Lab 2)

-   **Nama Mahasiswa:** Gede Troy Wiswama Bhareswara
-   **NIM:** 42430052
-   **Kelas / Prodi:** S1 Teknologi Informasi --- Universitas Pendidikan
    Nasional (UNDIKNAS) Denpasar
-   **Mata Kuliah:** DevOps (MK40)
-   **Branch GitHub:** `lab1-2/42430052-troy`
-   **Repository GitHub:**
    [praktikum-docker](https://github.com/troybhareswara/praktikum-docker.git)
-   **Nama Rekan Peer-Reviewer:** Belum dilakukan; akan dilengkapi
    setelah sesi peer-review bersama rekan.

------------------------------------------------------------------------

## 1. Tabel Komparasi Eksperimen Individu (Tahap 1 vs. Tahap 2 vs. Tahap 3)

  ---------------------------------------------------------------------------------------
  Parameter         Tahap 1: Tanpa Docker (Native di  Tahap 2: Docker   Tahap 3: Docker
  Pembanding        Host OS)                          Single-Stage      Multi-Stage
  ----------------- --------------------------------- ----------------- -----------------
  Kebutuhan Go di   Go SDK diperlukan untuk           Go SDK tidak      Go SDK tidak
  host              menjalankan aplikasi menggunakan  diperlukan di     diperlukan di
                    `go run main.go`.                 host untuk        host untuk
                                                      menjalankan       menjalankan image
                                                      container.        hasil build.

  Ukuran image      Tidak menghasilkan image Docker.  **1.32 GB**       **28.2 MB**
                                                      berdasarkan hasil berdasarkan hasil
                                                      pengukuran.       pengukuran.

  Output identitas  Endpoint `/healthz` menampilkan   Endpoint          Endpoint
                    informasi mahasiswa, mode         `/healthz`        `/healthz`
                    penyimpanan, dan hostname host.   menampilkan       menampilkan
                                                      informasi         informasi
                                                      aplikasi di dalam mahasiswa, mode
                                                      container.        penyimpanan, dan
                                                                        identitas
                                                                        container.

  User eksekusi     Aplikasi berjalan sebagai         Image             Aplikasi berjalan
                    pengguna host.                    single-stage      sebagai user
                                                      berjalan dengan   non-root
                                                      user default dari `appuser` dengan
                                                      image, yaitu      UID `10001`.
                                                      root.             

  Menjalankan       Dua proses pada port host yang    Beberapa          Image dapat
  beberapa instansi sama tidak dapat berjalan         container dapat   dijalankan
                    bersamaan; muncul error           diuji dengan      sebagai container
                    `bind: address already in use`.   pemetaan port     dan menggunakan
                                                      host yang         pemetaan port
                                                      berbeda, misalnya host sesuai
                                                      `8081` dan        kebutuhan.
                                                      `8082`.           

  Hasil utama       Memastikan aplikasi berjalan      Aplikasi berhasil Image jauh lebih
                    sebelum dikontainerisasi.         dikemas, tetapi   kecil dan
                                                      ukuran image      aplikasi berjalan
                                                      besar karena      sebagai user
                                                      lingkungan Go     non-root.
                                                      ikut dibawa.      
  ---------------------------------------------------------------------------------------

### Analisis komparasi

Pada Tahap 1, aplikasi dijalankan langsung di sistem operasi laptop
menggunakan Go. Endpoint `/healthz` berhasil menampilkan status aplikasi
dan identitas yang dikonfigurasi. Karena proses berjalan di host,
informasi hostname mengacu pada host. Saat dua proses mencoba
menggunakan port `8080` yang sama, proses kedua gagal karena port sudah
dipakai oleh proses pertama.

Pada Tahap 2, aplikasi dikemas menggunakan Dockerfile single-stage.
Image `backend-api:42430052-single` berhasil dibuat dengan ukuran
sekitar **1.32 GB**. Pengujian dua container menggunakan port host yang
berbeda menunjukkan bahwa beberapa instansi aplikasi dapat berjalan
bersamaan selama port host tidak bentrok. Ukuran image relatif besar
karena image masih membawa lingkungan Go yang lengkap.

Pada Tahap 3, proses kompilasi dipisahkan dari image runtime melalui
multi-stage build. Image `backend-api:42430052-multistage` berukuran
**28.2 MB**, turun sekitar **97.9%** dibandingkan image single-stage
berdasarkan ukuran yang teramati. Stage builder menangani pengunduhan
dependensi dan kompilasi, sedangkan stage produksi hanya menyalin binary
menggunakan `COPY --from=builder`. Pengaturan `CGO_ENABLED=0` dan
`-ldflags="-s -w"` membantu menghasilkan binary yang lebih ringkas.
Pengujian identitas proses menunjukkan bahwa aplikasi berjalan sebagai
`appuser` dengan UID `10001`, bukan root.

------------------------------------------------------------------------

## 2. Bukti Screenshot dan Analisis Setiap Tahapan

### A. Tahap 1 --- Eksekusi Aplikasi Tanpa Docker (Baseline)

![Bukti Screenshot Tahap
1](docs/screenshots/sc01-tahap1-tanpa-docker.png)

**Analisis Tahap 1**

Aplikasi dijalankan menggunakan `go run main.go`. Endpoint `/healthz`
berhasil memberikan respons yang memuat status aplikasi, NIM, nama
mahasiswa, hostname, dan mode penyimpanan. Karena aplikasi berjalan
langsung di laptop, hostname yang ditampilkan merupakan hostname host.
Saat aplikasi kedua dijalankan pada port `8080` yang sama, proses kedua
gagal karena port tersebut sudah digunakan.

### B. Tahap 2 --- Container dengan Single-Stage Dockerfile

![Bukti Screenshot Tahap
2](docs/screenshots/sc02-tahap2-single-stage.png)

**Analisis Tahap 2**

Aplikasi dikemas menggunakan `backend/Dockerfile.single` dengan base
image Go. Image `backend-api:42430052-single` berhasil dibuat dengan
ukuran **1.32 GB**. Dua container diuji dengan port host yang berbeda,
yaitu `8081` dan `8082`, sedangkan aplikasi di dalam container tetap
mendengarkan port `8080`. Dari percobaan ini terlihat bahwa container
memungkinkan aplikasi dijalankan terpisah dari lingkungan host, tetapi
ukuran image single-stage masih besar.

### C. Tahap 3 --- Optimasi Dockerfile dengan Multi-Stage Build

![Bukti Screenshot Tahap 3](docs/screenshots/sc03-tahap3-multistage.png)

**Analisis Tahap 3**

Image `backend-api:42430052-multistage` berhasil dibangun dengan ukuran
**28.2 MB**, sehingga memenuhi target di bawah 100 MB pada modul.
Dibandingkan image single-stage sebesar 1.32 GB, ukuran image turun
sekitar **97.9%**. Stage `builder` digunakan untuk mengunduh dependensi
dan mengompilasi aplikasi. Stage `production` hanya menerima binary
melalui `COPY --from=builder`, sehingga compiler dan source build tidak
ikut masuk ke image akhir. Hasil pemeriksaan user menunjukkan proses
berjalan sebagai `uid=10001(appuser)`, bukan root.

### D. Tahap 4 (Lab 2 RPS) --- Orkestrasi Multi-Tier dan Persistensi Database

![Bukti Screenshot Tahap
4](docs/screenshots/sc04-tahap4-compose-persistensi.png)

**Analisis Tahap 4**

Docker Compose digunakan untuk menjalankan tiga layanan, yaitu
`frontend` menggunakan Nginx, `backend` menggunakan Go API, dan `db`
menggunakan PostgreSQL 16. Aplikasi dapat diakses melalui
`http://localhost:8080`. Frontend meneruskan permintaan `/api/` ke
backend melalui jaringan Docker, sedangkan backend mengakses database
menggunakan hostname layanan `db` pada port `5432`.

Port PostgreSQL tidak dipublikasikan ke host karena layanan `db` tidak
memiliki konfigurasi `ports:`. Dengan demikian, database diakses melalui
jaringan internal Docker oleh layanan yang membutuhkannya. Data
PostgreSQL disimpan dalam named volume `mk40_pgdata_volume` yang
dipasang ke `/var/lib/postgresql/data`. Pengujian dengan menghentikan
dan menjalankan ulang layanan menunjukkan bahwa data yang sudah
ditambahkan tetap tersedia setelah container dibuat kembali, karena
volume tetap dipertahankan.

### E. Tahap 5 --- Troubleshooting dan Analisis Error

![Bukti Screenshot Tahap
5](docs/screenshots/sc05-tahap5-troubleshooting.png)

**Analisis Tahap 5**

Pada tahap troubleshooting, nilai `DB_HOST` diubah sementara dari `db`
menjadi `localhost`. Di dalam container backend, `localhost` menunjuk ke
container backend itu sendiri, bukan ke container PostgreSQL. Karena
database berada di container terpisah, koneksi ke port `5432` gagal.
Setelah pengujian, nilai `DB_HOST` dikembalikan menjadi `db` dan layanan
dijalankan kembali. Backend kemudian dapat terhubung ke PostgreSQL
melalui jaringan Docker.

Percobaan ini menunjukkan bahwa hostname layanan Docker Compose perlu
digunakan untuk komunikasi antarkontainer pada jaringan yang sama.
Selain itu, `docker compose down` menghapus container dan jaringan
Compose, tetapi named volume tetap dipertahankan. Sebaliknya,
`docker compose down -v` dapat menghapus named volume yang dikelola
Compose, sehingga data database di dalam volume tersebut berisiko
hilang.

------------------------------------------------------------------------

## F. Tantangan Mandiri --- Konfigurasi `.env`, Dual Bridge Network, dan Backend Healthcheck

![Bukti Screenshot Tantangan
Mandiri](docs/screenshots/sc06-tantangan-mandiri.png)

**Status pengerjaan: selesai.** Pengujian peer-review bersama rekan
masih belum dilakukan dan akan dilengkapi kemudian.

### 1. Pemisahan konfigurasi menggunakan `.env`

Konfigurasi database dan identitas mahasiswa dipindahkan ke file `.env`,
sedangkan `docker-compose.yml` menggunakan variabel seperti
`${POSTGRES_USER}`, `${POSTGRES_PASSWORD}`, `${POSTGRES_DB}`,
`${STUDENT_NIM}`, dan `${STUDENT_NAME}`. File `.env` telah tercantum di
`.gitignore`, sehingga tidak muncul pada hasil `git status --short`.
File `.env.example` disediakan sebagai contoh variabel yang dibutuhkan,
tanpa menggunakan password asli.

Validasi dengan `docker compose config --quiet` berhasil. Setelah
konfigurasi diterapkan, ketiga layanan tetap berjalan. Hal ini
menunjukkan bahwa variabel pada `.env` dapat digunakan oleh Docker
Compose untuk mengisi konfigurasi layanan.

### 2. Pemisahan jaringan menjadi dua bridge network

Konfigurasi awal menggunakan satu jaringan bersama. Pada tantangan ini,
jaringan dipisahkan menjadi dua:

-   `mk40_frontend_network` menghubungkan `mk40-frontend` dan
    `mk40-backend`.
-   `mk40_backend_network` menghubungkan `mk40-backend` dan `mk40-db`.

Backend menjadi layanan yang terhubung ke kedua jaringan. Frontend tidak
terhubung langsung ke jaringan database, sedangkan database tidak
terhubung langsung ke jaringan frontend. Hasil `docker network inspect`
menunjukkan bahwa anggota setiap jaringan sesuai dengan rancangan
tersebut. Jaringan lama `mk40_bridge_network` juga telah dihapus setelah
dipastikan tidak memiliki container yang terhubung.

Pemisahan ini membantu membatasi komunikasi antarlayanan berdasarkan
kebutuhan. Frontend dapat meneruskan permintaan ke backend, sementara
backend dapat mengakses PostgreSQL melalui jaringan backend.

### 3. Backend healthcheck

Healthcheck ditambahkan pada layanan `backend` di `docker-compose.yml`.
Pemeriksaan menggunakan `wget` untuk mengakses endpoint
`http://127.0.0.1:8080/healthz`, dengan interval 10 detik, timeout 3
detik, tiga kali percobaan ulang, dan waktu awal 10 detik.

Dockerfile backend telah menyediakan `wget`, sehingga tidak perlu
menambahkan paket tersebut lagi. Setelah konfigurasi diterapkan
menggunakan `docker compose up -d --build`, hasil `docker compose ps`
menunjukkan status backend `Up (healthy)`, database `Up (healthy)`, dan
frontend `Up`. Hasil ini menandakan bahwa endpoint healthcheck backend
berhasil merespons pemeriksaan.

### 4. Pengujian akhir

Setelah ketiga bagian tantangan diterapkan, layanan Docker Compose
berhasil berjalan. Backend dan database terdeteksi sehat, frontend tetap
tersedia melalui port host `8080`, dan jaringan baru sesuai dengan
konfigurasi Dual Bridge Network. Pemeriksaan endpoint `/healthz` dan
`/api/services` dapat digunakan sebagai validasi akhir aplikasi melalui
`http://localhost:8080`.

------------------------------------------------------------------------

## 3. Bukti Audit Silang (Peer-Review) di GitHub

-   **Tautan Pull Request (PR):** Belum tersedia / akan dilengkapi
    setelah proses peer-review.
-   **Nama dan NIM Peer-Reviewer:** Belum dilakukan; akan diisi setelah
    review bersama rekan.
-   **Ringkasan catatan peer-review:** Belum dilakukan.

Bagian peer-review sengaja belum diisi karena proses review bersama
rekan belum dilaksanakan. Setelah sesi peer-review selesai, tambahkan
tautan PR, identitas reviewer, dan ringkasan masukan atau perbaikan yang
benar-benar diberikan.

------------------------------------------------------------------------

## 5. Kesimpulan

Dari Tahap 1 sampai Tahap 5, saya berhasil menjalankan aplikasi Go
secara native, mengemasnya ke dalam Docker dengan single-stage dan
multi-stage build, serta menjalankan frontend, backend, dan PostgreSQL
menggunakan Docker Compose. Perbandingan ukuran image menunjukkan bahwa
multi-stage build membuat image jauh lebih kecil, dari 1.32 GB menjadi
28.2 MB. Selain itu, aplikasi multi-stage berjalan sebagai user
non-root.

Pada Tahap 4, saya mempraktikkan komunikasi antarlayanan, penggunaan
named volume untuk mempertahankan data database, dan pembatasan akses
langsung ke port PostgreSQL. Pada Tahap 5, saya memahami bahwa
penggunaan `localhost` dari dalam container berbeda dengan hostname
layanan Docker Compose, sehingga backend perlu menggunakan `DB_HOST=db`
untuk terhubung ke database.

Pada Tantangan Mandiri, saya berhasil memisahkan konfigurasi ke `.env`,
membuat dua bridge network untuk memisahkan jalur komunikasi frontend
dan database, serta menambahkan healthcheck pada backend. Hasil
pemeriksaan menunjukkan backend dan database berstatus `healthy`. Dengan
demikian, pengerjaan teknis Tahap 1--5 dan Tantangan Mandiri telah
selesai. Kegiatan yang masih tertunda adalah peer-review bersama rekan
dan pengisian bukti review di GitHub.
