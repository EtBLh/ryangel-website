package main

import (
    "context"
    "fmt"
    "log"

    "github.com/joho/godotenv"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/ryangel/ryangel-backend/internal/config"
)

func main() {
    // Try loading from up one or two dirs if needed, but since we run from backend root...
    // We will assume running from backend/ directory.
    // If running from backend/cmd/migrate_activated/, then ../../.env
    // We will try running `go run cmd/migrate_activated/main.go` from backend/
    
    godotenv.Load(".env") // Assumes running from backend root

    cfg, err := config.FromEnv()
    if err != nil {
        log.Fatalf("Config error (ensure .env is present): %v", err)
    }

    pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
    if err != nil {
        log.Fatalf("Pool connection error: %v", err)
    }
    defer pool.Close()

    // Check if column exists first
    checkSQL := `
    SELECT column_name 
    FROM information_schema.columns 
    WHERE table_name='client' AND column_name='activated';
    `
    var colName string
    err = pool.QueryRow(context.Background(), checkSQL).Scan(&colName)
    if err == nil && colName == "activated" {
         fmt.Println("Column 'activated' already exists. Skipping.")
         return
    }

    fmt.Println("Applying migration...")
    sql := `
    ALTER TABLE client ADD COLUMN activated BOOLEAN DEFAULT FALSE;
    UPDATE client SET activated = TRUE;
    `
    _, err = pool.Exec(context.Background(), sql)
    if err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

    fmt.Println("Migration successful: Added activated column and updated existing users.")
}
