package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ryangel/ryangel-backend/internal/models"
)

// ErrNotFound represents missing records.
var ErrNotFound = errors.New("record not found")

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) GetByIdentifier(ctx context.Context, identifier string) (*models.Admin, error) {
	identifier = strings.ToLower(identifier)
	const query = `
		SELECT admin_id, username, email, password_hash, is_active, last_login, token, token_expiry, created_at, updated_at
		FROM admin
		WHERE LOWER(username) = $1 OR LOWER(email) = $1
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, identifier)
	return scanAdmin(row)
}

func (r *AdminRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.Admin, error) {
	const query = `
		SELECT admin_id, username, email, password_hash, is_active, last_login, token, token_expiry, created_at, updated_at
		FROM admin
		WHERE token = $1 AND token_expiry > NOW()
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, tokenHash)
	return scanAdmin(row)
}

func (r *AdminRepository) UpdateToken(ctx context.Context, adminID int64, tokenHash string, expiry time.Time) error {
	const stmt = `UPDATE admin SET token = $1, token_expiry = $2, last_login = NOW() WHERE admin_id = $3`
	cmd, err := r.db.Exec(ctx, stmt, tokenHash, expiry, adminID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AdminRepository) ClearToken(ctx context.Context, adminID int64) error {
	const stmt = `UPDATE admin SET token = NULL, token_expiry = NULL WHERE admin_id = $1`
	cmd, err := r.db.Exec(ctx, stmt, adminID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func scanAdmin(row pgx.Row) (*models.Admin, error) {
	var admin models.Admin
	var lastLogin *time.Time
	var tokenHash *string
	var tokenExpiry *time.Time

	if err := row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.Email,
		&admin.PasswordHash,
		&admin.IsActive,
		&lastLogin,
		&tokenHash,
		&tokenExpiry,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	admin.LastLogin = lastLogin
	admin.TokenHash = tokenHash
	admin.TokenExpiry = tokenExpiry
	return &admin, nil
}
