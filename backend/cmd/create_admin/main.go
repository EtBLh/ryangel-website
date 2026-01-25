package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ryangel/ryangel-backend/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cmd/create_admin/main.go <username> <password> [email]")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]
	email := fmt.Sprintf("%s@admin.local", username)
	if len(os.Args) > 3 {
		email = os.Args[3]
	}

	// Load .env explicitly if it exists in current dir
	_ = godotenv.Load()

	// Load config
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Build DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	// Connect to DB
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Failed to hash password: %v\n", err)
		os.Exit(1)
	}
	passwordHash := string(hashedBytes)

	// Insert
	query := `
		INSERT INTO admin (username, email, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, true, NOW(), NOW())
		RETURNING admin_id
	`
	var id int64
	err = pool.QueryRow(context.Background(), query, username, email, passwordHash).Scan(&id)
	if err != nil {
		fmt.Printf("Failed to insert admin: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Admin user '%s' created successfully with ID: %d\n", username, id)
}
