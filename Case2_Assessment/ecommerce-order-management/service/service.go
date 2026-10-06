package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"ecommerce-order-management/models"
	"ecommerce-order-management/repository"
)

type Service interface {
	AddProduct(ctx context.Context, name string, price float64, stock int) (*models.Product, error)
	ListProducts(ctx context.Context) ([]models.Product, error)
	SearchProduct(ctx context.Context, query string) (*models.Product, error)
	UpdateProduct(ctx context.Context, id int64, price float64, stock int) error
	CreateOrder(ctx context.Context, productID int64, quantity int) (*models.Order, error)
	ListOrders(ctx context.Context) ([]models.Order, error)
	CancelOrder(ctx context.Context, orderID int64) error
	ProcessOrdersConcurrently(ctx context.Context, workerCount int) error
}

type OrderService struct {
	repo repository.Repository
}

func NewOrderService(repo repository.Repository) Service {
	return &OrderService{repo: repo}
}

func (s *OrderService) AddProduct(ctx context.Context, name string, price float64, stock int) (*models.Product, error) {
	if name == "" || price < 0 || stock < 0 {
		return nil, fmt.Errorf("invalid input parameters")
	}
	p := &models.Product{Name: name, Price: price, Stock: stock}
	if err := s.repo.AddProduct(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *OrderService) ListProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *OrderService) SearchProduct(ctx context.Context, search string) (*models.Product, error) {
	var id int64
	if _, err := fmt.Sscanf(search, "%d", &id); err == nil {
		p, err := s.repo.GetProductByID(ctx, id)
		if err == nil {
			return p, nil
		}
	}
	return s.repo.GetProductByName(ctx, search)
}

func (s *OrderService) UpdateProduct(ctx context.Context, id int64, price float64, stock int) error {
	if price < 0 || stock < 0 {
		return fmt.Errorf("price and stock must be non-negative")
	}
	return s.repo.UpdateProduct(ctx, id, price, stock)
}

func (s *OrderService) CreateOrder(ctx context.Context, productID int64, quantity int) (*models.Order, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}

	p, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if p.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock: requested %d, available %d", quantity, p.Stock)
	}

	order := &models.Order{
		ProductID:   p.ID,
		ProductName: p.Name,
		Quantity:    quantity,
		UnitPrice:   p.Price,
		TotalAmount: p.Price * float64(quantity),
		Status:      "PENDING",
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context) ([]models.Order, error) {
	return s.repo.ListOrders(ctx)
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID int64) error {
	return s.repo.CancelOrderAndRestoreStock(ctx, orderID)
}

func (s *OrderService) ProcessOrdersConcurrently(parentCtx context.Context, workerCount int) error {
	pendingOrders, err := s.repo.GetPendingOrders(parentCtx)
	if err != nil {
		return fmt.Errorf("failed to fetch pending orders: %w", err)
	}

	if len(pendingOrders) == 0 {
		fmt.Println("\n[Info] No pending orders found to process.")
		return nil
	}

	fmt.Printf("\nFound %d pending order(s). Processing with %d workers...\n", len(pendingOrders), workerCount)

	ctx, cancel := context.WithTimeout(parentCtx, 15*time.Second)
	defer cancel()

	orderChan := make(chan models.Order, len(pendingOrders))
	var wg sync.WaitGroup

	worker := func(workerID int) {
		defer wg.Done()
		for order := range orderChan {
			select {
			case <-ctx.Done():
				fmt.Printf("[Worker %d] Task cancelled for Order #%d: %v\n", workerID, order.ID, ctx.Err())
				return
			default:
				fmt.Printf("[Worker %d] Started processing Order #%d (%s)...\n", workerID, order.ID, order.ProductName)
				time.Sleep(300 * time.Millisecond) // Simulating async work

				if err := s.repo.UpdateOrderStatus(ctx, order.ID, "COMPLETED"); err != nil {
					fmt.Printf("[Worker %d] Failed to complete Order #%d: %v\n", workerID, order.ID, err)
				} else {
					fmt.Printf("[Worker %d] Successfully COMPLETED Order #%d\n", workerID, order.ID)
				}
			}
		}
	}

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(w)
	}

	for _, order := range pendingOrders {
		orderChan <- order
	}
	close(orderChan)

	wg.Wait()
	fmt.Println("All pending orders have been processed.")
	return nil
}