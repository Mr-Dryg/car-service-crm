package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Mr-Dryg/car-service-crm/internal/config"
	httphandler "github.com/Mr-Dryg/car-service-crm/internal/handler/http"
	"github.com/Mr-Dryg/car-service-crm/internal/repository/postgres"
	"github.com/Mr-Dryg/car-service-crm/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("loading dotenv with error: %v", err)
	}

	cnf, err := config.Load()
	if err != nil {
		log.Fatalf("loading config with error: %v", err)
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cnf.DatabaseURL)
	if err != nil {
		log.Fatalf("creating pool with error: %v", err)
	}
	defer pool.Close()

	log.Println("success connection to db")

	orderRepo := postgres.NewOrderRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	carRepo := postgres.NewCarRepository(pool)
	branchRepo := postgres.NewBranchRepository(pool)

	orderService := service.NewOrderService(orderRepo, userRepo, carRepo, branchRepo)

	orderHandler := httphandler.NewOrderHandler(orderService)

	router := httphandler.NewRouter(orderHandler)

	log.Println("HTTP-server started on :8080")
	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("HTTP-server start error: %v", err)
	}
}
