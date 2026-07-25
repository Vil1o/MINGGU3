package siakad

import (
	"context"
	"regexp"
	"strings"
)

type ServicePort interface {
	TambahMahasiswa(ctx context.Context, input InputMahasiswaDTO) error
	DaftarMahasiswa(ctx context.Context, filterJurusan string, limit, offset int) ([]Mahasiswa, error)
	DetailMahasiswa(ctx context.Context, nim string) (MahasiswaDetailDTO, error)
	UpdateMahasiswa(ctx context.Context, nim string, input InputMahasiswaDTO) error
	HapusMahasiswa(ctx context.Context, nim string) error
	InputNilai(ctx context.Context, nim string, input InputNilaiDTO) error
	Transkrip(ctx context.Context, nim string) (TranskripDTO, error)
	Ringkasan(ctx context.Context, nim string) (RingkasanDTO, error) // <-- TAMBAH INI, ganti PerJurusan & TopIPK
}

type coreService struct {
	mPort MahasiswaRepository 
	nPort NilaiRepository     
}

func NewService(mPort MahasiswaRepository, nPort NilaiRepository) ServicePort {
	return &coreService{mPort: mPort, nPort: nPort}
}

// HitungIPK: Fungsi Domain Murni
func HitungIPK(daftar []Nilai) float64 {
	var totalSKS int
	var totalBobot float64

	for _, n := range daftar {
		totalSKS += n.SKS
		totalBobot += float64(n.SKS) * n.Mutu
	}

	if totalSKS == 0 {
		return 0.0
	}
	return totalBobot / float64(totalSKS)
}

func (s *coreService) TambahMahasiswa(ctx context.Context, input InputMahasiswaDTO) error {
	if matched, _ := regexp.MatchString(`^[0-9]{8}$`, input.NIM); !matched {
		return ErrNIMTidakValid
	}
	m := Mahasiswa{NIM: input.NIM, Nama: input.Nama, Jurusan: input.Jurusan, Status: input.Status}
	return s.mPort.Tambah(ctx, m)
}

func (s *coreService) DaftarMahasiswa(ctx context.Context, filterJurusan string, limit, offset int) ([]Mahasiswa, error) {
	semua, err := s.mPort.Semua(ctx)
	if err != nil {
		return nil, err
	}
	var filtered []Mahasiswa
	for _, m := range semua {
		if filterJurusan == "" || strings.EqualFold(m.Jurusan, filterJurusan) {
			filtered = append(filtered, m)
		}
	}
	if offset >= len(filtered) {
		return []Mahasiswa{}, nil
	}
	end := offset + limit
	if end > len(filtered) || limit <= 0 {
		end = len(filtered)
	}
	return filtered[offset:end], nil
}

func (s *coreService) DetailMahasiswa(ctx context.Context, nim string) (MahasiswaDetailDTO, error) {
	m, err := s.mPort.Cari(ctx, nim)
	if err != nil {
		return MahasiswaDetailDTO{}, err
	}
	nilai, err := s.nPort.PerMahasiswa(ctx, nim)
	if err != nil {
		return MahasiswaDetailDTO{}, err
	}
	ipk := HitungIPK(nilai)
	statusCumlaude := ""
	if ipk >= 3.50 {
		statusCumlaude = "Cumlaude"
	}
	return MahasiswaDetailDTO{
		NIM: m.NIM, Nama: m.Nama, Jurusan: m.Jurusan, Status: m.Status, IPK: ipk, StatusCumlaude: statusCumlaude,
	}, nil
}

func (s *coreService) UpdateMahasiswa(ctx context.Context, nim string, input InputMahasiswaDTO) error {
	return s.mPort.Update(ctx, nim, Mahasiswa{NIM: nim, Nama: input.Nama, Jurusan: input.Jurusan, Status: input.Status})
}

func (s *coreService) HapusMahasiswa(ctx context.Context, nim string) error {
	return s.mPort.Hapus(ctx, nim)
}

func (s *coreService) InputNilai(ctx context.Context, nim string, input InputNilaiDTO) error {
	if input.SKS <= 0 {
		return ErrSKSTidakValid
	}
	if input.Mutu < 0.0 || input.Mutu > 4.0 {
		return ErrNilaiTidakValid
	}
	if _, err := s.mPort.Cari(ctx, nim); err != nil {
		return err
	}
	return s.nPort.Tambah(ctx, Nilai{
		NIM: nim, KodeMK: input.KodeMK, NamaMK: input.NamaMK, SKS: input.SKS, Mutu: input.Mutu,
	})
}

func (s *coreService) Transkrip(ctx context.Context, nim string) (TranskripDTO, error) {
	detail, err := s.DetailMahasiswa(ctx, nim)
	if err != nil {
		return TranskripDTO{}, err
	}
	nilai, err := s.nPort.PerMahasiswa(ctx, nim)
	if err != nil {
		return TranskripDTO{}, err
	}
	totalSKS := 0
	for _, n := range nilai {
		totalSKS += n.SKS
	}
	if nilai == nil {
		nilai = []Nilai{}
	}
	return TranskripDTO{MahasiswaDetailDTO: detail, DaftarNilai: nilai, TotalSKS: totalSKS}, nil
}

func (s *coreService) Ringkasan(ctx context.Context, nim string) (RingkasanDTO, error) {
	m, err := s.mPort.Cari(ctx, nim)
	if err != nil {
		return RingkasanDTO{}, err
	}
	
	nilai, err := s.nPort.PerMahasiswa(ctx, nim)
	if err != nil {
		return RingkasanDTO{}, err
	}
	
	totalSKS := 0
	for _, n := range nilai {
		totalSKS += n.SKS
	}
	
	ipk := HitungIPK(nilai)
	
	return RingkasanDTO{
		NIM:      m.NIM,
		Nama:     m.Nama,
		Jurusan:  m.Jurusan,
		Status:   m.Status,
		TotalSKS: totalSKS,
		IPK:      ipk,
	}, nil
}