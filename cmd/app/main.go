package main

import (
	"context"
	"log"

	"github.com/Mr-Dryg/car-service-crm/internal/config"
	"github.com/Mr-Dryg/car-service-crm/internal/domain"
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

	branchService := service.NewBranchService(postgres.NewBranchRepository(pool))
	firstBranch := &domain.Branch{
		Name: "THE COOLEST STO",
		Address: "Moscow",
		Phone: "+7 (777) 777-77-77",
	}
	err = branchService.Create(ctx, firstBranch)
	if err != nil {
		log.Fatalf("creating branch with error: %v", err)
	}

	branches, err := branchService.GetAll(ctx)
	if err != nil {
		log.Fatalf("getting braches with error: %v", err)
	}

	log.Println(branches)
}
