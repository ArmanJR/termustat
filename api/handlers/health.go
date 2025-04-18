package handlers

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	logger *zap.Logger
}

func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{logger: logger}
}

// HealthCheck responds with a simple “ok” status.
// @Summary      Health Check
// @Description  Returns a 200 OK if the service is running
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string  "status: ok"
// @Router       /v1/health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
