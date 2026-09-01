package bridge

import (
	"encoding/json"
	"github.com/Keith1039/foptimize/config"
	"github.com/Keith1039/foptimize/db"
	"github.com/Keith1039/foptimize/schema"
	"log"
)

var dbClient db.DatabaseClient

func init() {
	var err error
	dbClient, err = db.NewDatabaseClient(config.DB_URL)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseDeals(response map[string]interface{}) ([]schema.Deal, error) {
	var deals schema.Deals
	bytes, err := json.Marshal(response)
	if err != nil {
		return []schema.Deal{}, err
	}
	err = json.Unmarshal(bytes, &deals)
	if err != nil {
		return []schema.Deal{}, err
	}
	return deals.Deals, nil
}
