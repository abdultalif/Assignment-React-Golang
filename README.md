# Order & Settlement Processing System

## Overview

Sistem backend untuk menangani pemrosesan order dan settlement transaksi bervolume tinggi. Dibangun dengan fokus pada skalabilitas, pemrosesan asinkron, dan efisiensi data menggunakan worker pool dan batch processing.

---

## Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL
- **Architecture:** Clean Architecture (Handler → Service → Repository)
- **Concurrency:** Goroutines, Channels, Worker Pool

---

## Features

- Order management dengan RESTful API
- Atomic stock update untuk konsistensi data
- Asynchronous job processing untuk workload berat
- Worker pool dengan concurrency terkontrol
- Batch processing untuk agregasi transaksi skala besar
- Job lifecycle management: `QUEUED` → `RUNNING` → `COMPLETED` / `FAILED` / `CANCELLED`
- CSV report generation untuk hasil settlement

---

## How It Works

### Order Flow

1. User membuat order via API
2. Sistem memvalidasi produk dan mengupdate stok secara atomic
3. Order disimpan ke database

### Settlement Flow

1. User men-trigger settlement job dengan rentang tanggal
2. Job dibuat dengan status `QUEUED`
3. Worker pool memproses job secara asinkron
4. Transaksi diambil dalam batch dan diagregasi per merchant
5. Hasil disimpan dan diekspor sebagai file CSV

---

## Getting Started

### 1. Clone & Persiapan Environment

```bash
cp .env.example .env
```

Sesuaikan konfigurasi dalam `.env` dengan environment Anda.

### 2. Instalasi Dependencies

```bash
go mod tidy
```

### 3. Jalankan Docker Services

```bash
docker compose up -d
```

### 4. Migrasi Database

```bash
migrate -database "postgres://username:password@localhost:5432/db?sslmode=disable" \
  -path database/migrations up
```

> ⚠️ Ganti `username`, `password`, dan `db` sesuai konfigurasi PostgreSQL Anda.

### 5. Jalankan Aplikasi

```bash
go run main.go
```

---

## API Documentation

Dokumentasi lengkap API beserta contoh request/response tersedia di Postman:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://documenter.getpostman.com/view/47767603/2sB3HopzoE)
