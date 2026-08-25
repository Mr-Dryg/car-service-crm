package postgres

import (
	"context"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarRepository struct {
	db *pgxpool.Pool
}

func NewCarRepository(pool *pgxpool.Pool) *CarRepository {
	return &CarRepository{
		db: pool,
	}
}

func (r *CarRepository) Create(ctx context.Context, car *domain.Car) error {
	query := `INSERT INTO cars (user_id, license_plate, brand, model)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at`
	err := r.db.QueryRow(
		ctx, query, car.UserID,
		car.LicensePlate, car.Brand, car.Model,
	).Scan(&car.ID, &car.CreatedAt)
	return err
}

func (r *CarRepository) GetByLicensePlate(ctx context.Context, plate string) (*domain.Car, error) {
	query := `SELECT id, user_id, license_plate, brand, model, created_at
			  FROM cars WHERE license_plate = $1`
	var car domain.Car
	err := r.db.QueryRow(ctx, query, plate).Scan(
		&car.ID, &car.UserID, &car.LicensePlate,
		&car.Brand, &car.Model, &car.CreatedAt,
	)
	return &car, err
}

func (r *CarRepository) GetByUserID(ctx context.Context, userID int64) ([]domain.Car, error) {
	query := `SELECT id, user_id, license_plate, brand, model, created_at
			  FROM cars WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var cars []domain.Car
	for rows.Next() {
		var car domain.Car
		err = rows.Scan(
			&car.ID, &car.UserID, &car.LicensePlate,
			&car.Brand, &car.Model, &car.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, err
}
