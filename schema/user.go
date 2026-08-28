package schema

type User struct {
	Email               string   `json:"email" db:"email"`
	RelevantAirports    []string `json:"relevant_airports" db:"relevant_airports"`
	SubscribedCountries []string `json:"subscribed_countries" db:"subscribed_countries"`
	Config              Config   `json:"config" db:"config"`
	Threshold           int      `json:"threshold" db:"threshold"`
}

type Users struct {
	Users []User `json:"users"`
}
