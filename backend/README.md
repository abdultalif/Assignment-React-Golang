# 🚀 Cara Menjalankan Program

## 1. Persiapan Environment

```bash
# Salin file konfigurasi environment
cp .env.example .env
```

Sesuaikan konfigurasi dalam file `.env` dengan environment Anda.

## 2. Instalasi Dependencies

```bash
go mod tidy
```

## 2. Menjalankan Docker Services

```bash
docker compose up -d
```

## 4. Migrasi Database

```bash
migrate -database "postgres://username:password@localhost:5432/db?sslmode=disable" -path database/migrations up
```

**⚠️ Penting:** Ganti `username`, `password`, dan `db` sesuai dengan konfigurasi PostgreSQL Anda.

## 5. Menjalankan Aplikasi

```bash
# Jalankan aplikasi utama
go run main.go
```
