package main

import (
	"context"
	"os"

	"github.com/0xstxrless/punkt/backend/internal/db"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	ctx := context.Background()

	pool, err := db.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	queries := db.New(pool)

}
