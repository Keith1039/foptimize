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

// User section

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

func (client DatabaseClient) UpdateUser(ctx context.Context, user schema.User) error {
	args := pgx.NamedArgs{
		"EMAIL":                user.Email,
		"RELEVANT_AIRPORTS":    user.RelevantAirports,
		"SUBSCRIBED_COUNTRIES": user.SubscribedCountries,
		"CONFIG":               user.Config,
		"THRESHOLD":            user.Threshold,
	}
	_, err := client.db.Exec(ctx, updateUserQuery, args)
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

// will need cascade logic
func (client DatabaseClient) DeleteUserById(ctx context.Context, id int) error {
	_, err := client.db.Exec(ctx, deleteUserQueryID, id)
	if err != nil {
		return err
	}
	return nil
}

func (client DatabaseClient) DeleteUserByEmail(ctx context.Context, email string) error {
	_, err := client.db.Exec(ctx, deleteUserQueryEmail, email)
	if err != nil {
		return err
	}
	return nil
}

// Deals section

func (client DatabaseClient) AddDeal(ctx context.Context, deal schema.Deal) error {
	args := pgx.NamedArgs{
		"DESTINATION_ID":         deal.DestinationID,
		"NAME":                   deal.Name,
		"COUNTRY":                deal.Country,
		"PRICE":                  deal.Price,
		"AVERAGE_PRICE":          deal.AveragePrice,
		"DISCOUNT_PERCENTAGE":    deal.DiscountPercentage,
		"FLIGHT_LINK":            deal.FlightLink,
		"SERP_API_FLIGHT_LINK":   deal.SerpApiFlightLink,
		"THUMBNAIL":              deal.Thumbnail,
		"START_DATE":             deal.StartDate,
		"END_DATE":               deal.EndDate,
		"DEPARTURE_AIRPORT_CODE": deal.DepartureAirportCode,
		"ARRIVAL_AIRPORT_CODE":   deal.ArrivalAirportCode,
		"FLIGHT_DURATION":        deal.FlightDuration,
		"STOPS":                  deal.Stops,
		"AIRLINE":                deal.Airline,
		"AIRLINE_CODE":           deal.AirlineCode,
		"DESCRIPTION":            deal.Description,
		"HIGHLIGHTS":             deal.Highlights,
	}
	_, err := client.db.Exec(ctx, addDealQuery, args)
	if err != nil {
		return err
	}
	return nil
}

func (client DatabaseClient) GetDeals(ctx context.Context) ([]schema.Deal, error) {
	rows, err := client.db.Query(ctx, getDealsQuery)
	if err != nil {
		return []schema.Deal{}, err
	}
	deals, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Deal])
	if err != nil {
		return []schema.Deal{}, err
	}
	return deals, nil
}
