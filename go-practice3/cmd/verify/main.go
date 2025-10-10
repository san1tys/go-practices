package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://adilkanatov@localhost:5432/expense_tracker?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Successfully connected to database")

	tables := []string{"users", "categories", "expenses"}
	for _, table := range tables {
		var exists bool
		query := `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)`

		err = db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			log.Printf("Failed to check table %s: %v", table, err)
			continue
		}

		if exists {
			fmt.Printf("Table %s exists\n", table)

			var count int
			countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
			err = db.QueryRow(countQuery).Scan(&count)
			if err != nil {
				log.Printf("Could not count rows in %s: %v", table, err)
			} else {
				fmt.Printf("   Rows in %s: %d\n", table, count)
			}
		} else {
			fmt.Printf("Table %s does not exist\n", table)
		}
	}

	fmt.Println("\nChecking users table structure:")
	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns 
		WHERE table_name = 'users' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Printf("Failed to check users structure: %v", err)
	} else {
		defer rows.Close()

		for rows.Next() {
			var columnName, dataType, nullable string
			rows.Scan(&columnName, &dataType, &nullable)
			fmt.Printf("   %s: %s (%s)\n", columnName, dataType, nullable)
		}
	}

	var version int
	var dirty bool
	err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
	if err != nil {
		log.Printf("Could not check migration version: %v", err)
	} else {
		fmt.Printf("\nMigration version: %d, dirty: %t\n", version, dirty)

		if dirty {
			fmt.Println("Database is in dirty state - some migration failed")
			os.Exit(1)
		}
	}

	fmt.Println("\nAll migrations applied successfully! Database schema is ready.")
}
