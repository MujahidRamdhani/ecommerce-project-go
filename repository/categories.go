package repository

import (
	"database/sql"
	"ecommerce-project-go/entity"
	"time"
)

type CategoriesRepository interface {
	GetAll() ([]entity.Category, error)
	Save(cat entity.Category) (entity.Category, error)
	Update(cat entity.Category) (entity.Category, error)
	Delete(cat entity.Category) error
	FindById(id int) (entity.Category, error)
	FindByName(name string) (entity.Category, bool, error)
	GetAllProduct(id int) ([]entity.Product, error)
}

type categoriesRepository struct {
	db *sql.DB
}

func NewCategoriesRepository(db *sql.DB) *categoriesRepository {
	return &categoriesRepository{db}
}

func (r *categoriesRepository) GetAll() ([]entity.Category, error) {
	var result []entity.Category

	sqlStatement := `SELECT * FROM Product_categories`
	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		var cat entity.Category
		err = rows.Scan(&cat.Id, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			return result, err
		}
		result = append(result, cat)
	}

	return result, nil
}

func (r *categoriesRepository) Save(cat entity.Category) (entity.Category, error) {
	sqlStatement := `
	INSERT INTO Product_categories(name) 
	VALUES($1) 
	RETURNING *`

	err := r.db.QueryRow(sqlStatement, cat.Name).Scan(&cat.Id, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func (r *categoriesRepository) Update(cat entity.Category) (entity.Category, error) {
	sqlStatement := `
	UPDATE Product_categories 
	SET name = $1, updated_at = $2
	WHERE id = $3
	RETURNING id, name`

	err := r.db.QueryRow(sqlStatement, cat.Name, time.Now(), cat.Id).Scan(&cat.Id, &cat.Name)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func (r *categoriesRepository) Delete(cat entity.Category) error {
	sqlStatement := `DELETE FROM Product_categories WHERE id = $1`
	err := r.db.QueryRow(sqlStatement, cat.Id).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *categoriesRepository) FindById(id int) (entity.Category, error) {
	var cat entity.Category

	sqlStatement := `SELECT id, name FROM Product_categories WHERE id = $1`
	err := r.db.QueryRow(sqlStatement, id).Scan(&cat.Id, &cat.Name)
	if err != nil {
		return cat, err
	}

	return cat, nil
}

func (r *categoriesRepository) FindByName(name string) (entity.Category, bool, error) {
	var cat entity.Category

	sqlStatement := `SELECT id, name FROM Product_categories WHERE name = $1`
	err := r.db.QueryRow(sqlStatement, name).Scan(&cat.Id, &cat.Name)
	if err != nil {
		return cat, false, nil
	}

	return cat, true, nil
}

func (r *categoriesRepository) GetAllProduct(id int) ([]entity.Product, error) {
	var result []entity.Product

	sqlStatement := `SELECT * FROM products WHERE cat_id = $1`
	rows, err := r.db.Query(sqlStatement, id)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		var Product entity.Product
		err = rows.Scan(
			&Product.Id,
			&Product.CatId,
			&Product.Name,
			&Product.Description,
			&Product.IsAvailable,
			&Product.CreatedAt,
			&Product.UpdatedAt)
		if err != nil {
			return result, err
		}

		result = append(result, Product)
	}

	return result, nil
}
