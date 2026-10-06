package service

import (
	"context"
	"database/sql"
	"fmt"

	"ecommerce-order-management/models"
)

type Repository interface {
	AddProduct(ctx context.Context, p *models.Product) error
	ListProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int64) (*models.Product, error)
	GetProductByName(ctx context.Context, name string) (*models.Product, error)
	UpdateProduct(ctx context.Context, id int64, price float64, stock int) error

	CreateOrder(ctx context.Context, o *models.Order) error
	ListOrders(ctx context.Context) ([]models.Order, error)
	GetOrderByID(ctx context.Context, id int64) (*models.Order, error)
	GetPendingOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
	CancelOrderAndRestoreStock(ctx context.Context, orderID int64) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) AddProduct(ctx context.Context, p *models.Product) error {
	query := `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.Name, p.Price, p.Stock).Scan(&p.ID)
}

func (r *PostgresRepository) ListProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, price, stock FROM products ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return products, nil
}

func (r *PostgresRepository) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	query := `SELECT id, name, price, stock FROM products WHERE id = $1`
	p := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresRepository) GetProductByName(ctx context.Context, name string) (*models.Product, error) {
	query := `SELECT id, name, price, stock FROM products WHERE LOWER(name) = LOWER($1)`
	p := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostgresRepository) UpdateProduct(ctx context.Context, id int64, price float64, stock int) error {
	query := `UPDATE products SET price = $1, stock = $2 WHERE id = $3`
	res, err := r.db.ExecContext(ctx, query, price, stock, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("product not found or no changes made")
	}
	return nil
}

func (r *PostgresRepository) CreateOrder(ctx context.Context, o *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	updateStockQuery := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
	res, err := tx.ExecContext(ctx, updateStockQuery, o.Quantity, o.ProductID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("insufficient stock for product ID %d", o.ProductID)
	}

	orderQuery := `
		INSERT INTO orders (product_id, product_name, quantity, unit_price, total_amount, status)
		VALUES ($1, $2, $3, $4, $5, 'PENDING') RETURNING id, created_at`
	err = tx.QueryRowContext(ctx, orderQuery, o.ProductID, o.ProductName, o.Quantity, o.UnitPrice, o.TotalAmount).Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) ListOrders(ctx context.Context) ([]models.Order, error) {
	query := `SELECT id, product_id, product_name, quantity, unit_price, total_amount, status, created_at FROM orders ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.ProductName, &o.Quantity, &o.UnitPrice, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *PostgresRepository) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	query := `SELECT id, product_id, product_name, quantity, unit_price, total_amount, status, created_at FROM orders WHERE id = $1`
	o := &models.Order{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&o.ID, &o.ProductID, &o.ProductName, &o.Quantity, &o.UnitPrice, &o.TotalAmount, &o.Status, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (r *PostgresRepository) GetPendingOrders(ctx context.Context) ([]models.Order, error) {
	query := `SELECT id, product_id, product_name, quantity, unit_price, total_amount, status, created_at FROM orders WHERE status = 'PENDING'`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.ProductName, &o.Quantity, &o.UnitPrice, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *PostgresRepository) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, orderID)
	return err
}

func (r *PostgresRepository) CancelOrderAndRestoreStock(ctx context.Context, orderID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var o models.Order
	queryOrder := `SELECT id, product_id, quantity, status FROM orders WHERE id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, queryOrder, orderID).Scan(&o.ID, &o.ProductID, &o.Quantity, &o.Status)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	if o.Status != "PENDING" {
		return fmt.Errorf("only PENDING orders can be cancelled")
	}

	updateOrder := `UPDATE orders SET status = 'CANCELLED' WHERE id = $1`
	if _, err := tx.ExecContext(ctx, updateOrder, orderID); err != nil {
		return err
	}

	restoreStock := `UPDATE products SET stock = stock + $1 WHERE id = $2`
	if _, err := tx.ExecContext(ctx, restoreStock, o.Quantity, o.ProductID); err != nil {
		return err
	}

	return tx.Commit()
}