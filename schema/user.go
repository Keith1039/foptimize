package schema

import "time"

type User struct {
	Id                  int       `json:"id" db:"id"`
	Email               string    `json:"email" db:"email"`
	RelevantAirports    []string  `json:"relevant_airports" db:"relevant_airports"`
	SubscribedCountries []string  `json:"subscribed_countries" db:"subscribed_countries"`
	Config              Config    `json:"config" db:"config"`
	Threshold           int       `json:"threshold" db:"threshold"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

type Users struct {
	Users []User `json:"users"`
}
