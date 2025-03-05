package helper

import (
	"ecommerce-project-go/entity"
	"time"

	"github.com/go-playground/validator/v10"
)

func FormatError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}

type Response struct {
	Message string
	Code    int
	Status  string
	Data    any `json:"data,omitempty"` 
	Meta    any `json:"meta,omitempty"`
}

type Meta struct {
	TotalRecords int 
	TotalPages   int 
	CurrentPage  int 
	PerPage      int 
}

func APIResponse(message string, code int, status string,  data any, meta any) Response {

	return Response{
		Message: message,
		Code:    code,
		Status:  status,
		Data:    data,
		Meta:    meta,
	}
}

type UserResponseFormat struct {
	ID       int
	FullName string
	Email    string
	Token    string
}

func FormatUserResponse(user entity.Users, token string) UserResponseFormat {
	response := UserResponseFormat{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Token:    token,
	}

	return response
}

type UserEditResponseFormat struct {
	ID       int
	FullName string
	Email    string
}

func FormatUserEditResponse(user entity.Users) UserEditResponseFormat {
	response := UserEditResponseFormat{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}

	return response
}

type UserGetAllFormat struct {
	ID                   int
	FullName, Email      string
	CreatedAt, UpdatedAt time.Time
}

func FormatUserGetAllResponse(users []entity.Users) []UserGetAllFormat {
	var response []UserGetAllFormat
	for _, user := range users {
		res := UserGetAllFormat{
			ID:        user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		response = append(response, res)
	}

	return response
}

type CatGetAllFormat struct {
	ID   int
	Name string
}

func FormatCatGetAllResponse(cats []entity.Category) []CatGetAllFormat {
	var response []CatGetAllFormat
	for _, cat := range cats {
		res := CatGetAllFormat{
			ID:   cat.Id,
			Name: cat.Name,
		}
		response = append(response, res)
	}

	return response
}

type CatProductGetFormat struct {
	ID, CatID         int
	Name, Description string
	IsAvailable       bool
}

func FormatCatProductGetResponse(Products []entity.Product) []CatProductGetFormat {
	var response []CatProductGetFormat
	for _, Product := range Products {
		res := CatProductGetFormat{
			ID:          Product.Id,
			CatID:       Product.CatId,
			Name:        Product.Name,
			Description: Product.Description,
			IsAvailable: Product.IsAvailable,
		}
		response = append(response, res)
	}

	return response
}

type ProductFormat struct {
	ID, CatID               int
	Name, Description       string
	IsAvailable             bool
	StockUnit, PricePerUnit int
}

func FormatProductSaveResponse(Product entity.Product, stock entity.Stock) ProductFormat {
	response := ProductFormat{
		ID:           Product.Id,
		CatID:        Product.CatId,
		Name:         Product.Name,
		Description:  Product.Description,
		IsAvailable:  Product.IsAvailable,
		StockUnit:    stock.StockUnit,
		PricePerUnit: stock.PricePerUnit,
	}

	return response
}

type OutputTransaction struct {
	Information string
	Id          int
	UserId      int
	Item        string
	Unit        int
	TotalPrice  int
	Status      string
	StartAt     time.Time
	FinishAt    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiredAt   time.Time
}

func FormatTransactionResponse(transaction entity.Transaction, item string, info string) OutputTransaction {
	response := OutputTransaction{
		Information: info,
		Id:          transaction.Id,
		UserId:      transaction.UserId,
		Item:        item,
		Unit:        transaction.Unit,
		TotalPrice:  transaction.TotalPrice,
		Status:      transaction.Status,
		CreatedAt:   transaction.CreatedAt,
		UpdatedAt:   transaction.UpdatedAt,
		ExpiredAt:   transaction.ExpiredAt,
	}

	return response
}

