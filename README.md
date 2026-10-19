# AgriSmart

Platform pertanian digital yang membantu petani mengelola kebun dan mengakses pengetahuan budidaya yang bersumber serta tervalidasi.

> **Status:** fondasi backend dan knowledge base sedang dikembangkan. Ini belum merupakan deklarasi siap produksi. Frontend petani dan verifikasi deployment produksi masih perlu diselesaikan.

## Stack

- Backend: Go 1.22
- Database: PostgreSQL 16
- API: JSON over HTTP, prefix `/api/v1`
- Local database: Docker Compose

## Prasyarat

- Go 1.22+
- Docker Engine dan Docker Compose
- Git

## Menjalankan database lokal

Salin konfigurasi contoh dan sesuaikan hanya untuk lingkungan lokal:

```sh
cp .env.example .env
docker compose up -d postgres
sh scripts/verify-database.sh
```

Nilai `change-me-for-local-use` dan `replace-with-at-least-32-random-characters` adalah placeholder pengembangan. Ganti dengan secret kuat yang dikelola secara aman sebelum deployment; jangan gunakan placeholder tersebut di produksi dan jangan commit file `.env`.

## Menjalankan backend

Pastikan PostgreSQL aktif dan migrasi telah diterapkan. Atur `DATABASE_URL` dan `JWT_SECRET` di environment, lalu:

```sh
cd backend
go test ./...
go vet ./...
go build ./cmd/api
go run ./cmd/api
```

Secara default API mendengarkan pada `:8080`; ubah dengan `API_ADDR` bila diperlukan.

## Endpoint API yang tersedia

Semua path di bawah memakai prefix `/api/v1`.

- `GET /health` — pemeriksaan kesehatan API.
- `POST /auth/register`, `POST /auth/login`, `GET /auth/me` — autentikasi.
- `GET/POST /farms`, `GET /farms/{farmID}` — data kebun.
- `GET/POST /farms/{farmID}/plots` — petak kebun.
- `GET/POST /farms/{farmID}/plots/{plotID}/crop-cycles` — siklus tanam.
- `GET /crops` dan endpoint varietas/tahap pertumbuhan — katalog tanaman.
- `GET /knowledge`, `GET /knowledge/fertilizers/recommendations`, `POST /diagnosis` — akses knowledge base.
- Endpoint `/knowledge/documents` — pengelolaan dokumen pengetahuan untuk role berwenang.

Endpoint terlindungi memerlukan token Bearer. Operasi validasi/publikasi pengetahuan memiliki batas role tersendiri. Detail request/response perlu dikonfirmasi melalui handler dan pengujian API sebelum integrasi frontend.

## Verifikasi otomatis

GitHub Actions menjalankan pemeriksaan migrasi PostgreSQL serta tes, `go vet`, dan build backend. Lihat tab Actions untuk hasil terbaru.

## Batasan sebelum pilot/produksi

- Frontend mobile-first untuk alur petani belum tersedia di repository pada status ini.
- Uji end-to-end dari registrasi sampai alur penggunaan petani masih diperlukan.
- Deployment, HTTPS, pengelolaan secret produksi, backup, dan uji restore harus diverifikasi pada lingkungan target.
- Rekomendasi pertanian—terutama dosis, interval, PHI, dan pencampuran pestisida—harus bersumber dan divalidasi sebelum dipublikasikan.

PR dan perubahan pada branch fitur tidak berarti aplikasi sudah di-merge atau dideploy ke server produksi.
