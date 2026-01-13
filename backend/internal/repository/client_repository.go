package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ryangel/ryangel-backend/internal/models"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) GetByPhone(ctx context.Context, phone string) (*models.Client, error) {
	const query = `
		SELECT client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at
		FROM client
		WHERE phone = $1
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, phone)
	return scanClient(row)
}

func (r *ClientRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.Client, error) {
	const query = `
		SELECT client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at
		FROM client
		WHERE google_id = $1
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, googleID)
	return scanClient(row)
}

func (r *ClientRepository) GetByIdentifier(ctx context.Context, identifier string) (*models.Client, error) {
	// Legacy support or alias to GetByPhone if identifier is phone-like
	return r.GetByPhone(ctx, identifier)
}

func (r *ClientRepository) GetByUsername(ctx context.Context, username string) (*models.Client, error) {
	const query = `
		SELECT client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at
		FROM client
		WHERE username = $1
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, username)
	return scanClient(row)
}

func (r *ClientRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.Client, error) {
	const query = `
		SELECT client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at
		FROM client
		WHERE token = $1 AND token_expiry > NOW()
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, tokenHash)
	return scanClient(row)
}

func (r *ClientRepository) UpdateToken(ctx context.Context, clientID int64, tokenHash string, expiry time.Time) error {
	const stmt = `UPDATE client SET token = $1, token_expiry = $2 WHERE client_id = $3`
	cmd, err := r.db.Exec(ctx, stmt, tokenHash, expiry, clientID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ClientRepository) ClearToken(ctx context.Context, clientID int64) error {
	const stmt = `UPDATE client SET token = NULL, token_expiry = NULL WHERE client_id = $1`
	cmd, err := r.db.Exec(ctx, stmt, clientID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ClientRepository) GetByEmail(ctx context.Context, email string) (*models.Client, error) {
	const query = `
		SELECT client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at
		FROM client
		WHERE email = $1
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, email)
	return scanClient(row)
}

func (r *ClientRepository) UpdateClient(ctx context.Context, clientID int64, email *string, username *string, dateOfBirth *time.Time) (*models.Client, error) {
	const query = `
		UPDATE client 
		SET email = COALESCE($2, email),
		    username = COALESCE($3, username),
			date_of_birth = COALESCE($4, date_of_birth),
			updated_at = NOW()
		WHERE client_id = $1
		RETURNING client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at
	`
	row := r.db.QueryRow(ctx, query, clientID, email, username, dateOfBirth)
	return scanClient(row)
}

func (r *ClientRepository) CreateClient(ctx context.Context, phone *string, email *string, username *string, passwordHash *string, googleID *string, isActive bool) (*models.Client, error) {
	const query = `
		INSERT INTO client (phone, email, username, password_hash, google_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, phone, email, username, passwordHash, googleID, isActive)
	return scanClient(row)
}

func (r *ClientRepository) ActivateClient(ctx context.Context, clientID int64) error {
	const stmt = `UPDATE client SET is_active = TRUE, otp_code = NULL, otp_code_expiry = NULL, updated_at = NOW() WHERE client_id = $1`
	cmd, err := r.db.Exec(ctx, stmt, clientID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ClientRepository) UpdateClientRegistration(ctx context.Context, clientID int64, username *string, passwordHash *string) (*models.Client, error) {
	const query = `
		UPDATE client 
		SET username = $2, password_hash = $3, updated_at = NOW()
		WHERE client_id = $1
		RETURNING client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, clientID, username, passwordHash)
	return scanClient(row)
}

func (r *ClientRepository) UpdateGoogleID(ctx context.Context, clientID int64, googleID string) error {
	const stmt = `UPDATE client SET google_id = $1, updated_at = NOW() WHERE client_id = $2`
	cmd, err := r.db.Exec(ctx, stmt, googleID, clientID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ClientRepository) SetOTP(ctx context.Context, phone string, otpCode string, expiry time.Time) error {
	const stmt = `UPDATE client SET otp_code = $1, otp_code_expiry = $2 WHERE phone = $3`
	cmd, err := r.db.Exec(ctx, stmt, otpCode, expiry, phone)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *ClientRepository) VerifyOTP(ctx context.Context, phone string, otpCode string) (*models.Client, error) {
	const query = `
		SELECT client_id, email, username, phone, google_id, password_hash, date_of_birth, is_active, token, token_expiry, otp_code, otp_code_expiry, created_at, updated_at
		FROM client
		WHERE phone = $1 AND otp_code = $2 AND otp_code_expiry > NOW()
		LIMIT 1`

	row := r.db.QueryRow(ctx, query, phone, otpCode)
	client, err := scanClient(row)
	if err != nil {
		return nil, err
	}
	
	// Clear OTP after successful verification
	_, _ = r.db.Exec(ctx, "UPDATE client SET otp_code = NULL, otp_code_expiry = NULL WHERE client_id = $1", client.ID)
	
	return client, nil
}

func scanClient(row pgx.Row) (*models.Client, error) {
	var client models.Client
	var email *string
	var username *string
	var phone *string
	var googleID *string
	var passwordHash *string
	var dateOfBirth *time.Time
	var tokenHash *string
	var tokenExpiry *time.Time
	var otpCode *string
	var otpCodeExpiry *time.Time

	if err := row.Scan(
		&client.ID,
		&email,
		&username,
		&phone,
		&googleID,
		&passwordHash,
		&dateOfBirth,
		&client.IsActive,
		&tokenHash,
		&tokenExpiry,
		&otpCode,
		&otpCodeExpiry,
		&client.CreatedAt,
		&client.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	client.Email = email
	client.Username = username
	client.Phone = phone
	client.GoogleID = googleID
	client.PasswordHash = passwordHash
	client.DateOfBirth = dateOfBirth
	client.TokenHash = tokenHash
	client.TokenExpiry = tokenExpiry
	client.OTPCode = otpCode
	client.OTPCodeExpiry = otpCodeExpiry
	return &client, nil
}
