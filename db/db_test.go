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
	"time"
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

func genDeal() schema.Deal {
	// get the time dates
	startDate := gofakeit.Date()
	afterDate := startDate.Add(time.Hour * time.Duration(24*gofakeit.Number(14, 21)))
	return schema.Deal{
		DestinationID:        gofakeit.AirlineAirportIATA(),
		Name:                 gofakeit.City(),
		Country:              gofakeit.Country(),
		Price:                gofakeit.Price(500, 4000),
		AveragePrice:         gofakeit.Price(500, 2000),
		DiscountPercentage:   gofakeit.Price(0, 100),
		FlightLink:           gofakeit.URL(),
		SerpApiFlightLink:    gofakeit.URL(),
		Thumbnail:            gofakeit.URL(),
		StartDate:            startDate.String(),
		EndDate:              afterDate.String(),
		DepartureAirportCode: gofakeit.AirlineAirportIATA(),
		FlightDuration:       gofakeit.Number(0, 2000),
		Stops:                gofakeit.Number(0, 10),
		Airline:              gofakeit.AirlineAirport(),
		AirlineCode:          "",
		Description:          "",
		Highlights:           "",
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
	id, err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	testUser.Id = id
	// grab latest user
	query := `SELECT ID, EMAIL, RELEVANT_AIRPORTS, SUBSCRIBED_COUNTRIES, CONFIG, THRESHOLD, CREATED_AT FROM USERS WHERE ID=$1`
	rows, err := pool.Query(ctx, query, id)
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
	id, err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	testUser.Id = id
	user, err := dbClient.GetUser(ctx, id)
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
	id, err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	testUser.Id = id
	// regenerate some fields for update
	testUser.Email = gofakeit.Email()
	testUser.RelevantAirports = []string{gofakeit.AirlineAirportIATA()}
	testUser.SubscribedCountries = []string{gofakeit.Country()}
	testUser.Config = schema.Config{TravelDuration: "1"}
	testUser.Threshold = gofakeit.Number(0, 100)

	err = dbClient.UpdateUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	user, err := dbClient.GetUser(ctx, id)
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
	id, err := dbClient.AddUser(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	testUser.Id = id
	beforeTotal := getCountForTable("USERS")
	err = dbClient.DeleteUser(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	afterTotal := getCountForTable("USERS")
	if (beforeTotal - afterTotal) != 1 {
		t.Fatalf("more than 1 deletion occured. beforeTotal: %d, afterTotal: %d", beforeTotal, afterTotal)
	}
}

// Deal stuff
func TestDatabaseClient_AddDeal(t *testing.T) {
	testDeal := genDeal()
	link := testDeal.FlightLink
	ctx := context.Background()
	err := dbClient.AddDeal(ctx, testDeal)
	if err != nil {
		t.Fatal(err)
	}
	query := `SELECT
				DESTINATION_ID,
				NAME,
				COUNTRY,
				PRICE,
				AVERAGE_PRICE,
				DISCOUNT_PERCENTAGE,
				FLIGHT_LINK,
				SERP_API_FLIGHT_LINK,
				THUMBNAIL,
				START_DATE,
				END_DATE,
				DEPARTURE_AIRPORT_CODE,
				ARRIVAL_AIRPORT_CODE,
				FLIGHT_DURATION,
				STOPS,
				AIRLINE,
				AIRLINE_CODE,
				DESCRIPTION,
				HIGHLIGHTS
			FROM DEALS
			WHERE FLIGHT_LINK=$1
				`
	rows, err := pool.Query(ctx, query, link)
	if err != nil {
		t.Fatal(err)
	}
	deal, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Deal])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(deal, testDeal) {
		t.Fatalf("deal: %+v\n is not equal to testDeal: %+v", deal, testDeal)
	}

	// dupes should not affect it
	err = dbClient.AddDeal(ctx, testDeal)
	if err != nil {
		t.Fatal(err)
	}
}
