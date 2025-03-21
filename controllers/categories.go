package controllers

import (
	"ecommerce-project-go/entity"
	"ecommerce-project-go/helper"
	"ecommerce-project-go/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type categoriesHandler struct {
	catService service.CategoriesService
}

func NewCatHandler(categoriesService service.CategoriesService) *categoriesHandler {
	return &categoriesHandler{categoriesService}
}

func (h *categoriesHandler) GetAllCategories(ctx *gin.Context) {
	categories, err := h.catService.GetAll()
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Get all categories failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := helper.FormatCatGetAllResponse(categories)
	response := helper.APIResponse("Get all users success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *categoriesHandler) InsertCategory(ctx *gin.Context) {
	var cat entity.Category
	currentUser := ctx.MustGet("currentUser").(entity.Users)

	err := ctx.ShouldBindJSON(&cat)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse(
			"Insert category failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newCat, err := h.catService.InsertCategory(cat, currentUser.IsAdmin)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Insert category failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := gin.H{
		"id":    newCat.Id,
		"name":  newCat.Name,
		"admin": currentUser.FullName,
	}
	response := helper.APIResponse("Insert category success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *categoriesHandler) UpdateCategory(ctx *gin.Context) {
	var cat entity.Category

	id, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Update category failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(entity.Users)

	err = ctx.ShouldBindJSON(&cat)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse(
			"Update category failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updatedCat, err := h.catService.EditCategory(cat, id, currentUser.IsAdmin)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Update category failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := gin.H{
		"id":    updatedCat.Id,
		"name":  updatedCat.Name,
		"admin": currentUser.FullName,
	}
	response := helper.APIResponse("Update category success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *categoriesHandler) DeleteCategories(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Delete category failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(entity.Users)

	err = h.catService.DeleteCategory(id, currentUser.IsAdmin)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Delete category failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := gin.H{
		"message": "Delete success",
	}
	response := helper.APIResponse("Delete category success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *categoriesHandler) GetAllProductsByCatId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Get all item with such categories id failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	Products, err := h.catService.GetAllProduct(id)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Get all item with such categories id failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	format := helper.FormatCatProductGetResponse(Products)
	response := helper.APIResponse("Get success", http.StatusOK, "success", format, nil)
	ctx.JSON(http.StatusOK, response)
}
