package siakad

import (
	"context"
	"errors"
	"math"
	"testing"
)

// --- Table-Driven Test untuk HitungIPK ---
func TestHitungIPK(t *testing.T) {
	tests := []struct {
		name   string
		daftar []Nilai
		want   float64
	}{
		{
			name:   "Skenario 1: Daftar nilai kosong",
			daftar: []Nilai{},
			want:   0.0,
		},
		{
			name:   "Skenario 2: 1 mata kuliah",
			daftar: []Nilai{{SKS: 3, Mutu: 4.0}},
			want:   4.0,
		},
		{
			name:   "Skenario 3: 2 mata kuliah beda SKS",
			daftar: []Nilai{{SKS: 3, Mutu: 4.0}, {SKS: 2, Mutu: 3.0}}, // (12 + 6) / 5 = 18/5 = 3.6
			want:   3.6,
		},
		{
			name: "Skenario 4: Kasus kompleks",
			daftar: []Nilai{
				{SKS: 4, Mutu: 3.5}, // 14
				{SKS: 2, Mutu: 4.0}, // 8
				{SKS: 3, Mutu: 2.0}, // 6
				{SKS: 1, Mutu: 1.0}, // 1
			}, // Total: 29 / 10 SKS = 2.9
			want: 2.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HitungIPK(tt.daftar)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("HitungIPK() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Mocks ---
type mockMahasiswaRepo struct {
	cariFunc func(nim string) (Mahasiswa, error)
}

func (m *mockMahasiswaRepo) Tambah(ctx context.Context, mhs Mahasiswa) error { return nil }
func (m *mockMahasiswaRepo) Semua(ctx context.Context) ([]Mahasiswa, error)  { return nil, nil }
func (m *mockMahasiswaRepo) Update(ctx context.Context, nim string, mhs Mahasiswa) error { return nil }
func (m *mockMahasiswaRepo) Hapus(ctx context.Context, nim string) error     { return nil }
func (m *mockMahasiswaRepo) Cari(ctx context.Context, nim string) (Mahasiswa, error) {
	return m.cariFunc(nim)
}

type mockNilaiRepo struct {
	tambahFunc func(n Nilai) error
}

func (m *mockNilaiRepo) PerMahasiswa(ctx context.Context, nim string) ([]Nilai, error) {
	return nil, nil
}
func (m *mockNilaiRepo) Tambah(ctx context.Context, n Nilai) error {
	return m.tambahFunc(n)
}

// --- Mock Test: InputNilai jika Mahasiswa Tidak Ada ---
func TestInputNilai_MahasiswaTidakAda(t *testing.T) {
	mRepo := &mockMahasiswaRepo{
		cariFunc: func(nim string) (Mahasiswa, error) {
			return Mahasiswa{}, ErrMahasiswaTidakAda
		},
	}
	nRepo := &mockNilaiRepo{}

	svc := NewService(mRepo, nRepo)

	err := svc.InputNilai(context.Background(), "12345678", InputNilaiDTO{
		KodeMK: "CS101",
		NamaMK: "Intro",
		SKS:    3,
		Mutu:   3.5,
	})

	if !errors.Is(err, ErrMahasiswaTidakAda) {
		t.Errorf("Diharapkan error ErrMahasiswaTidakAda, namun mendapatkan: %v", err)
	}
}