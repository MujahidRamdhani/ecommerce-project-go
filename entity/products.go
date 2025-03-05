package entity

import "time"

type Product struct {
	Id          int
	CatId       int
	Name        string
	Description string
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InputProduct struct {
	Id           int
	CatId        int    `json:"cat_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	IsAvailable  bool   `json:"is_available"`
	StockUnit    int    `json:"stock_unit"`
	PricePerUnit int    `json:"price_per_unit"`
}

type UpdateProduct struct {
	Id           int
	CatId        int    `json:"cat_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsAvailable  bool   `json:"is_available"`
	StockUnit    int    `json:"stock_unit"`
	PricePerUnit int    `json:"price_per_unit"`
}

type Stock struct {
	Id           int
	InvenId      int
	StockUnit    int
	PricePerUnit int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProductStock struct {
	ProductId      int
	CatId          int
	Name           string
	Description    string
	IsAvailable    bool
	ProductCreatedAt time.Time
	ProductUpdatedAt time.Time
	StockId        int
	ProductRefId   int
	StockUnit      int
	PricePerUnit   int
	StockCreatedAt time.Time
	StockUpdatedAt time.Time
}
