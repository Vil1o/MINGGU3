# ==========================================
# STAGE 1: Builder
# ==========================================
FROM golang:1.26-alpine AS builder

# Mematikan CGO agar binary benar-benar statis dan berjalan tanpa libc tambahan
ENV CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

# Memanfaatkan layer cache: copy go.mod & go.sum terlebih dahulu
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary dengan nama "akademik-app"
RUN go build -o akademik-app .

# ==========================================
# STAGE 2: Runner
# ==========================================
FROM alpine:3.21

WORKDIR /app

# Menambahkan ca-certificates untuk kebutuhan panggilan HTTPS keluar (jika ada nanti)
RUN apk --no-cache add ca-certificates

# Membuat user non-root bernama "appuser"
RUN adduser -D appuser

# Copy binary dari stage builder ke stage runner
COPY --from=builder /app/akademik-app .

# Ubah kepemilikan file ke appuser
RUN chown -R appuser:appuser /app

# Gunakan user non-root
USER appuser

# Expose port internal (di dalam container selalu 8080 sesuai spesifikasi)
EXPOSE 8080

# Perintah untuk menjalankan aplikasi
CMD ["./akademik-app"]