package reserv

import "time"

// BookingFilter is the filter for the bookings.
type BookingFilter struct {
	// PropertyID is the id of the property that the booking is for.
	PropertyID string
	// GuestID is the id of the guest(user) who made the booking.
	GuestID string
}

// Booking is the entity that represents a booking of a property by a guest(user).
type Booking struct {
	ID         string `json:"id" db:"id"`
	PropertyID string `json:"property_id" db:"property_id"`
	GuestID    string `json:"guest_id" db:"guest_id"`
	// CheckInDate and CheckOutDate are the dates of the booking.
	// They must be stored in UTC timezone. and must be in the format YYYY-MM-DD.
	// Format: 2025-01-01T00:00:00Z
	CheckInDate     time.Time `json:"check_in_date" db:"check_in_date"`
	CheckOutDate    time.Time `json:"check_out_date" db:"check_out_date"`
	TotalPriceCents int       `json:"total_price_cents" db:"total_price_cents"`
	Currency        string    `json:"currency" db:"currency"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
