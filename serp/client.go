package serp

import (
	"github.com/Keith1039/foptimize/config"
	"github.com/serpapi/serpapi-golang"
)

var client serpapi.SerpApiClient

func init() {
	// replace with your SerpApi key
	setting := serpapi.NewSerpApiClientSetting(config.Config.SERPAPI_KEY)
	setting.Engine = "google_flights_deals"
	// initialize the client
	client = serpapi.NewClient(setting)
}

// GetClient returns the configured client
func GetClient() serpapi.SerpApiClient {
	return client
}
