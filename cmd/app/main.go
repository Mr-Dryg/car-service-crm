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
	branch := &domain.Branch{
		Name:    "THE COOLEST STO",
		Address: "Moscow",
		Phone:   "+7 (777) 777-77-77",
	}
	err = branchService.Create(ctx, branch)
	if err != nil {
		log.Fatalf("creating branch with error: %v", err)
	}

	userService := service.NewUserService(postgres.NewUserRepository(pool))
	user := &domain.User{
		Name:  "Alex",
		Phone: "+7 (111) 111-11-11",
		Role:  "client",
	}
	err = userService.Create(ctx, user, "")
	if err != nil {
		log.Fatalf("creating user with error: %v", err)
	}

	carService := service.NewCarService(postgres.NewCarRepository(pool))
	car := &domain.Car{
		UserID:       user.ID,
		LicensePlate: "Y777YB777",
		Brand:        "MAZDA",
		Model:        "MX-5",
	}
	err = carService.Create(ctx, car)
	if err != nil {
		log.Fatalf("creating car with error: %v", err)
	}

	cars, err := carService.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Fatalf("getting car with error: %v", err)
	}

	log.Printf("cars for client %v: %v", user.Name, cars)
}
