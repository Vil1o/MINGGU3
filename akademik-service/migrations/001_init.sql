CREATE TABLE IF NOT EXISTS mahasiswa (
    nim VARCHAR(8) PRIMARY KEY CHECK (LENGTH(nim) = 8),
    nama VARCHAR NOT NULL,
    jurusan VARCHAR NOT NULL,
    status VARCHAR NOT NULL CHECK (status IN ('Aktif', 'Cuti', 'Lulus'))
);

CREATE TABLE IF NOT EXISTS nilai (
    id BIGSERIAL PRIMARY KEY,
    nim VARCHAR(8) NOT NULL REFERENCES mahasiswa(nim) ON DELETE CASCADE,
    kode_mk VARCHAR NOT NULL,
    nama_mk VARCHAR NOT NULL,
    sks INT NOT NULL CHECK (sks > 0),
    mutu NUMERIC(3, 2) NOT NULL CHECK (mutu >= 0.0 AND mutu <= 4.0)
);