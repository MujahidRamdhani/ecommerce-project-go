package entity

import "time"

type Transaction struct {
	Id             int
	UserId         int
	ProductId      int
	Unit           int
	TotalPrice     int
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExpiredAt      time.Time
	StockRetrieved bool
}

type InputTransaction struct {
	Item     string    `json:"item" binding:"required"`
	Unit     int       `json:"unit" binding:"required,gt=1"`
}
