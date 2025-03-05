package controllers

import (
	"ecommerce-project-go/entity"
	"ecommerce-project-go/helper"
	"ecommerce-project-go/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(ProductService service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService}
}

func (h *ProductHandler) InsertProduct(ctx *gin.Context) {
	var inputProduct entity.InputProduct
	currentUser := ctx.MustGet("currentUser").(entity.Users)

	err := ctx.ShouldBindJSON(&inputProduct)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse(
			"Insert Product failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	Product, stock, err := h.service.AddItem(inputProduct, currentUser.IsAdmin)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Insert Product failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := helper.FormatProductSaveResponse(Product, stock)
	response := helper.APIResponse("Insert success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var input entity.UpdateProduct
	currentUser := ctx.MustGet("currentUser").(entity.Users)
	id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Update Product failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"error": errors}
		response := helper.APIResponse(
			"Update Product failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	Product, stock, err := h.service.UpdateItem(input, currentUser.IsAdmin, id)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Update Product failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	msg := helper.FormatProductSaveResponse(Product, stock)
	response := helper.APIResponse("Update success", http.StatusOK, "success", msg, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(entity.Users)
	id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Delete Product failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.service.DeleteItem(currentUser.IsAdmin, id)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Delete Product failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Delete success", http.StatusOK, "success", nil, nil)
	ctx.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "5"))
	products, meta, err := h.service.GetAll(page, limit)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Get all Product failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Get success", http.StatusOK, "success", products, meta)
	ctx.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("product_id"))
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Get Product failed",
			http.StatusUnprocessableEntity,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	Product, err := h.service.GetById(id)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse(
			"Get Product failed",
			http.StatusBadRequest,
			"error",
			errorMessage,
			nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Get success", http.StatusOK, "success", Product, nil)
	ctx.JSON(http.StatusOK, response)
}
