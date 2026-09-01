package db

import (
	"context"
	"fmt"
	"github.com/Keith1039/foptimize/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"slices"
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

func (client DatabaseClient) AddUser(ctx context.Context, user schema.User) (int, error) {
	args := pgx.NamedArgs{
		"EMAIL":                user.Email,
		"RELEVANT_AIRPORTS":    user.RelevantAirports,
		"SUBSCRIBED_COUNTRIES": user.SubscribedCountries,
		"CONFIG":               user.Config,
		"THRESHOLD":            user.Threshold,
	}
	rows, err := client.db.Query(ctx, addUserQuery, args)
	if err != nil {
		return -1, err
	}
	id, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (client DatabaseClient) UpdateUser(ctx context.Context, user schema.User) error {
	args := pgx.NamedArgs{
		"ID":                   user.Id,
		"EMAIL":                user.Email,
		"RELEVANT_AIRPORTS":    user.RelevantAirports,
		"SUBSCRIBED_COUNTRIES": user.SubscribedCountries,
		"CONFIG":               user.Config,
		"THRESHOLD":            user.Threshold,
	}
	_, err := client.db.Exec(ctx, updateUserQuery, args)
	return err
}

func (client DatabaseClient) GetUser(ctx context.Context, id int) (schema.User, error) {
	rows, err := client.db.Query(ctx, findUserQuery, id)
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
func (client DatabaseClient) DeleteUser(ctx context.Context, id int) error {
	_, err := client.db.Exec(ctx, deleteUserQuery, id)
	if err != nil {
		return err
	}
	return nil
}

// Deals section

func getDealArgs(deal schema.Deal) pgx.NamedArgs {
	return pgx.NamedArgs{
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
}

//func (client DatabaseClient) AddDeal(ctx context.Context, deal schema.Deal) error {
//	args := pgx.NamedArgs{
//		"DESTINATION_ID":         deal.DestinationID,
//		"NAME":                   deal.Name,
//		"COUNTRY":                deal.Country,
//		"PRICE":                  deal.Price,
//		"AVERAGE_PRICE":          deal.AveragePrice,
//		"DISCOUNT_PERCENTAGE":    deal.DiscountPercentage,
//		"FLIGHT_LINK":            deal.FlightLink,
//		"SERP_API_FLIGHT_LINK":   deal.SerpApiFlightLink,
//		"THUMBNAIL":              deal.Thumbnail,
//		"START_DATE":             deal.StartDate,
//		"END_DATE":               deal.EndDate,
//		"DEPARTURE_AIRPORT_CODE": deal.DepartureAirportCode,
//		"ARRIVAL_AIRPORT_CODE":   deal.ArrivalAirportCode,
//		"FLIGHT_DURATION":        deal.FlightDuration,
//		"STOPS":                  deal.Stops,
//		"AIRLINE":                deal.Airline,
//		"AIRLINE_CODE":           deal.AirlineCode,
//		"DESCRIPTION":            deal.Description,
//		"HIGHLIGHTS":             deal.Highlights,
//	}
//	_, err := client.db.Exec(ctx, addDealQuery, args)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (client DatabaseClient) MapUserAndDeal(ctx context.Context, deal schema.Deal, id int) error {
//	args := pgx.NamedArgs{
//		"USER_ID":     id,
//		"FLIGHT_LINK": deal.FlightLink,
//	}
//	_, err := client.db.Exec(ctx, mapUserAndDeal, args)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (client DatabaseClient) SaveDeals(ctx context.Context, id int, deals []schema.Deal) error {
	var rollBackErr error
	user, err := client.GetUser(ctx, id)
	if err != nil {
		return err
	}
	tx, err := client.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	for _, deal := range deals {
		// check if it's a subscribed country... maybe this should be a map for quick access?
		if slices.Contains(user.SubscribedCountries, deal.Country) {
			// add the deal
			args := getDealArgs(deal)
			_, err = tx.Exec(ctx, addDealQuery, args)
			if err != nil {
				// doesn't matter
				rollBackErr = tx.Rollback(ctx)
				fmt.Printf("Rollaback Error occured: %v", rollBackErr)
				return err
			}

			mappingArgs := pgx.NamedArgs{
				"USER_ID":     id,
				"FLIGHT_LINK": deal.FlightLink,
			}
			// map to user
			_, err = tx.Exec(ctx, mapUserAndDeal, mappingArgs)
			if err != nil {
				// doesn't matter
				rollBackErr = tx.Rollback(ctx)
				fmt.Printf("Rollaback Error occured: %v", rollBackErr)
				return err
			}
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (client DatabaseClient) GetUserDeals(ctx context.Context, id int) ([]schema.Deal, error) {
	rows, err := client.db.Query(ctx, getDealsQuery, id)
	if err != nil {
		return []schema.Deal{}, err
	}
	deals, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.Deal])
	if err != nil {
		return []schema.Deal{}, err
	}
	return deals, nil
}
