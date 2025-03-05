package repository

import (
	"database/sql"
	"ecommerce-project-go/entity"
	"time"
)

type ProductRepository interface {
	Save(inputProduct entity.InputProduct) (entity.Product, entity.Stock, error)
	FindById(id int) (entity.ProductStock, bool, error)
	FindByName(name string) (entity.Product, bool, error)
	Update(inputProduct entity.ProductStock) (entity.Product, entity.Stock, error)
	Delete(Product entity.Product) error
	GetAll(page int, limit int) ([]entity.InputProduct, map[string]interface{}, error)
	GetById(id int) (entity.InputProduct, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{db}
}

func (r *productRepository) Save(inputProduct entity.InputProduct) (entity.Product, entity.Stock, error) {
	var Product entity.Product
	var stock entity.Stock

	sqlStatement := `
	INSERT INTO products(cat_id, name, description, is_available)
	VALUES ($1, $2, $3, $4)
	RETURNING id, cat_id, name, description, is_available`

	sqlStatementStock := `
	INSERT INTO Product_stocks(product_id, stock_unit, price_per_unit)
	VALUES ($1, $2, $3)
	RETURNING stock_unit, price_per_unit`

	err := r.db.QueryRow(
		sqlStatement,
		inputProduct.CatId,
		inputProduct.Name,
		inputProduct.Description,
		inputProduct.IsAvailable).Scan(
		&Product.Id,
		&Product.CatId,
		&Product.Name,
		&Product.Description,
		&Product.IsAvailable)
	if err != nil {
		return Product, stock, err
	}

	err = r.db.QueryRow(
		sqlStatementStock,
		Product.Id,
		inputProduct.StockUnit,
		inputProduct.PricePerUnit).Scan(
		&stock.StockUnit,
		&stock.PricePerUnit)
	if err != nil {
		return Product, stock, err
	}

	return Product, stock, nil
}

func (r *productRepository) FindById(id int) (entity.ProductStock, bool, error) {
	var Product entity.ProductStock

	sqlStatement := `
	SELECT * FROM products
	JOIN Product_stocks ON products.id = Product_stocks.product_id
	WHERE products.id = $1`
	err := r.db.QueryRow(
		sqlStatement,
		id).Scan(
		&Product.ProductId,
		&Product.CatId,
		&Product.Name,
		&Product.Description,
		&Product.IsAvailable,
		&Product.ProductCreatedAt,
		&Product.ProductUpdatedAt,
		&Product.StockId,
		&Product.ProductId,
		&Product.StockUnit,
		&Product.PricePerUnit,
		&Product.StockCreatedAt,
		&Product.StockUpdatedAt)
	if err != nil {
		return Product, false, nil
	}

	return Product, true, nil
}

func (r *productRepository) FindByName(name string) (entity.Product, bool, error) {
	var Product entity.Product

	sqlStatement := `SELECT id, name FROM products WHERE name = $1`
	err := r.db.QueryRow(sqlStatement, name).Scan(&Product.Id, &Product.Name)
	if err != nil {
		return Product, false, nil
	}

	return Product, true, nil
}

func (r *productRepository) Update(input entity.ProductStock) (entity.Product, entity.Stock, error) {
	var Product entity.Product
	var stock entity.Stock

	sqlStatement := `
	UPDATE products 
	SET cat_id=$1, name=$2, description=$3, is_available=$4, updated_at=$5
	WHERE id = $6
	RETURNING id, cat_id, name, description, is_available`

	sqlStatementStock := `
	UPDATE Product_stocks
	SET stock_unit=$1, price_per_unit=$2, updated_at=$3
	WHERE product_id = $4
	RETURNING stock_unit, price_per_unit`

	err := r.db.QueryRow(
		sqlStatement,
		input.CatId,
		input.Name,
		input.Description,
		input.IsAvailable,
		time.Now(),
		input.ProductId).Scan(
		&Product.Id,
		&Product.CatId,
		&Product.Name,
		&Product.Description,
		&Product.IsAvailable)

	if err != nil {
		return Product, stock, err
	}

	err = r.db.QueryRow(
		sqlStatementStock,
		input.StockUnit,
		input.PricePerUnit,
		time.Now(),
		Product.Id).Scan(
		&stock.StockUnit,
		&stock.PricePerUnit)

	if err != nil {
		return Product, stock, err
	}

	return Product, stock, nil
}

func (r *productRepository) Delete(Product entity.Product) error {
	sqlStatement := `DELETE FROM products WHERE id = $1`

	err := r.db.QueryRow(sqlStatement, Product.Id).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) GetAll(page int, limit int) ([]entity.InputProduct, map[string]interface{}, error) {
    var result []entity.InputProduct
    var totalProducts int

    offset := (page - 1) * limit

    // count total produk
    err := r.db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&totalProducts)
    if err != nil {
        return result, nil, err
    }

    // get data produk
    sqlStatement := `
    SELECT i.id, i.cat_id, i.name, i.description, i.is_available, s.stock_unit, s.price_per_unit
    FROM products i 
    JOIN product_stocks s ON i.id = s.product_id
    LIMIT $1 OFFSET $2`

    rows, err := r.db.Query(sqlStatement, limit, offset)
    if err != nil {
        return result, nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var product entity.InputProduct
        err = rows.Scan(
            &product.Id,
            &product.CatId,
            &product.Name,
            &product.Description,
            &product.IsAvailable,
            &product.StockUnit,
            &product.PricePerUnit,
        )
        if err != nil {
            return result, nil, err
        }
        result = append(result, product)
    }

    totalPages := 1
    if limit > 0 {
        totalPages = (totalProducts + limit - 1) / limit
    }

    meta := map[string]interface{}{
        "TotalProducts": totalProducts,
        "TotalPages":    totalPages,
        "CurrentPage":   page,
        "Limit":         limit,
    }

    return result, meta, nil
}


func (r *productRepository) GetById(id int) (entity.InputProduct, error) {
	var Product entity.InputProduct

	sqlStatement := `
	SELECT i.id, i.cat_id, i.name, i.description, i.is_available, s.stock_unit, s.price_per_unit
	FROM products i JOIN Product_stocks s 
	    ON i.id = s.product_id
	    WHERE i.id = $1`

	err := r.db.QueryRow(
		sqlStatement,
		id).Scan(
		&Product.Id,
		&Product.CatId,
		&Product.Name,
		&Product.Description,
		&Product.IsAvailable,
		&Product.StockUnit,
		&Product.PricePerUnit)

	if err != nil {
		return Product, err
	}

	return Product, nil
}
