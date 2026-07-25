package rekap

import "errors"

// Error domain yang disalin dari Akademik Service
var (
	ErrNIMTidakValid     = errors.New("NIM tidak valid")
	ErrMahasiswaSudahAda = errors.New("mahasiswa sudah terdaftar")
	ErrMahasiswaTidakAda = errors.New("mahasiswa tidak ditemukan")
	ErrNilaiTidakValid   = errors.New("nilai tidak valid")
	ErrSKSTidakValid     = errors.New("sks tidak valid")
)

// BARU: Error khusus kegagalan jaringan (Tahap 5)
var (
	ErrAkademikTimeout       = errors.New("layanan akademik tidak merespons, coba lagi")
	ErrAkademikTidakTersedia = errors.New("layanan akademik sedang tidak tersedia")
)