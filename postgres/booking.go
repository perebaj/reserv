package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/perebaj/reserv"
)

// CreateBooking creates a new booking considering the existing bookings to avoid overlapping.
func (r *Repository) CreateBooking(ctx context.Context, newBooking reserv.Booking) (string, error) {

	// get all curBookings between check in and check out date
	var curBookings []reserv.Booking
	query := `
		SELECT * FROM bookings WHERE property_id = $1 AND check_in_date <= $2 AND check_out_date >= $3
	`

	if err := r.db.SelectContext(ctx, &curBookings, query, newBooking.PropertyID, newBooking.CheckInDate, newBooking.CheckOutDate); err != nil {
		slog.Error("failed to get bookings", "error", err)
		return "", fmt.Errorf("failed to get bookings: %v", err)
	}

	// if there are any bookings that overlap with the new booking, then we can't create the booking
	for _, curBooking := range curBookings {
		// if the new booking is between the current booking it's overlapping. Then we can't create the booking
		if curBooking.CheckInDate < newBooking.CheckOutDate && curBooking.CheckOutDate > newBooking.CheckInDate {
			slog.Error("booking overlaps", "new_booking", newBooking, "cur_booking", curBooking)
			return "", errors.New("booking overlaps")
		}
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

// GetBookings returns all bookings.
func (r *Repository) GetBookings(ctx context.Context) ([]reserv.Booking, error) {
	slog.Info("getting all bookings")
	query := `
		SELECT * FROM bookings ORDER BY created_at DESC
	`

	var bookings []reserv.Booking
	if err := r.db.SelectContext(ctx, &bookings, query); err != nil {
		return nil, fmt.Errorf("failed to get bookings: %v", err)
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
