package response

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Sukses bool `json:"sukses"`
	Data   any  `json:"data"`
}

type ErrorResponse struct {
	Sukses bool   `json:"sukses"`
	Error  string `json:"error"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Sukses: true,
		Data:   data,
	})
}

func Error(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponse{
		Sukses: false,
		Error:  err.Error(),
	})
}