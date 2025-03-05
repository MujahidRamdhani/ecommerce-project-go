package service

import (
	"ecommerce-project-go/entity"
	"ecommerce-project-go/repository"
	"errors"
)

type CategoriesService interface {
	GetAll() ([]entity.Category, error)
	InsertCategory(cat entity.Category, isAdmin bool) (entity.Category, error)
	EditCategory(cat entity.Category, id int, isAdmin bool) (entity.Category, error)
	DeleteCategory(id int, isAdmin bool) error
	GetAllProduct(id int) ([]entity.Product, error)
}

type categoriesService struct {
	catRepository repository.CategoriesRepository
}

func NewCategoryService(categoriesRepository repository.CategoriesRepository) *categoriesService {
	return &categoriesService{categoriesRepository}
}

func (s *categoriesService) GetAll() ([]entity.Category, error) {
	categories, err := s.catRepository.GetAll()
	if err != nil {
		return categories, err
	}

	return categories, nil
}

func (s *categoriesService) InsertCategory(cat entity.Category, isAdmin bool) (entity.Category, error) {
	if !isAdmin {
		return entity.Category{}, errors.New("you're not authorized")
	}

	catFound, found, err := s.catRepository.FindByName(cat.Name)
	if err != nil || found {
		return catFound, errors.New("category with such name already exist")
	}

	cat, err = s.catRepository.Save(cat)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func (s *categoriesService) EditCategory(cat entity.Category, id int, isAdmin bool) (entity.Category, error) {
	if !isAdmin {
		return entity.Category{}, errors.New("you're not authorized")
	}

	existingCat, err := s.catRepository.FindById(id)
	if err != nil {
		return entity.Category{}, errors.New("category with such id not found")
	}

	catFound, found, err := s.catRepository.FindByName(cat.Name)
	if err != nil {
		return entity.Category{}, err
	}

	
	if found && catFound.Id != existingCat.Id {
		return catFound, errors.New("category with such name already exists")
	}

	cat.Id = id

	updatedCat, err := s.catRepository.Update(cat)
	if err != nil {
		return updatedCat, err
	}

	return updatedCat, nil
}


func (s *categoriesService) DeleteCategory(id int, isAdmin bool) error {
	if !isAdmin {
		return errors.New("you're not authorized")
	}

	cat, err := s.catRepository.FindById(id)
	if err != nil {
		return errors.New("category with such id not found")
	}

	err = s.catRepository.Delete(cat)
	if err != nil {
		return err
	}

	return nil
}

func (s *categoriesService) GetAllProduct(id int) ([]entity.Product, error) {
	_, err := s.catRepository.FindById(id)
	if err != nil {
		return nil, errors.New("category with such id not found")
	}
	Products, err := s.catRepository.GetAllProduct(id)
	if err != nil {
		return Products, err
	}

	return Products, nil
}
