package controllers

import (
	"net/http"
	"working-day-api/internal/domain"
	"working-day-api/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	LoginService *services.LoginService
}

func NewAuthController(loginService *services.LoginService) *AuthController {
	return &AuthController{
		LoginService: loginService,
	}
}

func (ctrl *AuthController) LoginHandler(c *gin.Context) {
	var loginRequest domain.Login

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON payload: " + err.Error(),
		})
		return
	}

	token, err := ctrl.LoginService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
