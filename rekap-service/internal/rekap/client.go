package rekap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)


type AkademikClient interface {
	DaftarMahasiswa(ctx context.Context) ([]Mahasiswa, error)
	Ringkasan(ctx context.Context, nim string) (Ringkasan, error)
}

type HTTPAkademikClient struct {
	baseURL string
	client  *http.Client
}

// NewHTTPAkademikClient membuat instance client baru
func NewHTTPAkademikClient(baseURL string) *HTTPAkademikClient {
	return &HTTPAkademikClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 2 * time.Second, // Batas waktu dari level client
		},
	}
}

// mapNetworkError menerjemahkan error jaringan Golang menjadi error spesifik domain
func (c *HTTPAkademikClient) mapNetworkError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrAkademikTimeout // Timeout spesifik jika melewati batas waktu
	}
	// Membungkus (wrap) error asli dengan %w
	return fmt.Errorf("%w: %v", ErrAkademikTidakTersedia, err)
}

func (c *HTTPAkademikClient) DaftarMahasiswa(ctx context.Context) ([]Mahasiswa, error) {
	// Batasi waktu request ini maksimal 2 detik sesuai anggaran timeout[cite: 1]
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel() // Wajib agar tidak terjadi memory leak

	url := fmt.Sprintf("%s/api/v1/mahasiswa", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) // Meneruskan context[cite: 1]
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, c.mapNetworkError(err)
	}
	defer resp.Body.Close() // Wajib dipanggil untuk mencegah kebocoran koneksi[cite: 1]

	if resp.StatusCode >= 500 {
		return nil, ErrAkademikTidakTersedia
	}

	var apiResp struct {
		Sukses bool        `json:"sukses"`
		Data   []Mahasiswa `json:"data"`
		Error  string      `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("gagal decode JSON: %w", err)
	}

	return apiResp.Data, nil
}

func (c *HTTPAkademikClient) Ringkasan(ctx context.Context, nim string) (Ringkasan, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second) // Anggaran batas HTTP 2s[cite: 1]
	defer cancel()

	url := fmt.Sprintf("%s/api/v1/mahasiswa/%s/ringkasan", c.baseURL, nim)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Ringkasan{}, fmt.Errorf("gagal membuat request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Ringkasan{}, c.mapNetworkError(err)
	}
	defer resp.Body.Close()

	// Menerjemahkan status 404 dari Akademik Service menjadi error domain[cite: 1]
	if resp.StatusCode == http.StatusNotFound {
		return Ringkasan{}, ErrMahasiswaTidakAda
	}
	if resp.StatusCode >= 500 {
		return Ringkasan{}, ErrAkademikTidakTersedia
	}

	var apiResp struct {
		Sukses bool      `json:"sukses"`
		Data   Ringkasan `json:"data"`
		Error  string    `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return Ringkasan{}, fmt.Errorf("gagal decode JSON: %w", err)
	}

	return apiResp.Data, nil
}