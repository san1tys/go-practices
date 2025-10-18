package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	connStr := "postgres://postgres4:postgres4@localhost:5433/practice4?sslmode=disable"

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println(" Successfully connected to PostgreSQL database")

	err = applyInitialSchema(db)
	if err != nil {
		log.Fatal("Failed to apply initial schema:", err)
	}

	fmt.Println("\n=== DEMO CRUD OPERATIONS ===")

	fmt.Println("\n1. GetAllUsers:")
	users, err := GetAllUsers(db)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, user := range users {
			fmt.Printf("   ID: %d, Name: %s, Email: %s, Balance: $%.2f\n",
				user.ID, user.Name, user.Email, user.Balance)
		}
	}

	fmt.Println("\n2. GetUserByID (ID: 1):")
	user, err := GetUserByID(db, 1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   User: %s, Balance: $%.2f\n", user.Name, user.Balance)
	}

	fmt.Println("\n3. InsertUser:")
	newUser := User{
		Name:    "Eve Wilson",
		Email:   "eve@example.com",
		Balance: 300.00,
	}
	err = InsertUser(db, newUser)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("   New user inserted successfully")
	}

	fmt.Println("\n4. TransferBalance ($100 from Alice to Bob):")
	err = TransferBalance(db, 1, 2, 100.00)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("   Transfer completed successfully")
	}

	fmt.Println("\n5. Final user balances:")
	users, err = GetAllUsers(db)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		for _, user := range users {
			fmt.Printf("   %s: $%.2f\n", user.Name, user.Balance)
		}
	}

	fmt.Println("\n=== TESTING ERROR CASES ===")

	fmt.Println("6. Transfer with insufficient balance:")
	err = TransferBalance(db, 1, 2, 5000.00)
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}

	fmt.Println("7. Transfer to non-existent user:")
	err = TransferBalance(db, 1, 999, 10.00)
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}
}

func applyInitialSchema(db *sqlx.DB) error {
	var tableExists bool
	err := db.Get(&tableExists,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')")
	if err != nil {
		return err
	}

	if !tableExists {
		schema := `
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				email VARCHAR(100) UNIQUE NOT NULL,
				balance DECIMAL(10,2) DEFAULT 0.00
			);

			INSERT INTO users (name, email, balance) VALUES
			('Alice Johnson', 'alice@example.com', 1000.00),
			('Bob Smith', 'bob@example.com', 500.00),
			('Charlie Brown', 'charlie@example.com', 750.00),
			('Diana Prince', 'diana@example.com', 1200.00);
		`
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
		fmt.Println("Database schema initialized")
	}
	return nil
}

func InsertUser(db *sqlx.DB, user User) error {
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`

	_, err := db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User

	query := `SELECT id, name, email, balance FROM users ORDER BY id`
	err := db.Select(&users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User

	query := `SELECT id, name, email, balance FROM users WHERE id = $1`
	err := db.Get(&user, query, id)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user with ID %d: %w", id, err)
	}

	return user, nil
}

func TransferBalance(db *sqlx.DB, fromID int, toID int, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var fromBalance float64
	err = tx.Get(&fromBalance, "SELECT balance FROM users WHERE id = $1", fromID)
	if err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}

	if fromBalance < amount {
		return fmt.Errorf("insufficient balance: sender has $%.2f, tried to send $%.2f", fromBalance, amount)
	}

	var receiverExists bool
	err = tx.Get(&receiverExists, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", toID)
	if err != nil {
		return fmt.Errorf("failed to check receiver: %w", err)
	}
	if !receiverExists {
		return fmt.Errorf("receiver with ID %d not found", toID)
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return fmt.Errorf("failed to deduct from sender: %w", err)
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return fmt.Errorf("failed to add to receiver: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
