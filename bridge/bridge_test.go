package bridge_test

import (
	"encoding/json"
	"github.com/Keith1039/foptimize/bridge"
	"os"
	"testing"
)

func TestParseDeals(t *testing.T) {
	var m map[string]interface{}
	bytes, err := os.ReadFile("../example.json")
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		t.Fatal(err)
	}
	deals, err := bridge.ParseDeals(m)
	if err != nil {
		t.Fatal(err)
	}
	if len(deals) == 0 {
		t.Fatal("deals could not be parsed")
	}

}
