package controllers

import (
	"net/http"
	"working-day-api/internal/services"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	HealthService *services.HealthService
}

func NewHealthController(healthService *services.HealthService) *HealthController {
	return &HealthController{HealthService: healthService}
}

func (ctrl *HealthController) AliveHandler(c *gin.Context) {
	status := ctrl.HealthService.GetStatus()

	c.JSON(http.StatusOK, gin.H{
		"message": status,
	})
}
