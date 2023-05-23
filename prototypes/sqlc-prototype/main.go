package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/humanbeeng/lepo/server/prototypes/sqlc-prototype/database"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello from sqlc prototype")

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Unable to load .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	log.Println("Successfully connected to PlanetScale!")

	queries := database.New(db)

	author, err := queries.GetAuthor(ctx, 1)

	log.Printf("Author %v", author)

	result, err := queries.CreateAuthor(ctx, database.CreateAuthorParams{
		Name: "Ullas",
		Bio:  sql.NullString{String: "asdfasdfl", Valid: true},
	})

	id, _ := result.LastInsertId()

	log.Printf("last inserted id %v", id)
}
