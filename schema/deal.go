package schema

type Deal struct {
	DestinationID        string  `json:"destination_id" db:"destination_id"`
	Name                 string  `json:"name" db:"name"`
	Country              string  `json:"country" db:"country"`
	Price                float64 `json:"price" db:"price"`
	AveragePrice         float64 `json:"average_price" db:"average_price"`
	DiscountPercentage   float64 `json:"discount_percentage" db:"discount_percentage"`
	FlightLink           string  `json:"flight_link" db:"flight_link"`
	SerpApiFlightLink    string  `json:"serp_api_flight_link" db:"serp_api_flight_link"`
	Thumbnail            string  `json:"thumbnail" db:"thumbnail"`
	StartDate            string  `json:"start_date" db:"start_date"`
	EndDate              string  `json:"end_date" db:"end_date"`
	DepartureAirportCode string  `json:"departure_airport_code" db:"departure_airport_code"`
	FlightDuration       int     `json:"flight_duration" db:"flight_duration"`
	Stops                int     `json:"stops" db:"stops"`
	Airline              string  `json:"airline" db:"airline"`
	AirlineCode          string  `json:"airline_code" db:"airline_code"`
	Description          string  `json:"description" db:"description"`
	Highlights           string  `json:"highlights" db:"highlights"`
}

type Deals struct {
	Deals []Deal `json:"deals"`
}
