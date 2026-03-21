package controllers

import (
	"net/http"

	"employee-platform/backend/internal/dto"
	"employee-platform/backend/internal/services"
	"employee-platform/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	service *services.EmployeeService
}

func NewEmployeeController(service *services.EmployeeService) *EmployeeController {
	return &EmployeeController{service: service}
}

func (ctl *EmployeeController) CreateEmployee(c *gin.Context) {
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	actorUserID := c.GetString("userID")
	user, profile, loginLink, err := ctl.service.CreateEmployee(c.Request.Context(), actorUserID, req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "employee created successfully", gin.H{
		"user":      user,
		"profile":   profile,
		"loginLink": loginLink,
	})
}

func (ctl *EmployeeController) ListEmployees(c *gin.Context) {
	users, err := ctl.service.ListEmployees(c.Request.Context())
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "employees fetched successfully", users)
}