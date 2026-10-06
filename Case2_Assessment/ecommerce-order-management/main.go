package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"ecommerce-order-management/config"
	"ecommerce-order-management/controller"
	"ecommerce-order-management/repository"
	"ecommerce-order-management/service"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Database Connection
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection check failed: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL database!")

	// 3. Dependency Injection (Repository -> Service -> Controller)[cite: 26, 27]
	repo := repository.NewPostgresRepository(db)
	orderService := service.NewOrderService(repo)
	ctrl := controller.NewConsoleController(orderService)

	// 4. Start Application
	ctrl.Start()
}