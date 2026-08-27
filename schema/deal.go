package schema

type Deal struct {
	DestinationID        string  `json:"destination_id"`
	Name                 string  `json:"name"`
	Country              string  `json:"country"`
	Price                float64 `json:"price"`
	AveragePrice         float64 `json:"average_price"`
	DiscountPercentage   float64 `json:"discount_percentage"`
	FlightLink           string  `json:"flight_link"`
	SerpApiFlightLink    string  `json:"serp_api_flight_link"`
	Thumbnail            string  `json:"thumbnail"`
	StartDate            string  `json:"start_date"`
	EndDate              string  `json:"end_date"`
	DepartureAirportCode string  `json:"departure_airport_code"`
	FlightDuration       int     `json:"flight_duration"`
	Stops                int     `json:"stops"`
	Airline              string  `json:"airline"`
	AirlineCode          string  `json:"airline_code"`
	Description          string  `json:"description"`
	Highlights           string  `json:"highlights"`
}

type Deals struct {
	Deals []Deal `json:"deals"`
}
