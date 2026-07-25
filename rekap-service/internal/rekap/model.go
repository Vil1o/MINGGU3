package rekap

// Mahasiswa dibutuhkan untuk menyusun rekap jurusan
type Mahasiswa struct {
	NIM     string `json:"nim"`
	Nama    string `json:"nama"`
	Jurusan string `json:"jurusan"`
	Status  string `json:"status"`
}

// Ringkasan adalah kontrak data dari Akademik Service
type Ringkasan struct {
	NIM      string  `json:"nim"`
	Nama     string  `json:"nama"`
	Jurusan  string  `json:"jurusan"`
	Status   string  `json:"status"`
	TotalSKS int     `json:"total_sks"`
	IPK      float64 `json:"ipk"`
}

// Amplop API dari Akademik Service
type ResponseAPI struct {
	Sukses bool        `json:"sukses"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}