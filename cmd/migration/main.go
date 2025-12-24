package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"app/config"
	_ "github.com/jackc/pgx/v5/stdlib" // Use pgx standard library adapter
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		fmt.Printf("Running migration: %s\n", fileName)
		content, err := ioutil.ReadFile(filepath.Join("migrations", fileName))
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", fileName, err)
		}

		// Split Up/Down (Simplistic approach, executing the whole file for now assuming it's Up)
		// Usually we parse -- Up and -- Down. For init, we can just run the whole thing if we are careful,
		// but better to only run the Up part.
		
		requests := strings.Split(string(content), "-- Down")
		upScript := requests[0]

		if _, err := db.Exec(upScript); err != nil {
			// Ignore "already exists" errors for idempotency if simple scripts, 
			// but for production use a real migration tool like golang-migrate
			log.Printf("Error running migration %s: %v", fileName, err)
		} else {
			fmt.Printf("Migration %s completed successfully.\n", fileName)
		}
	}
}
