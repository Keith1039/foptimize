package schema

import (
	"encoding/json"
	"log"
)

type Config struct {
	TravelDuration string `json:"travel_duration"`
}

func (c Config) TransformParams() map[string]string {
	var m map[string]any
	config := make(map[string]string)
	bytes, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range m {
		val, ok := v.(string)
		if ok && val != "" {
			config[k] = val
		}
	}
	return config
}
