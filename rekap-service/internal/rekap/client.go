package rekap

import (
	"context"
)

// AkademikClient sejajar dengan repository, namun mengambil data via HTTP
type AkademikClient interface {
	DaftarMahasiswa(ctx context.Context) ([]Mahasiswa, error)
	Ringkasan(ctx context.Context, nim string) (Ringkasan, error)
}