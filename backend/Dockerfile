# =====================================================================
# STAGE 1: LINGKUNGAN KOMPILASI (BUILDER)
# =====================================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

# [OPTIMASI LAYER CACHING]:
# Salin file definisi library (go.mod & go.sum) lebih dulu sebelum main.go.
# Tujuannya: Jika kita hanya mengubah isi main.go, Docker TIDAK perlu mengunduh ulang library!
COPY go.mod go.sum ./
RUN go mod download

# Salin kode sumber aplikasi
COPY main.go ./

# Kompilasi menjadi file biner mandiri (CGO_ENABLED=0) dan pangkas simbol debug (-ldflags="-s -w")
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /build/app-server main.go

# =====================================================================
# STAGE 2: LINGKUNGAN PRODUKSI MINIMALIS (TARGET < 100 MB)
# =====================================================================
FROM alpine:3.20 AS production

# Tambahkan wget untuk keperluan healthcheck dan buat akun user non-root (keamanan)
RUN apk --no-cache add ca-certificates wget && \
    addgroup -S appgroup && \
    adduser -S -u 10001 -G appgroup appuser

WORKDIR /app

# KUNCI MULTI-STAGE BUILD:
# Salin HANYA file biner '/build/app-server' dari Stage 'builder' ke image produksi ini!
COPY --from=builder --chown=appuser:appgroup /build/app-server /app/app-server

# Jalankan aplikasi sebagai user biasa (bukan root)
USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/app-server"]