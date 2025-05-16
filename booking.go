package reserv

import "time"

// Booking is the entity that represents a booking of a property by a guest(user).
type Booking struct {
	ID              string    `json:"id" db:"id"`
	PropertyID      string    `json:"property_id" db:"property_id"`
	GuestID         string    `json:"guest_id" db:"guest_id"`
	CheckInDate     string    `json:"check_in_date" db:"check_in_date"`
	CheckOutDate    string    `json:"check_out_date" db:"check_out_date"`
	TotalPriceCents int       `json:"total_price_cents" db:"total_price_cents"`
	Currency        string    `json:"currency" db:"currency"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
