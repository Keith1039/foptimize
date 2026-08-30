package db_test

import (
	"context"
	"github.com/Keith1039/foptimize/config"
	"github.com/Keith1039/foptimize/db"
	"github.com/Keith1039/foptimize/schema"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"reflect"
	"testing"
)

var dbClient db.DatabaseClient

var testUser = schema.User{
	Email:               gofakeit.Email(),
	RelevantAirports:    []string{"YYW"},
	SubscribedCountries: []string{"China"},
	Config:              schema.Config{TravelDuration: "2"},
	Threshold:           50,
}

var pool *pgxpool.Pool

func init() {
	var err error
	dbClient, err = db.NewDatabaseClient(config.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	pool, err = pgxpool.New(context.Background(), config.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
}

func TestDatabaseClient_AddUser(t *testing.T) {
	err := dbClient.AddUser(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	// grab latest user
	query := `SELECT EMAIL, RELEVANT_AIRPORTS, SUBSCRIBED_COUNTRIES, CONFIG, THRESHOLD FROM USERS ORDER BY ID DESC LIMIT 1`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.User])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(user, testUser) {
		t.Fatalf("user: %+v is not equal to test user: %+v", user, testUser)
	}
}
