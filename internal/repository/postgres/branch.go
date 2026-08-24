package postgres

import (
	"context"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BranchRepository struct {
	db *pgxpool.Pool
}

func NewBranchRepository(pool *pgxpool.Pool) *BranchRepository {
	return &BranchRepository{
		db: pool,
	}
}

func (r *BranchRepository) Create(ctx context.Context, branch *domain.Branch) error {
	query := `INSERT INTO branches (name, address, phone) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, branch.Name, branch.Address, branch.Phone)
	return err
}

func (r *BranchRepository) GetAll(ctx context.Context) ([]domain.Branch, error) {
	query := `SELECT id, name, address, phone, created_at FROM branches`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []domain.Branch
	for rows.Next() {
		var b domain.Branch
		if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.Phone, &b.CreatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}

	return branches, rows.Err()
}
