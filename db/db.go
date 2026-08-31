package db

import (
	"context"
	"github.com/Keith1039/foptimize/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (client DatabaseClient) AddUser(ctx context.Context, user schema.User) error {
	args := pgx.NamedArgs{
		"EMAIL":                user.Email,
		"RELEVANT_AIRPORTS":    user.RelevantAirports,
		"SUBSCRIBED_COUNTRIES": user.SubscribedCountries,
		"CONFIG":               user.Config,
		"THRESHOLD":            user.Threshold,
	}
	_, err := client.db.Exec(ctx, addUserQuery, args)
	return err
}

func (client DatabaseClient) GetUserByID(ctx context.Context, id int) (schema.User, error) {
	rows, err := client.db.Query(ctx, findUserQueryId, id)
	if err != nil {
		return schema.User{}, err
	}
	// closes row and puts into struct
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.User])
	if err != nil {
		return schema.User{}, err
	}
	return user, nil
}

func (client DatabaseClient) GetUserByEmail(ctx context.Context, email string) (schema.User, error) {
	rows, err := client.db.Query(ctx, findUserQueryEmail, email)
	if err != nil {
		return schema.User{}, err
	}
	// closes row and puts into struct
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.User])
	if err != nil {
		return schema.User{}, err
	}
	return user, nil
}
