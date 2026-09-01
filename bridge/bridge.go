package bridge

import (
	"encoding/json"
	"github.com/Keith1039/foptimize/schema"
)

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
