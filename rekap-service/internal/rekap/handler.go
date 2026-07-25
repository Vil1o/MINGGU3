package rekap

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HttpAdapter struct {
	service ServicePort
}

func NewHttpAdapter(service ServicePort) *HttpAdapter {
	return &HttpAdapter{service: service}
}


func responseError(c *gin.Context, status int, err error) {
	c.JSON(status, ResponseAPI{Sukses: false, Error: err.Error()})
}

func responseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResponseAPI{Sukses: true, Data: data})
}

func (h *HttpAdapter) TopIPK(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		responseError(c, http.StatusBadRequest, ErrSKSTidakValid) 
		return
	}

	data, err := h.service.TopIPK(c.Request.Context(), limit)
	if err != nil {
		responseError(c, http.StatusInternalServerError, err)
		return
	}

	responseSuccess(c, data)
}

func (h *HttpAdapter) PerJurusan(c *gin.Context) {
	jurusan := c.Query("jurusan")
	if jurusan == "" {
		responseError(c, http.StatusBadRequest, ErrNIMTidakValid)
		return
	}

	data, err := h.service.PerJurusan(c.Request.Context(), jurusan)
	if err != nil {
		responseError(c, http.StatusInternalServerError, err)
		return
	}

	responseSuccess(c, data)
}