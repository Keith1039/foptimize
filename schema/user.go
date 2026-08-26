package schema

type User struct {
	Email               string   `json:"email"`
	RelevantAirports    []string `json:"relevant_airports"`
	SubscribedCountries []string `json:"subscribed_countries"`
	Config              Config   `json:"config"`
	Threshold           int      `json:"threshold"`
}

type Users struct {
	Users []User `json:"users"`
}
