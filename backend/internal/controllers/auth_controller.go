package controllers

import (
	"net/http"
	"strings"

	"employee-platform/backend/internal/dto"
	"employee-platform/backend/internal/services"
	"employee-platform/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ctl *AuthController) AdminSignup(c *gin.Context) {
	var req dto.AdminSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := ctl.service.AdminSignup(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "admin signup successful", gin.H{
		"user":  user,
		"token": token,
	})
}

func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := ctl.service.Login(c.Request.Context(), req)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

func (ctl *AuthController) Me(c *gin.Context) {
	userID := c.GetString("userID")
	user, profile, err := ctl.service.Me(c.Request.Context(), userID)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "me fetched successfully", gin.H{
		"user":    user,
		"profile": profile,
	})
}

func ExtractBearerToken(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}