package schema_test

import (
	"github.com/Keith1039/foptimize/schema"
	"github.com/brianvoe/gofakeit/v7"
	"testing"
	"time"
)

func genDeal() schema.Deal {
	// get the time dates
	startDate := gofakeit.Date()
	afterDate := startDate.Add(time.Hour * time.Duration(24*gofakeit.Number(14, 21)))
	return schema.Deal{
		DestinationID:        gofakeit.AirlineAirportIATA(),
		Name:                 gofakeit.City(),
		Country:              gofakeit.Country(),
		Price:                gofakeit.Price(500, 4000),
		AveragePrice:         gofakeit.Price(500, 2000),
		DiscountPercentage:   gofakeit.Price(0, 100),
		FlightLink:           gofakeit.URL(),
		SerpApiFlightLink:    gofakeit.URL(),
		Thumbnail:            gofakeit.URL(),
		StartDate:            startDate.String(),
		EndDate:              afterDate.String(),
		DepartureAirportCode: gofakeit.AirlineAirportIATA(),
		ArrivalAirportCode:   gofakeit.AirlineAirportIATA(),
		FlightDuration:       gofakeit.Number(0, 2000),
		Stops:                gofakeit.Number(0, 10),
		Airline:              gofakeit.AirlineAirport(),
		AirlineCode:          "",
		Description:          "",
		Highlights:           "",
	}
}

func TestDeal_Validate(t *testing.T) {
	deal := genDeal()

	deal.StartDate = ""
	err := deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'StartDate' was left as empty string")
	}
	deal.StartDate = gofakeit.Date().String()
	deal.EndDate = ""
	err = deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'EndDate' was left as empty string")
	}
	deal.EndDate = gofakeit.Date().String()
	deal.DepartureAirportCode = ""
	err = deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'DepartureAirportCode' was left as empty string")
	}
	deal.DepartureAirportCode = gofakeit.AirlineAirportIATA()
	deal.ArrivalAirportCode = ""
	err = deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'ArrivalAirportCode' was left as empty string")
	}
	deal.ArrivalAirportCode = gofakeit.AirlineAirportIATA()
	deal.FlightDuration = 0
	err = deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'FlightDuration' was left as default value of 0")
	}
	deal.FlightDuration = gofakeit.Number(0, 2000)
	deal.FlightLink = ""
	err = deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'FlightLink' was left as empty string")
	}
	deal.FlightLink = gofakeit.URL()
	deal.Country = ""
	err = deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'Country' was left as empty")
	}
	deal.Country = gofakeit.Country()
	deal.Price = 0
	err = deal.Validate()
	if err == nil {
		t.Errorf("deal.Validate() should have returned an error since 'Price' was left as 0")
	}
	deal.Price = gofakeit.Price(500, 4000)
	err = deal.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
