package postgres

import (
	"context"
	"time"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: pool,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order, parsedDate time.Time) error {
	query := `INSERT INTO orders (branch_id, car_id, service_type, status,
			  preferred_date, preferred_time, price, client_confirmed, notes)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  RETURNING id, created_at, updated_at`
	var branchID *int64
	if order.BranchID > 0 {
		branchID = &order.BranchID
	}
	var notes *string
	if order.Notes != "" {
		notes = &order.Notes
	}
	err := r.db.QueryRow(
		ctx, query, branchID, order.CarID, order.ServiceType,
		order.Status, parsedDate, order.PreferredTime, order.Price,
		order.ClientConfirmed, notes,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	return err
}

func (r *OrderRepository) GetByOrderID(ctx context.Context, orderID int64) (*domain.Order, error) {
	query := `SELECT id, branch_id, car_id, service_type,
			  status, preferred_date, preferred_time, price,
			  client_confirmed, notes, created_at, updated_at
			  FROM orders WHERE id = $1`
	row := r.db.QueryRow(ctx, query, orderID)
	return scanOrder(row)
}

func (r *OrderRepository) GetByBranchID(ctx context.Context, branchID int64, includeCommon bool) ([]domain.Order, error) {
	query := `SELECT id, branch_id, car_id, service_type,
			  status, preferred_date, preferred_time, price,
			  client_confirmed, notes, created_at, updated_at
			  FROM orders WHERE branch_id = $1`
	if includeCommon {
		query += " OR branch_id IS NULL"
	}
	rows, err := r.db.Query(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) GetByCarID(ctx context.Context, carID int64) ([]domain.Order, error) {
	query := `SELECT id, branch_id, car_id, service_type,
			  status, preferred_date, preferred_time, price,
			  client_confirmed, notes, created_at, updated_at
			  FROM orders WHERE car_id = $1`
	rows, err := r.db.Query(ctx, query, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, orderID)
	return err
}

func (r *OrderRepository) UpdateSchedule(ctx context.Context, orderID int64, prefDate time.Time, prefTime string) error {
	query := `UPDATE orders SET preferred_date = $1, preferred_time = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, prefDate, prefTime, orderID)
	return err
}

func (r *OrderRepository) UpdatePrice(ctx context.Context, orderID int64, price float64) error {
	query := `UPDATE orders SET price = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, price, orderID)
	return err
}

func (r *OrderRepository) UpdateClientConfirmed(ctx context.Context, orderID int64, flag bool) error {
	query := `UPDATE orders SET client_confirmed = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, flag, orderID)
	return err
}

func (r *OrderRepository) UpdateNotes(ctx context.Context, orderID int64, notes string) error {
	query := `UPDATE orders SET notes = $1, updated_at = NOW() WHERE id = $2`
	var n *string
	if notes != "" {
		n = &notes
	}
	_, err := r.db.Exec(ctx, query, n, orderID)
	return err
}

func scanOrder(row rowScanner) (*domain.Order, error) {
	var order domain.Order
	var prefDate, prefTime time.Time
	var branchID *int64
	var notes *string
	err := row.Scan(
		&order.ID, &branchID, &order.CarID, &order.ServiceType,
		&order.Status, &prefDate, &prefTime, &order.Price,
		&order.ClientConfirmed, &notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	order.PreferredDate = prefDate.Format(domain.DateLayout)
	order.PreferredTime = prefTime.Format(domain.TimeLayout)
	if branchID != nil {
		order.BranchID = *branchID
	}
	if notes != nil {
		order.Notes = *notes
	}
	return &order, nil
}
