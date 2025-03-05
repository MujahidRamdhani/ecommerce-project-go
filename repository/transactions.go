package repository

import (
	"database/sql"
	"ecommerce-project-go/entity"
	"errors"
	"log"
	"time"
)

type TransactionRepository interface {
	Save(transaction entity.Transaction) (entity.Transaction, error)
	FindProductId(item string) (entity.Product, error)
	FindProductById(ProductID int) (string, error)
	FindProductByName(productName string) (entity.Product, error)
	GetStockAndPrice(ProductId int) (entity.Stock, error)
	UpdateStock(ProductId int, unit int) error
	FindById(transactionId int) (entity.Transaction, error)
	Update(transaction entity.Transaction) (entity.Transaction, error)
	GetAll(userID int, page int, limit int) ([]entity.Transaction,  map[string]interface{}, error)
	GetByStatus(userID int, status string) ([]entity.Transaction, error)
	FindUserId(userID int) error
	UpdateStockRetrieved(transaction entity.Transaction) error
	GetAllAdmin() ([]entity.Transaction, error)
	GetByStatusAdmin(status string) ([]entity.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Save(transaction entity.Transaction) (entity.Transaction, error) {
	sqlStatement := `
	INSERT INTO transactions(user_id, product_id, unit, total_price, status, expired_at)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING *`

	err := r.db.QueryRow(
		sqlStatement,
		transaction.UserId,
		transaction.ProductId,
		transaction.Unit,
		transaction.TotalPrice,
		transaction.Status,
		transaction.ExpiredAt).Scan(
		&transaction.Id,
		&transaction.UserId,
		&transaction.ProductId,
		&transaction.Unit,
		&transaction.TotalPrice,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.StockRetrieved,
		&transaction.ExpiredAt)

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) FindProductId(item string) (entity.Product, error) {
	var Product entity.Product

	sqlStatement := `SELECT id, name FROM products WHERE name = $1`
	err := r.db.QueryRow(sqlStatement, item).Scan(&Product.Id, &Product.Name)
	if err != nil {
		return Product, err
	}

	return Product, nil
}

func (r *transactionRepository) FindProductById(ProductID int) (string, error) {
	var item string

	sqlStatement := `SELECT name FROM products WHERE id = $1`
	err := r.db.QueryRow(sqlStatement, ProductID).Scan(&item)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (r *transactionRepository) FindProductByName(ProductName string) (entity.Product, error) {
	var Product entity.Product
	log.Println("Searching for product:", ProductName)

	sqlStatement := `SELECT id, name FROM products WHERE LOWER(name) = LOWER($1)`
	err := r.db.QueryRow(sqlStatement, ProductName).Scan(&Product.Id, &Product.Name)
	if err != nil {
		return Product, err
	}

	log.Println("product result:", Product)

	return Product, nil
}


func (r *transactionRepository) GetStockAndPrice(ProductId int) (entity.Stock, error) {
	var stock entity.Stock

	sqlStatement := `SELECT stock_unit, price_per_unit FROM Product_stocks WHERE product_id = $1`
	err := r.db.QueryRow(sqlStatement, ProductId).Scan(&stock.StockUnit, &stock.PricePerUnit)
	if err != nil {
		return stock, err
	}

	return stock, nil
}

func (r *transactionRepository) UpdateStock(ProductId int, unit int) error {
	sqlStatement := `UPDATE Product_stocks SET stock_unit = $1, updated_at=$3 WHERE product_id = $2`
	err := r.db.QueryRow(sqlStatement, unit, ProductId, time.Now()).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) FindById(transactionId int) (entity.Transaction, error) {
	var transaction entity.Transaction
	log.Println("transactionId:", transactionId)
	sqlStatement := `SELECT * FROM transactions WHERE id = $1`
	err := r.db.QueryRow(sqlStatement, transactionId).Scan(
		&transaction.Id,
		&transaction.UserId,
		&transaction.ProductId,
		&transaction.Unit,
		&transaction.TotalPrice,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.StockRetrieved,
		&transaction.ExpiredAt)

	if err != nil {
		log.Println("FindById error:", err)
		return transaction, err
	}
	log.Println("transaction:", transaction)
	return transaction, nil
}


func (r *transactionRepository) Update(transaction entity.Transaction) (entity.Transaction, error) {
	sqlStatement := `UPDATE transactions SET status = $1, updated_at=$4 WHERE id = $2 AND user_id = $3 RETURNING *`
	err := r.db.QueryRow(
		sqlStatement,
		transaction.Status,
		transaction.Id,
		transaction.UserId,
		time.Now()).Scan(
		&transaction.Id,
		&transaction.UserId,
		&transaction.ProductId,
		&transaction.Unit,
		&transaction.TotalPrice,
		&transaction.Status,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.StockRetrieved,
		&transaction.ExpiredAt)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) GetAll(userID, page, limit int) ([]entity.Transaction, map[string]interface{}, error) {
    var result []entity.Transaction
    var totalTransactions int

    offset := (page - 1) * limit

    // count total transactions
    err := r.db.QueryRow(`SELECT COUNT(*) FROM transactions WHERE user_id = $1`, userID).Scan(&totalTransactions)
    if err != nil {
        return result, nil, err
    }

    // get transactions with pagination
    sqlStatement := `
    SELECT id, user_id, product_id, unit, total_price, status, created_at, updated_at, stock_retrieved, expired_at
    FROM transactions WHERE user_id = $1
    LIMIT $2 OFFSET $3`

    rows, err := r.db.Query(sqlStatement, userID, limit, offset)
    if err != nil {
        return result, nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var transaction entity.Transaction
        err = rows.Scan(&transaction.Id, &transaction.UserId, &transaction.ProductId, &transaction.Unit, &transaction.TotalPrice, 
            &transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.StockRetrieved, &transaction.ExpiredAt)
        if err != nil {
            return result, nil, err
        }
        result = append(result, transaction)
    }

    totalPages := 1
    if limit > 0 {
        totalPages = (totalTransactions + limit - 1) / limit
    }

    meta := map[string]interface{}{
        "TotalTransactions": totalTransactions,
        "TotalPages":        totalPages,
        "CurrentPage":       page,
        "Limit":             limit,
    }

    return result, meta, nil
}

func (r *transactionRepository) GetByStatus(userID int, status string) ([]entity.Transaction, error) {
	var result []entity.Transaction

	sqlStatement := `SELECT * FROM transactions WHERE user_id = $1 AND status = $2`
	rows, err := r.db.Query(sqlStatement, userID, status)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.Id,
			&transaction.UserId,
			&transaction.ProductId,
			&transaction.Unit,
			&transaction.TotalPrice,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.StockRetrieved,
			&transaction.ExpiredAt)

		if err != nil {
			return result, err
		}

		result = append(result, transaction)
	}

	return result, nil
}

func (r *transactionRepository) FindUserId(userID int) error {
	sqlStatement := `SELECT user_id FROM transactions WHERE user_id = $1`
	err := r.db.QueryRow(sqlStatement, userID)
	if err != nil {
		return errors.New("user doesn't have transactions")
	}

	return nil
}

func (r *transactionRepository) UpdateStockRetrieved(transaction entity.Transaction) error {
	sqlStatement := `UPDATE transactions SET stock_retrieved = $1, updated_at=$3 WHERE id = $2`
	err := r.db.QueryRow(sqlStatement, transaction.StockRetrieved, transaction.Id, time.Now()).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) GetAllAdmin() ([]entity.Transaction, error) {
	var result []entity.Transaction

	sqlStatement := `SELECT * FROM transactions`
	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.Id,
			&transaction.UserId,
			&transaction.ProductId,
			&transaction.Unit,
			&transaction.TotalPrice,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.StockRetrieved,
			&transaction.ExpiredAt)

		if err != nil {
			return result, err
		}

		result = append(result, transaction)
	}

	return result, nil
}

func (r *transactionRepository) GetByStatusAdmin(status string) ([]entity.Transaction, error) {
	var result []entity.Transaction

	sqlStatement := `SELECT * FROM transactions WHERE status = $1`
	rows, err := r.db.Query(sqlStatement, status)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		var transaction entity.Transaction
		err = rows.Scan(&transaction.Id,
			&transaction.UserId,
			&transaction.ProductId,
			&transaction.Unit,
			&transaction.TotalPrice,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.StockRetrieved,
			&transaction.ExpiredAt)

		if err != nil {
			return result, err
		}

		result = append(result, transaction)
	}

	return result, nil
}
