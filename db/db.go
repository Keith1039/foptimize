package db

import (
	"context"
	"errors"
	"fmt"
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

func NewDatabaseClientWithPool(pool *pgxpool.Pool) (DatabaseClient, error) {
	return DatabaseClient{db: pool}, nil
}

// User section

func (client DatabaseClient) AddUser(ctx context.Context, user schema.User) (int, error) {
	args := pgx.NamedArgs{
		"EMAIL":                user.Email,
		"RELEVANT_AIRPORTS":    user.RelevantAirports,
		"SUBSCRIBED_COUNTRIES": user.SubscribedCountries,
		"CONFIG":               user.Config,
		"EMAIL_CONFIG":         user.EmailConfig,
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
		"EMAIL_CONFIG":         user.EmailConfig,
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

func (client DatabaseClient) GetUsers(ctx context.Context) ([]schema.User, error) {
	rows, err := client.db.Query(ctx, getAllUsersQuery)
	if err != nil {
		return []schema.User{}, err
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.User])
	if err != nil {
		return []schema.User{}, err
	}
	return users, nil
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
func getAddDealArgs(deal schema.Deal) pgx.NamedArgs {
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

func (client DatabaseClient) dealExistsForUser(ctx context.Context, id int, deal schema.Deal) int {
	args := pgx.NamedArgs{
		"USER_ID":                id,
		"NAME":                   deal.Name,
		"COUNTRY":                deal.Country,
		"START_DATE":             deal.StartDate,
		"END_DATE":               deal.EndDate,
		"DEPARTURE_AIRPORT_CODE": deal.DepartureAirportCode,
		"ARRIVAL_AIRPORT_CODE":   deal.ArrivalAirportCode,
		"FLIGHT_DURATION":        deal.FlightDuration,
		"STOPS":                  deal.Stops,
	}
	rows, err := client.db.Query(ctx, getDealForUserQuery, args)
	if err != nil {
		log.Fatal(err)
	}
	dealId, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0 // default value
		} else {
			log.Fatal(err)
		}
	}
	return dealId
}

func arrToMap(countries []string) map[string]bool {
	m := make(map[string]bool)
	for _, country := range countries {
		m[country] = true
	}
	return m
}

func (client DatabaseClient) addDeal(ctx context.Context, tx pgx.Tx, deal schema.Deal) (int, error) {
	// add the deal
	args := getAddDealArgs(deal)
	rows, err := tx.Query(ctx, addDealQuery, args)
	if err != nil {
		// doesn't matter
		return 0, err
	}
	dealId, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return 0, err
	}
	return dealId, nil
}

func (client DatabaseClient) updateDeal(ctx context.Context, tx pgx.Tx, dealId int, deal schema.Deal) error {

	args := pgx.NamedArgs{
		"DEAL_ID":             dealId,
		"PRICE":               deal.Price,
		"AVERAGE_PRICE":       deal.AveragePrice,
		"DISCOUNT_PERCENTAGE": deal.DiscountPercentage,
	}
	_, err := tx.Exec(ctx, updateDealPriceQuery, args)
	if err != nil {
		return err
	}
	return nil
}

func (client DatabaseClient) mapUserAndDealNoTx(ctx context.Context, userId, dealId int) error {
	mappingArgs := pgx.NamedArgs{
		"USER_ID": userId,
		"DEAL_ID": dealId,
	}
	_, err := client.db.Exec(ctx, mapUserAndDeal, mappingArgs)
	if err != nil {
		return err
	}
	return nil
}

func (client DatabaseClient) mapUserAndDealWithTx(ctx context.Context, tx pgx.Tx, userId int, dealId int) error {
	mappingArgs := pgx.NamedArgs{
		"USER_ID": userId,
		"DEAL_ID": dealId,
	}
	_, err := tx.Exec(ctx, mapUserAndDeal, mappingArgs)
	if err != nil {
		return err
	}
	return nil
}

func (client DatabaseClient) MapUserAndDeal(ctx context.Context, tx pgx.Tx, userId, dealId int) error {
	if tx != nil {
		return client.mapUserAndDealWithTx(ctx, tx, userId, dealId)
	}
	return client.mapUserAndDealNoTx(ctx, userId, dealId)
}

func (client DatabaseClient) RemoveDealAssociation(ctx context.Context, userId, dealId int) error {
	args := pgx.NamedArgs{
		"USER_ID": userId,
		"DEAL_ID": dealId,
	}
	_, err := client.db.Exec(ctx, removeUserAndDealAssociation, args)
	if err != nil {
		return err
	}
	return nil
}

// SaveDeals saves deals and updates existing deal prices to ensure accuracy
// deals given to this function will NOT have ids but deals returned by this function WILL
func (client DatabaseClient) SaveDeals(ctx context.Context, id int, deals []schema.Deal) ([]schema.Deal, error) {
	var rollBackErr error
	var metThreshold []schema.Deal
	user, err := client.GetUser(ctx, id)
	if err != nil {
		return []schema.Deal{}, err
	}
	tx, err := client.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return []schema.Deal{}, err
	}
	countryMap := arrToMap(user.SubscribedCountries)
	for _, deal := range deals {
		// validate deal
		err = deal.Validate()
		if err != nil {
			return []schema.Deal{}, fmt.Errorf("deal with link '%s' failed validation with error %w", deal.FlightLink, err)
		}

		// check if it's a subscribed country... if not, skip iteration
		if _, ok := countryMap[deal.Country]; !ok {
			continue
		}

		// update deals that already exist for the user and skip to the next one
		dealId := client.dealExistsForUser(ctx, id, deal)
		if dealId != 0 {
			// update price and the fields around price such as average price and discount percentage
			err = client.updateDeal(ctx, tx, dealId, deal)
			// if it fails it's whatever, doesn't matter just log for now
			if err != nil {
				fmt.Printf("failed to update price fields for deal with id: '%d' with error '%v' \n", dealId, err)
			}
			continue
		}

		dealId, err = client.addDeal(ctx, tx, deal)
		if err != nil {
			rollBackErr = tx.Rollback(ctx)
			fmt.Printf("Rollback Error occured: %v", rollBackErr)
			return []schema.Deal{}, err
		}
		// map to user
		err = client.MapUserAndDeal(ctx, tx, id, dealId)
		if err != nil {
			// doesn't matter
			rollBackErr = tx.Rollback(ctx)
			fmt.Printf("Rollback Error occured: %v", rollBackErr)
			return []schema.Deal{}, err
		}
		// next we check if threshold is met, so we can return them
		if deal.DiscountPercentage >= float64(user.Threshold) {
			// set the id before adding to array
			deal.Id = dealId
			metThreshold = append(metThreshold, deal)
		}

	}

	err = tx.Commit(ctx)
	if err != nil {
		return []schema.Deal{}, err
	}
	return metThreshold, err
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

func (client DatabaseClient) GetDeal(ctx context.Context, dealId int) (schema.Deal, error) {
	rows, err := client.db.Query(ctx, getDealQuery, dealId)
	if err != nil {
		return schema.Deal{}, err
	}
	deal, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[schema.Deal])
	if err != nil {
		return schema.Deal{}, err
	}
	return deal, nil
}

func (client DatabaseClient) AddSentEmail(ctx context.Context, id int, emailId string, deals []schema.Deal) error {
	tx, err := client.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	for _, deal := range deals {
		if !deal.HasValidId() {
			return fmt.Errorf("deal with link '%s' has no id attached, cannot associate with sent email", deal.FlightLink)
		}
		mappingArgs := pgx.NamedArgs{
			"EMAIL_ID": emailId,
			"USER_ID":  id,
			"DEAL_ID":  deal.Id,
		}
		_, err = tx.Exec(ctx, addEmailTrackingQuery, mappingArgs)
		if err != nil {
			rollBackErr := tx.Rollback(ctx)
			if rollBackErr != nil {
				fmt.Printf("Rollback Error occured: %v", rollBackErr)
			}
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (client DatabaseClient) GetUnsentEmails(ctx context.Context) ([]schema.UserDealLink, error) {
	rows, err := client.db.Query(ctx, getUnsentEmailsQuery)
	if err != nil {
		return []schema.UserDealLink{}, err
	}
	dealLinks, err := pgx.CollectRows(rows, pgx.RowToStructByName[schema.UserDealLink])
	if err != nil {
		return []schema.UserDealLink{}, err
	}
	return dealLinks, nil
}

func (client DatabaseClient) Close() {
	client.db.Close()
}
