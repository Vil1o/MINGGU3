package siakad

import (
	"errors"
	"net/http"
	"strconv"

	"siakad-service/internal/response"

	"github.com/gin-gonic/gin"
)

type HttpAdapter struct {
	core ServicePort // Memanggil Core via Driving Port
}

func NewHttpAdapter(core ServicePort) *HttpAdapter {
	return &HttpAdapter{core: core}
}

func (h *HttpAdapter) RegisterRoutes(rg *gin.RouterGroup) {
	api := rg.Group("/api/v1")
	{
		api.POST("/mahasiswa", h.TambahMahasiswa)
		api.GET("/mahasiswa", h.DaftarMahasiswa)
		api.GET("/mahasiswa/:nim", h.DetailMahasiswa)
		api.PUT("/mahasiswa/:nim", h.UpdateMahasiswa)
		api.DELETE("/mahasiswa/:nim", h.HapusMahasiswa)
		api.POST("/mahasiswa/:nim/nilai", h.InputNilai)
		api.GET("/mahasiswa/:nim/transkrip", h.Transkrip)
		api.GET("/mahasiswa/:nim/ringkasan", h.Ringkasan) 
	}
}

func mapErrorToStatus(err error) int {
	if errors.Is(err, ErrNIMTidakValid) || errors.Is(err, ErrNilaiTidakValid) || errors.Is(err, ErrSKSTidakValid) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrMahasiswaTidakAda) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrMahasiswaSudahAda) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func (h *HttpAdapter) TambahMahasiswa(c *gin.Context) {
	var input InputMahasiswaDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	if err := h.core.TambahMahasiswa(c.Request.Context(), input); err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusCreated, "Mahasiswa berhasil ditambahkan")
}

func (h *HttpAdapter) DaftarMahasiswa(c *gin.Context) {
	jurusan := c.Query("jurusan")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	data, err := h.core.DaftarMahasiswa(c.Request.Context(), jurusan, limit, offset)
	if err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusOK, data)
}

func (h *HttpAdapter) DetailMahasiswa(c *gin.Context) {
	nim := c.Param("nim")
	data, err := h.core.DetailMahasiswa(c.Request.Context(), nim)
	if err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusOK, data)
}

func (h *HttpAdapter) UpdateMahasiswa(c *gin.Context) {
	nim := c.Param("nim")
	var input InputMahasiswaDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	if err := h.core.UpdateMahasiswa(c.Request.Context(), nim, input); err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusOK, "Mahasiswa berhasil diupdate")
}

func (h *HttpAdapter) HapusMahasiswa(c *gin.Context) {
	nim := c.Param("nim")
	if err := h.core.HapusMahasiswa(c.Request.Context(), nim); err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusOK, "Mahasiswa berhasil dihapus")
}

func (h *HttpAdapter) InputNilai(c *gin.Context) {
	nim := c.Param("nim")
	var input InputNilaiDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	if err := h.core.InputNilai(c.Request.Context(), nim, input); err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusCreated, "Nilai berhasil diinput")
}

func (h *HttpAdapter) Transkrip(c *gin.Context) {
	nim := c.Param("nim")
	data, err := h.core.Transkrip(c.Request.Context(), nim)
	if err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusOK, data)
}

func (h *HttpAdapter) Ringkasan(c *gin.Context) {
	nim := c.Param("nim")
	data, err := h.core.Ringkasan(c.Request.Context(), nim)
	if err != nil {
		response.Error(c, mapErrorToStatus(err), err)
		return
	}
	response.Success(c, http.StatusOK, data)
}

