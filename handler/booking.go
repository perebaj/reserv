package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/perebaj/reserv"
)

//go:generate mockgen -source booking.go -destination ../mock/booking.go -package mock
type BookingRepository interface {
	CreateBooking(ctx context.Context, booking reserv.Booking) (string, error)
	DeleteBooking(ctx context.Context, id string) error
	GetBooking(ctx context.Context, id string) (int, reserv.Booking, error)
}

// CreateBooking is the request body for creating a booking.
type CreateBooking struct {
	PropertyID string `json:"property_id"`
	GuestID    string `json:"guest_id"`
	// CheckInDate and CheckOutDate are the dates of the booking.
	// They must be stored in UTC timezone. and must be in the format YYYY-MM-DD.
	// Format: 2025-01-01T00:00:00Z
	CheckInDate  string `json:"check_in_date"`
	CheckOutDate string `json:"check_out_date"`
}

// CreateBookingHandler is the handler for creating a booking.
func (h *Handler) CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateBooking
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode request body", "error", err)
		NewAPIError("invalid_request_body", "invalid request body", http.StatusBadRequest).Write(w)
		return
	}

	if req.PropertyID == "" || req.GuestID == "" || req.CheckInDate == "" || req.CheckOutDate == "" {
		NewAPIError("missing_required_fields", "missing required fields", http.StatusBadRequest).Write(w)
		return
	}

	checkInDate, err := time.Parse(time.RFC3339, req.CheckInDate)
	if err != nil {
		slog.Error("invalid date format", "error", err)
		NewAPIError("invalid_date_format", "invalid date format", http.StatusBadRequest).Write(w)
		return
	}

	checkOutDate, err := time.Parse(time.RFC3339, req.CheckOutDate)
	if err != nil {
		slog.Error("invalid date format", "error", err)
		NewAPIError("invalid_date_format", "invalid date format", http.StatusBadRequest).Write(w)
		return
	}

	booking := reserv.Booking{
		PropertyID: req.PropertyID,
		GuestID:    req.GuestID,
		// convert data to time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC)
		CheckInDate:  time.Date(checkInDate.Year(), checkInDate.Month(), checkInDate.Day(), 0, 0, 0, 0, time.UTC),
		CheckOutDate: time.Date(checkOutDate.Year(), checkOutDate.Month(), checkOutDate.Day(), 0, 0, 0, 0, time.UTC),
		// TODO(@perebaj): Add total price. and currency.
	}

	id, err := h.bookingRepo.CreateBooking(r.Context(), booking)
	if err != nil {
		slog.Error("failed to create booking", "error", err)
		NewAPIError("failed_to_create_booking", "failed to create booking", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"id": id})
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		NewAPIError("failed_to_encode_response", "failed to encode response", http.StatusInternalServerError).Write(w)
		return
	}
}
