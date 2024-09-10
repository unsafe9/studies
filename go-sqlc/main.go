package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/unsafe9/studies/go-sqlc/db"
	"log"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://admin:admin@localhost:5432/database")
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	ctx = WithDB(ctx, pool)
	doTesting(ctx)
}

func doTesting(ctx context.Context) {
	var id int64
	txErr := Transaction(ctx, func(ctx context.Context) error {
		queries := db.New(DB(ctx))

		_ = OnCommitSuccess(ctx, func(ctx context.Context) error {
			log.Println("transaction committed")
			return nil
		})

		if _, err := queries.AddGreeting(ctx, "Hello, World!"); err != nil {
			return err
		}

		var err error
		id, err = queries.AddGreeting(ctx, "This will be a return value")
		return err
	})
	if txErr != nil {
		log.Fatalf("transaction failed: %v", txErr)
	}

	queries := db.New(DB(ctx))
	greeting, err := queries.GetGreeting(ctx, id)
	if err != nil {
		log.Fatalf("failed to get greeting: %v", err)
	}

	log.Println(greeting.Content)
}
