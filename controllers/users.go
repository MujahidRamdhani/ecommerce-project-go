package controllers

import (
	"ecommerce-project-go/auth"
	"ecommerce-project-go/entity"
	"ecommerce-project-go/helper"
	"ecommerce-project-go/service"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(ctx *gin.Context) {
	var input entity.InputRegisterUsers

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse(
			"Register failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Register failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := auth.GenerateToken(newUser.ID, newUser.IsAdmin)
	if err != nil {
		response := helper.APIResponse(
			"Register failed",
			http.StatusBadRequest,
			"error",
			"Failed Generate Token",
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	format := helper.FormatUserResponse(newUser, token)
	response := helper.APIResponse("Register success", http.StatusOK, "success", format, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(ctx *gin.Context) {
	var input entity.InputLogin

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse(
			"Login failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Login failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.IsAdmin)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Login failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	format := helper.FormatUserResponse(user, token)
	response := helper.APIResponse("Login success", http.StatusOK, "success", format, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) UpdateUser(ctx *gin.Context) {
	var input entity.InputUpdateUser

	// obtain current user from token
	currentUser := ctx.MustGet("currentUser").(entity.Users)

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse(
			"Update failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedUser, err := h.userService.UpdateUser(currentUser.ID, input)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Update failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := helper.FormatUserEditResponse(updatedUser)
	response := helper.APIResponse("Update success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) DeleteUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(entity.Users)
	err := h.userService.DeleteUser(currentUser.ID)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Delete failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := gin.H{"message": "user has been deleted"}
	response := helper.APIResponse("Delete success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) GetAllUsers(ctx *gin.Context) {
    currentUser := ctx.MustGet("currentUser").(entity.Users)

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))

    users, meta, err := h.userService.GetAll(currentUser.IsAdmin, page, limit)
    if err != nil {
        errorMessage := gin.H{"error": err.Error()}
        response := helper.APIResponse("Get all users failed", http.StatusUnauthorized, "error", errorMessage,
		nil)
        ctx.JSON(http.StatusUnauthorized, response)
        return
    }

    msg := helper.FormatUserGetAllResponse(users)
    response := helper.APIResponse("Get all users success", http.StatusOK, "success", msg, meta)
    ctx.JSON(http.StatusOK, response)
}


