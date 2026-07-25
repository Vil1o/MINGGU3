package rekap

import (
	"context"
	"sort"
)

type ServicePort interface {
	TopIPK(ctx context.Context, limit int) ([]Ringkasan, error)
	PerJurusan(ctx context.Context, jurusan string) ([]Ringkasan, error)
}

type coreService struct {
	client AkademikClient
}

func NewService(client AkademikClient) ServicePort {
	return &coreService{client: client}
}

func (s *coreService) TopIPK(ctx context.Context, limit int) ([]Ringkasan, error) {
	mahasiswa, err := s.client.DaftarMahasiswa(ctx)
	if err != nil {
		return nil, err
	}

	var daftarRingkasan []Ringkasan
	for _, m := range mahasiswa {

		ringkasan, err := s.client.Ringkasan(ctx, m.NIM)
		if err == nil { 
			daftarRingkasan = append(daftarRingkasan, ringkasan)
		}
	}


	sort.Slice(daftarRingkasan, func(i, j int) bool {
		return daftarRingkasan[i].IPK > daftarRingkasan[j].IPK
	})


	if len(daftarRingkasan) > limit {
		daftarRingkasan = daftarRingkasan[:limit]
	}

	return daftarRingkasan, nil
}

func (s *coreService) PerJurusan(ctx context.Context, jurusan string) ([]Ringkasan, error) {
	mahasiswa, err := s.client.DaftarMahasiswa(ctx)
	if err != nil {
		return nil, err
	}

	var daftarRingkasan []Ringkasan
	for _, m := range mahasiswa {
		if m.Jurusan == jurusan {
			ringkasan, err := s.client.Ringkasan(ctx, m.NIM)
			if err == nil {
				daftarRingkasan = append(daftarRingkasan, ringkasan)
			}
		}
	}

	return daftarRingkasan, nil
}