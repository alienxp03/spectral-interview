package sqlite

import "time"

// A Usage represents a Usage in the database
type Usage struct {
	ID    int       `json:"id"`
	Time  time.Time `json:"time"`
	Usage float64   `json:"usage"`
}
