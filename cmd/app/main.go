package main

import (
	"context"
	"log"

	"github.com/Mr-Dryg/car-service-crm/internal/config"
	"github.com/Mr-Dryg/car-service-crm/internal/repository/postgres"
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
}
