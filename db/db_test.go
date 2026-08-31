package db_test

import (
	"context"
	"fmt"
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

var pool *pgxpool.Pool

func genUser() schema.User {
	return schema.User{
		Email:               gofakeit.Email(),
		RelevantAirports:    []string{gofakeit.AirlineAirportIATA()},
		SubscribedCountries: []string{gofakeit.Country()},
		Config:              schema.Config{TravelDuration: "2"},
		Threshold:           gofakeit.Number(0, 100),
	}
}

func getCountForTable(tableName string) int {
	query := fmt.Sprintf("SELECT COUNT(*) as total FROM %s", tableName)
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	rows.Next()
	var total int
	err = rows.Scan(&total)
	if err != nil {
		log.Fatal(err)
	}
	return total
}

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
	testUser := genUser()
	ctx := context.Background()
	err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	// grab latest user
	query := `SELECT EMAIL, RELEVANT_AIRPORTS, SUBSCRIBED_COUNTRIES, CONFIG, THRESHOLD, CREATED_AT FROM USERS ORDER BY ID DESC LIMIT 1`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		t.Fatal(err)
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.User])
	if err != nil {
		t.Fatal(err)
	}
	// account for any slight difference
	testUser.CreatedAt = user.CreatedAt
	if !reflect.DeepEqual(user, testUser) {
		t.Fatalf("user: %+v is not equal to test user: %+v", user, testUser)
	}
}

func TestDatabaseClient_GetUser(t *testing.T) {
	testUser := genUser()
	ctx := context.Background()
	err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	query := `SELECT ID FROM USERS WHERE EMAIL=$1`
	rows, err := pool.Query(ctx, query, testUser.Email)
	if err != nil {
		t.Fatal(err)
	}
	var id int
	check := rows.Next()
	if !check {
		t.Log("no next row!")
	}
	err = rows.Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	user, err := dbClient.GetUserByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	// account for any slight difference
	testUser.CreatedAt = user.CreatedAt
	if !reflect.DeepEqual(user, testUser) {
		t.Fatalf("user: %+v is not equal to test user: %+v", user, testUser)
	}
}

func TestDatabaseClient_GetUserByEmail(t *testing.T) {
	testUser := genUser()
	ctx := context.Background()
	err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	user, err := dbClient.GetUserByEmail(ctx, testUser.Email)
	if err != nil {
		t.Fatal(err)
	}
	// account for any slight difference
	testUser.CreatedAt = user.CreatedAt
	if !reflect.DeepEqual(user, testUser) {
		t.Fatalf("user: %+v is not equal to test user: %+v", user, testUser)
	}
}

func TestDatabaseClient_UpdateUser(t *testing.T) {
	testUser := genUser()
	ctx := context.Background()
	err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	// regenerate some fields for update
	testUser.RelevantAirports = []string{gofakeit.AirlineAirportIATA()}
	testUser.SubscribedCountries = []string{gofakeit.Country()}
	testUser.Config = schema.Config{TravelDuration: "1"}
	testUser.Threshold = gofakeit.Number(0, 100)

	err = dbClient.UpdateUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	user, err := dbClient.GetUserByEmail(ctx, testUser.Email)
	if err != nil {
		t.Fatal(err)
	}
	// account for any slight difference
	testUser.CreatedAt = user.CreatedAt
	if !reflect.DeepEqual(user, testUser) {
		t.Fatalf("user: %+v is not equal to test user: %+v", user, testUser)
	}
}

func TestDatabaseClient_DeleteUserById(t *testing.T) {
	testUser := genUser()
	ctx := context.Background()
	err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	query := `SELECT ID FROM USERS WHERE EMAIL=$1`
	rows, err := pool.Query(ctx, query, testUser.Email)
	if err != nil {
		t.Fatal(err)
	}
	var id int
	check := rows.Next()
	if !check {
		t.Log("no next row!")
	}
	err = rows.Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	beforeTotal := getCountForTable("USERS")
	err = dbClient.DeleteUserById(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	afterTotal := getCountForTable("USERS")
	if (beforeTotal - afterTotal) != 1 {
		t.Fatalf("more than 1 deletion occured. beforeTotal: %d, afterTotal: %d", beforeTotal, afterTotal)
	}
}

func TestDatabaseClient_DeleteUserByEmail(t *testing.T) {
	testUser := genUser()
	ctx := context.Background()
	err := dbClient.AddUser(ctx, testUser)
	beforeTotal := getCountForTable("USERS")
	err = dbClient.DeleteUserByEmail(ctx, testUser.Email)
	if err != nil {
		t.Fatal(err)
	}
	afterTotal := getCountForTable("USERS")
	if (beforeTotal - afterTotal) != 1 {
		t.Fatalf("more than 1 deletion occured. beforeTotal: %d, afterTotal: %d", beforeTotal, afterTotal)
	}
}
