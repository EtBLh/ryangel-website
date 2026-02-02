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
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run cmd/create_client/main.go <username> <password> <phone> [email]")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]
	phone := os.Args[3]
	email := ""
	if len(os.Args) > 4 {
		email = os.Args[4]
	}

	// Load .env explicitly if it exists in current dir or parent dirs
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

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

	// Insert Client
	var query string
	var args []interface{}

	if email != "" {
		query = `INSERT INTO client (username, password_hash, phone, email, is_active, activated) VALUES ($1, $2, $3, $4, true, true)`
		args = []interface{}{username, passwordHash, phone, email}
	} else {
		query = `INSERT INTO client (username, password_hash, phone, is_active, activated) VALUES ($1, $2, $3, true, true)`
		args = []interface{}{username, passwordHash, phone}
	}

	_, err = pool.Exec(context.Background(), query, args...)
	if err != nil {
		fmt.Printf("Failed to insert client: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Client '%s' created successfully.\n", username)
}
