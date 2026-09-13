package schema

type UserDealLink struct {
	UserId     int    `db:"user_id"`
	FlightLink string `db:"flight_link"`
}
