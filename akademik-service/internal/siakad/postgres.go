package siakad

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type mahasiswaRepo struct {
	db *sql.DB
}

func NewMahasiswaRepository(db *sql.DB) MahasiswaRepository {
	return &mahasiswaRepo{db: db}
}

func (r *mahasiswaRepo) Tambah(ctx context.Context, m Mahasiswa) error {
	query := `INSERT INTO mahasiswa (nim, nama, jurusan, status) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, m.NIM, m.Nama, m.Jurusan, m.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: %v", ErrMahasiswaSudahAda, err)
		}
		return err
	}
	return nil
}

func (r *mahasiswaRepo) Cari(ctx context.Context, nim string) (Mahasiswa, error) {
	query := `SELECT nim, nama, jurusan, status FROM mahasiswa WHERE nim = $1`
	var m Mahasiswa
	err := r.db.QueryRowContext(ctx, query, nim).Scan(&m.NIM, &m.Nama, &m.Jurusan, &m.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Mahasiswa{}, ErrMahasiswaTidakAda
		}
		return Mahasiswa{}, err
	}
	return m, nil
}

func (r *mahasiswaRepo) Semua(ctx context.Context) ([]Mahasiswa, error) {
	query := `SELECT nim, nama, jurusan, status FROM mahasiswa`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Mahasiswa
	for rows.Next() {
		var m Mahasiswa
		if err := rows.Scan(&m.NIM, &m.Nama, &m.Jurusan, &m.Status); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *mahasiswaRepo) Update(ctx context.Context, nim string, m Mahasiswa) error {
	m.NIM = nim
	query := `UPDATE mahasiswa SET nama = $1, jurusan = $2, status = $3 WHERE nim = $4`
	res, err := r.db.ExecContext(ctx, query, m.Nama, m.Jurusan, m.Status, m.NIM)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMahasiswaTidakAda
	}
	return nil
}

func (r *mahasiswaRepo) Hapus(ctx context.Context, nim string) error {
	query := `DELETE FROM mahasiswa WHERE nim = $1`
	res, err := r.db.ExecContext(ctx, query, nim)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrMahasiswaTidakAda
	}
	return nil
}

type nilaiRepo struct {
	db *sql.DB
}

func NewNilaiRepository(db *sql.DB) NilaiRepository {
	return &nilaiRepo{db: db}
}

func (r *nilaiRepo) Tambah(ctx context.Context, n Nilai) error {
	query := `INSERT INTO nilai (nim, kode_mk, nama_mk, sks, mutu) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, n.NIM, n.KodeMK, n.NamaMK, n.SKS, n.Mutu)
	return err
}

func (r *nilaiRepo) PerMahasiswa(ctx context.Context, nim string) ([]Nilai, error) {
	query := `SELECT id, nim, kode_mk, nama_mk, sks, mutu FROM nilai WHERE nim = $1`
	rows, err := r.db.QueryContext(ctx, query, nim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Nilai
	for rows.Next() {
		var n Nilai
		if err := rows.Scan(&n.ID, &n.NIM, &n.KodeMK, &n.NamaMK, &n.SKS, &n.Mutu); err != nil {
			return nil, err
		}
		list = append(list, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return list, nil
}