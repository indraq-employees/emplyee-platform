package controllers

import (
	"net/http"

	"employee-platform/backend/internal/dto"
	"employee-platform/backend/internal/services"
	"employee-platform/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	service *services.ProfileService
}

func NewProfileController(service *services.ProfileService) *ProfileController {
	return &ProfileController{service: service}
}

func (ctl *ProfileController) CompleteProfile(c *gin.Context) {
	var req dto.CompleteProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("userID")
	if err := ctl.service.CompleteProfile(c.Request.Context(), userID, req); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "profile completed successfully", nil)
}