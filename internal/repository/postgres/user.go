package postgres

import (
	"context"

	"github.com/Mr-Dryg/car-service-crm/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	var email *string
	if user.Email != "" {
		email = &user.Email
	}
	var passwordHash *string
	if user.PasswordHash != "" {
		passwordHash = &user.PasswordHash
	}
	var branchID *int64
	if user.BranchID > 0 {
		branchID = &user.BranchID
	}
	var telegramID *int64
	if user.TelegramID > 0 {
		telegramID = &user.TelegramID
	}

	query := `INSERT INTO users
			  (name, phone, email, password_hash, role, telegram_id, branch_id)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)
			  RETURNING id, created_at`
	err := r.db.QueryRow(
		ctx, query, user.Name, user.Phone, email,
		passwordHash, user.Role, telegramID, branchID,
	).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT
			 id, name, phone, email, password_hash,
			 role, telegram_id, branch_id, created_at
			 FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)
	return scanUser(row)
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query := `SELECT
			 id, name, phone, email, password_hash,
			 role, telegram_id, branch_id, created_at
			 FROM users WHERE phone = $1`
	row := r.db.QueryRow(ctx, query, phone)
	return scanUser(row)
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, tgID int64) (*domain.User, error) {
	query := `SELECT
			 id, name, phone, email, password_hash,
			 role, telegram_id, branch_id, created_at
			 FROM users WHERE telegram_id = $1`
	row := r.db.QueryRow(ctx, query, tgID)
	return scanUser(row)
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row rowScanner) (*domain.User, error) {
	var u domain.User
	var email, passwordHash *string
	var telegramID, branchID *int64

	err := row.Scan(
		&u.ID, &u.Name, &u.Phone, &email, &passwordHash,
		&u.Role, &telegramID, &branchID, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if email != nil {
		u.Email = *email
	}
	if passwordHash != nil {
		u.PasswordHash = *passwordHash
	}
	if telegramID != nil {
		u.TelegramID = *telegramID
	}
	if branchID != nil {
		u.BranchID = *branchID
	}

	return &u, nil
}
