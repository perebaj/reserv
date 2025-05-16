// Package postgres contains all the database operations.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/perebaj/reserv"
)

// CreateBooking creates a new booking considering the existing bookings to avoid overlapping.
// An important detail about this implementation is that the booking will never accept overlapping. So, if a property has booked 2025-01-01 to 2025-01-05
// and another booking is requested for 2025-01-05 to 2025-01-10, the booking will be rejected. We are not considering hours of check in and check out.
func (r *Repository) CreateBooking(ctx context.Context, newBooking reserv.Booking) (string, error) {
	query := `
		SELECT EXISTS (
			-- SELECT will return 1 if there is a booking that overlaps with the new booking
			SELECT 1 FROM bookings WHERE property_id = $1 AND (
				(check_in_date <= $2 AND check_out_date >= $2) OR
				(check_in_date <= $3 AND check_out_date >= $3) OR
				(check_in_date >= $2 AND check_out_date <= $3)
			)
		)
	`

	var isBooked bool
	if err := r.db.GetContext(ctx, &isBooked, query, newBooking.PropertyID, newBooking.CheckInDate, newBooking.CheckOutDate); err != nil {
		return "", fmt.Errorf("failed to check if booking exists: %v", err)
	}

	if isBooked {
		return "", fmt.Errorf("booking overlaps")
	}

	slog.Info("creating booking")
	q := `
		INSERT INTO bookings (
			property_id,
			guest_id,
			check_in_date,
			check_out_date,
			total_price_cents,
			currency,
			created_at,
			updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id string
	if err := r.db.QueryRowxContext(ctx, q,
		newBooking.PropertyID,
		newBooking.GuestID,
		newBooking.CheckInDate,
		newBooking.CheckOutDate,
		newBooking.TotalPriceCents,
		newBooking.Currency,
		newBooking.CreatedAt,
		newBooking.UpdatedAt,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("failed to create booking: %v", err)
	}

	return id, nil
}

// GetBooking returns a booking by id.
func (r *Repository) GetBooking(ctx context.Context, id string) (int, reserv.Booking, error) {
	slog.Info("getting booking", "id", id)
	query := `
		SELECT * FROM bookings WHERE id = $1
	`

	var booking reserv.Booking
	if err := r.db.GetContext(ctx, &booking, query, id); err != nil {
		if err == sql.ErrNoRows {
			return 0, reserv.Booking{}, nil
		}
		return 0, reserv.Booking{}, fmt.Errorf("failed to get booking: %v", err)
	}

	return 1, booking, nil
}

// DeleteBooking deletes a booking by id.
func (r *Repository) DeleteBooking(ctx context.Context, id string) error {
	slog.Info("deleting booking", "id", id)
	query := `
		DELETE FROM bookings WHERE id = $1
	`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete booking: %v", err)
	}

	return nil
}

// Bookings returns all bookings.
func (r *Repository) Bookings(ctx context.Context, filter reserv.BookingFilter) ([]reserv.Booking, error) {
	slog.Info("getting all bookings")
	query := `
		SELECT * FROM bookings
	`

	// To avoid passing arguments that wont be used, we will build the query and the arguments separately.
	args := make(map[string]interface{}, 2)

	if filter.PropertyID != "" {
		query += " WHERE property_id = :property_id"
		args["property_id"] = filter.PropertyID
	}

	if filter.GuestID != "" {
		query += " WHERE guest_id = :guest_id"
		args["guest_id"] = filter.GuestID
	}

	query += " ORDER BY created_at DESC"
	slog.Info("final query for bookings", "query", query, "args", args)
	var bookings []reserv.Booking
	res, err := r.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %v", err)
	}

	for res.Next() {
		var booking reserv.Booking
		if err := res.StructScan(&booking); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %v", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// GetBookingsByHostID returns all bookings by host id.
func (r *Repository) GetBookingsByHostID(ctx context.Context, hostID string) ([]reserv.Booking, error) {
	slog.Info("getting bookings by host id", "host_id", hostID)
	query := `
		SELECT * FROM bookings WHERE host_id = $1 ORDER BY created_at DESC
	`

	var bookings []reserv.Booking
	if err := r.db.SelectContext(ctx, &bookings, query, hostID); err != nil {
		return nil, fmt.Errorf("failed to get bookings by host id: %v", err)
	}

	return bookings, nil
}
