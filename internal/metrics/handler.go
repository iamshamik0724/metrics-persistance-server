package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service IService
}

func NewHandler(metricsService IService) *Handler {
	return &Handler{
		Service: metricsService,
	}
}

func (h *Handler) GetMetrics(c *gin.Context) {
	data, err := h.Service.FetchLast10MinutesMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
