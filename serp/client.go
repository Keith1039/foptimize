package serp

import (
	"github.com/Keith1039/foptimize/config"
	"github.com/serpapi/serpapi-golang"
)

var Client serpapi.SerpApiClient

func init() {
	// replace with your SerpApi key
	setting := serpapi.NewSerpApiClientSetting(config.SERPAPI_KEY)
	setting.Engine = "google_flights_deals"
	// initialize the client
	Client = serpapi.NewClient(setting)
}
