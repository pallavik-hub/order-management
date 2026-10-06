package controller

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ecommerce-order-management/service"
)

type ConsoleController struct {
	svc    service.Service
	reader *bufio.Reader
}

func NewConsoleController(svc service.Service) *ConsoleController {
	return &ConsoleController{
		svc:    svc,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (c *ConsoleController) Start() {
	for {
		fmt.Println("\n===========================================")
		fmt.Println("   E-COMMERCE ORDER MANAGEMENT SYSTEM      ")
		fmt.Println("===========================================")
		fmt.Println("1. Add Product")
		fmt.Println("2. List Products")
		fmt.Println("3. Search Product")
		fmt.Println("4. Update Product")
		fmt.Println("5. Create Order")
		fmt.Println("6. List Orders")
		fmt.Println("7. Cancel Order")
		fmt.Println("8. Process Orders Concurrently")
		fmt.Println("9. Exit")
		fmt.Print("Select an option (1-9): ")

		input := c.readLine()
		ctx := context.Background()

		switch input {
		case "1":
			c.handleAddProduct(ctx)
		case "2":
			c.handleListProducts(ctx)
		case "3":
			c.handleSearchProduct(ctx)
		case "4":
			c.handleUpdateProduct(ctx)
		case "5":
			c.handleCreateOrder(ctx)
		case "6":
			c.handleListOrders(ctx)
		case "7":
			c.handleCancelOrder(ctx)
		case "8":
			c.handleProcessConcurrently(ctx)
		case "9":
			fmt.Println("Exiting... Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please choose 1 to 9.")
		}
	}
}

func (c *ConsoleController) handleAddProduct(ctx context.Context) {
	fmt.Print("Enter Product Name: ")
	name := c.readLine()

	fmt.Print("Enter Price: ")
	price, _ := strconv.ParseFloat(c.readLine(), 64)

	fmt.Print("Enter Stock: ")
	stock, _ := strconv.Atoi(c.readLine())

	p, err := c.svc.AddProduct(ctx, name, price, stock)
	if err != nil {
		fmt.Printf("Error adding product: %v\n", err)
		return
	}
	fmt.Printf("Product added successfully! ID: %d\n", p.ID)
}

func (c *ConsoleController) handleListProducts(ctx context.Context) {
	products, err := c.svc.ListProducts(ctx)
	if err != nil {
		fmt.Printf("Error listing products: %v\n", err)
		return
	}
	fmt.Println("\n--- PRODUCT LIST ---")
	fmt.Printf("%-5s %-25s %-10s %-8s\n", "ID", "Name", "Price", "Stock")
	for _, p := range products {
		fmt.Printf("%-5d %-25s %-10.2f %-8d\n", p.ID, p.Name, p.Price, p.Stock)
	}
}

func (c *ConsoleController) handleSearchProduct(ctx context.Context) {
	fmt.Print("Enter Product ID or Name: ")
	query := c.readLine()

	p, err := c.svc.SearchProduct(ctx, query)
	if err != nil || p == nil {
		fmt.Println("Product not found.")
		return
	}
	fmt.Printf("Found Product -> ID: %d | Name: %s | Price: %.2f | Stock: %d\n", p.ID, p.Name, p.Price, p.Stock)
}

func (c *ConsoleController) handleUpdateProduct(ctx context.Context) {
	fmt.Print("Enter Product ID: ")
	id, _ := strconv.ParseInt(c.readLine(), 10, 64)

	fmt.Print("Enter New Price: ")
	price, _ := strconv.ParseFloat(c.readLine(), 64)

	fmt.Print("Enter New Stock: ")
	stock, _ := strconv.Atoi(c.readLine())

	err := c.svc.UpdateProduct(ctx, id, price, stock)
	if err != nil {
		fmt.Printf("Error updating product: %v\n", err)
		return
	}
	fmt.Println("Product updated successfully!")
}

func (c *ConsoleController) handleCreateOrder(ctx context.Context) {
	fmt.Print("Enter Product ID: ")
	id, _ := strconv.ParseInt(c.readLine(), 10, 64)

	fmt.Print("Enter Quantity: ")
	qty, _ := strconv.Atoi(c.readLine())

	order, err := c.svc.CreateOrder(ctx, id, qty)
	if err != nil {
		fmt.Printf("Error creating order: %v\n", err)
		return
	}
	fmt.Printf("Order created successfully! ID: %d | Total: %.2f | Status: %s\n", order.ID, order.TotalAmount, order.Status)
}

func (c *ConsoleController) handleListOrders(ctx context.Context) {
	orders, err := c.svc.ListOrders(ctx)
	if err != nil {
		fmt.Printf("Error listing orders: %v\n", err)
		return
	}
	fmt.Println("\n--- ORDER LIST ---")
	fmt.Printf("%-5s %-20s %-8s %-10s %-12s\n", "ID", "Product", "Qty", "Total", "Status")
	for _, o := range orders {
		fmt.Printf("%-5d %-20s %-8d %-10.2f %-12s\n", o.ID, o.ProductName, o.Quantity, o.TotalAmount, o.Status)
	}
}

func (c *ConsoleController) handleCancelOrder(ctx context.Context) {
	fmt.Print("Enter Order ID to Cancel: ")
	id, _ := strconv.ParseInt(c.readLine(), 10, 64)

	err := c.svc.CancelOrder(ctx, id)
	if err != nil {
		fmt.Printf("Error cancelling order: %v\n", err)
		return
	}
	fmt.Println("Order cancelled successfully and stock restored!")
}

func (c *ConsoleController) handleProcessConcurrently(ctx context.Context) {
	fmt.Print("Enter worker pool size (e.g., 3): ")
	workers, _ := strconv.Atoi(c.readLine())
	if workers <= 0 {
		workers = 3
	}

	err := c.svc.ProcessOrdersConcurrently(ctx, workers)
	if err != nil {
		fmt.Printf("Error processing orders: %v\n", err)
	}
}

func (c *ConsoleController) readLine() string {
	input, _ := c.reader.ReadString('\n')
	return strings.TrimSpace(input)
}