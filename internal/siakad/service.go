package siakad

import (
	"context"
	"regexp"
	"sort"
	"strings"
)

// ServicePort adalah Primary Port (Inbound) untuk orkestrasi bisnis
type ServicePort interface {
	TambahMahasiswa(ctx context.Context, input InputMahasiswaDTO) error
	DaftarMahasiswa(ctx context.Context, filterJurusan string, limit, offset int) ([]Mahasiswa, error)
	DetailMahasiswa(ctx context.Context, nim string) (MahasiswaDetailDTO, error)
	UpdateMahasiswa(ctx context.Context, nim string, input InputMahasiswaDTO) error
	HapusMahasiswa(ctx context.Context, nim string) error
	InputNilai(ctx context.Context, nim string, input InputNilaiDTO) error
	Transkrip(ctx context.Context, nim string) (TranskripDTO, error)
	PerJurusan(ctx context.Context) (map[string][]Mahasiswa, error)
	TopIPK(ctx context.Context, n int) ([]MahasiswaDetailDTO, error)
}

// coreService adalah implementasi dari Primary Port, tidak tahu menahu soal SQL atau HTTP
type coreService struct {
	mPort MahasiswaRepository // Menggunakan Secondary Port
	nPort NilaiRepository     // Menggunakan Secondary Port
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

func (s *coreService) PerJurusan(ctx context.Context) (map[string][]Mahasiswa, error) {
	semua, err := s.mPort.Semua(ctx)
	if err != nil {
		return nil, err
	}
	rekap := make(map[string][]Mahasiswa)
	for _, m := range semua {
		rekap[m.Jurusan] = append(rekap[m.Jurusan], m)
	}
	return rekap, nil
}

func (s *coreService) TopIPK(ctx context.Context, n int) ([]MahasiswaDetailDTO, error) {
	if n <= 0 {
		n = 3
	}
	semua, err := s.mPort.Semua(ctx)
	if err != nil {
		return nil, err
	}
	var daftarDetail []MahasiswaDetailDTO
	for _, m := range semua {
		detail, err := s.DetailMahasiswa(ctx, m.NIM)
		if err != nil {
			return nil, err
		}
		daftarDetail = append(daftarDetail, detail)
	}
	sort.Slice(daftarDetail, func(i, j int) bool {
		return daftarDetail[i].IPK > daftarDetail[j].IPK
	})
	if len(daftarDetail) < n {
		n = len(daftarDetail)
	}
	if daftarDetail == nil {
		return []MahasiswaDetailDTO{}, nil
	}
	return daftarDetail[:n], nil
}