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

	orderService := service.NewOrderService(postgres.NewOrderRepository(pool))
	order := &domain.Order{
		BranchID:      branch.ID,
		CarID:         car.ID,
		ServiceType:   "window tinting",
		PreferredDate: "27-08-2026",
		PreferredTime: "15:00",
	}
	err = orderService.Create(ctx, order)
	if err != nil {
		log.Fatalf("creating order with error: %v", err)
	}

	err = orderService.ChangeStatus(ctx, order.ID, domain.StatusConfirmed)
	if err != nil {
		log.Fatalf("changing status with error: %v", err)
	}

	err = orderService.RescheduleOrder(ctx, order.ID, "30-08-2026", "10:00")
	if err != nil {
		log.Fatalf("rescheduling order with error: %v", err)
	}

	err = orderService.UpdateCost(ctx, order.ID, 300.39)
	if err != nil {
		log.Fatalf("updating cost with error: %v", err)
	}

	err = orderService.ChangeStatus(ctx, order.ID, domain.StatusInProgress)
	if err != nil {
		log.Fatalf("changing status with error: %v", err)
	}

	err = orderService.ChangeStatus(ctx, order.ID, domain.StatusReady)
	if err != nil {
		log.Fatalf("changing status with error: %v", err)
	}

	err = orderService.ConfirmReceipt(ctx, order.ID)
	if err != nil {
		log.Fatalf("confirming receipt with error: %v", err)
	}

	err = orderService.UpdateNotes(ctx, order.ID, "some new notes")
	if err != nil {
		log.Fatalf("confirming receipt with error: %v", err)
	}

	orders, err := orderService.GetBranchOrders(ctx, branch.ID, true)
	if err != nil {
		log.Fatalf("getting orders by branch_id with error: %v", err)
	}

	log.Printf("orders: %+v", orders)
}
