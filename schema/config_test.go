package schema_test

import (
	"github.com/Keith1039/foptimize/schema"
	"maps"
	"testing"
)

func TestConfig_TransformParams(t *testing.T) {
	// full config
	config := schema.Config{
		Type:           "1",
		TravelClass:    "1",
		TravelDuration: "1",
		Stops:          "1",
		MaxDuration:    "1",
	}
	expected := map[string]string{
		"type":            "1",
		"travel_class":    "1",
		"travel_duration": "1",
		"stops":           "1",
		"max_duration":    "1",
	}
	actual := config.TransformParams()
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	config = schema.Config{
		Type:           "1",
		TravelClass:    "1",
		TravelDuration: "1",
	}
	expected = map[string]string{
		"type":            "1",
		"travel_class":    "1",
		"travel_duration": "1",
	}
	actual = config.TransformParams()
	if !maps.Equal(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
