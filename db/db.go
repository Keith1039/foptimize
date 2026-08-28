package db

import (
	"context"
	"github.com/Keith1039/foptimize/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DatabaseClient struct {
	db *pgxpool.Pool
}

func NewDatabaseClient(connString string) (DatabaseClient, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return DatabaseClient{db: nil}, err
	}
	return DatabaseClient{db: pool}, nil
}

func (client DatabaseClient) AddUser(user schema.User) {
	args := pgx.NamedArgs{
		"email":  user.Email,
		"config": user.Config,
	}
}

func (client DatabaseClient) GetUserByID(ctx context.Context, id int) schema.User {
	rows, err := client.db.Query(ctx, findUserQueryId, id)
	if err != nil {
		log.Fatal(err)
	}
	// closes row and puts into struct
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.User])
	if err != nil {
		log.Fatalf("User with id '%d' does not exist %v", id, err)
	}
	return user
}

func (client DatabaseClient) GetUserByEmail(ctx context.Context, email string) schema.User {
	rows, err := client.db.Query(ctx, findUserQueryEmail, email)
	if err != nil {
		log.Fatal(err)
	}
	// closes row and puts into struct
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.User])
	if err != nil {
		log.Fatalf("User with email '%s' does not exist %v", email, err)
	}
	return user
}
