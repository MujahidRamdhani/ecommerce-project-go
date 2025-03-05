package service

import (
	"ecommerce-project-go/entity"
	"ecommerce-project-go/repository"
	"errors"
)

type ProductService interface {
	AddItem(inputProduct entity.InputProduct, isAdmin bool) (entity.Product, entity.Stock, error)
	UpdateItem(input entity.UpdateProduct, isAdmin bool, id int) (entity.Product, entity.Stock, error)
	DeleteItem(isAdmin bool, id int) error
	GetAll(page int, limit int) ([]entity.InputProduct, map[string]interface{}, error)
	GetById(id int) (entity.InputProduct, error)
}

type productService struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) *productService {
	return &productService{repository}
}

func (s *productService) AddItem(input entity.InputProduct, isAdmin bool) (entity.Product, entity.Stock, error) {
	if !isAdmin {
		return entity.Product{}, entity.Stock{}, errors.New("you're not authorize")
	}

	_, found, _ := s.repository.FindByName(input.Name)
	if found {
		return entity.Product{}, entity.Stock{}, errors.New("name already exist")
	}

	Product, stock, err := s.repository.Save(input)
	if err != nil {
		return Product, stock, err
	}

	return Product, stock, nil
}

func (s *productService) UpdateItem(input entity.UpdateProduct, isAdmin bool, id int) (entity.Product, entity.Stock, error) {
	if !isAdmin {
		return entity.Product{}, entity.Stock{}, errors.New("you're not authorized")
	}

	realInput, idExist, _ := s.repository.FindById(id)
	if !idExist {
		return entity.Product{}, entity.Stock{}, errors.New("product id not found")
	}

	if input.CatId != realInput.CatId {
		realInput.CatId = input.CatId
	}

	if input.Name != "" {
		existingProduct, found, _ := s.repository.FindByName(input.Name)
		if found && existingProduct.Id != realInput.ProductId {
			return entity.Product{}, entity.Stock{}, errors.New("name already exist")
		}
		realInput.Name = input.Name
	}

	if input.Description != "" {
		realInput.Description = input.Description
	}

	if realInput.IsAvailable != input.IsAvailable {
		realInput.IsAvailable = input.IsAvailable
	}

	if input.StockUnit != realInput.StockUnit {
		realInput.StockUnit = input.StockUnit
	}

	if input.PricePerUnit != realInput.PricePerUnit {
		realInput.PricePerUnit = input.PricePerUnit
	}

	Product, stock, err := s.repository.Update(realInput)
	if err != nil {
		return Product, stock, err
	}

	return Product, stock, nil
}

func (s *productService) DeleteItem(isAdmin bool, id int) error {
	var Product entity.Product

	if !isAdmin {
		return errors.New("you're not authorized")
	}

	_, idExist, _ := s.repository.FindById(id)
	if !idExist {
		return errors.New("product id not found")
	}

	Product.Id = id

	err := s.repository.Delete(Product)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) GetAll(page int, limit int) ([]entity.InputProduct, map[string]interface{}, error) {
	products, meta, err := s.repository.GetAll(page, limit)
	if err != nil {
		return products, nil, err
	}

	return products, meta, nil
}

func (s *productService) GetById(id int) (entity.InputProduct, error) {
	_, idExist, _ := s.repository.FindById(id)
	if !idExist {
		return entity.InputProduct{}, errors.New("product id not found")
	}

	Product, err := s.repository.GetById(id)
	if err != nil {
		return Product, err
	}

	return Product, nil
}
