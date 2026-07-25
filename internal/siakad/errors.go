package siakad

import "errors"

var (
	ErrNIMTidakValid = errors.New("NIM tidak valid")
	ErrMahasiswaSudahAda = errors.New("NIM sudah terdaftar")
	ErrMahasiswaTidakAda = errors.New("NIM tidak ditemukan")
	ErrNilaiTidakValid = errors.New("Mutu di luar 0.0-4.0")
	ErrSKSTidakValid = errors.New("SKS di luar <= 0")
)