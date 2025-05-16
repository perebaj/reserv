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

// BookingRepository is the repository for the booking. Gathers all the methods to interact with the booking.
type BookingRepository interface {
	CreateBooking(ctx context.Context, booking reserv.Booking) (string, error)
	DeleteBooking(ctx context.Context, id string) error
	GetBooking(ctx context.Context, id string) (int, reserv.Booking, error)
	Bookings(ctx context.Context, filter reserv.BookingFilter) ([]reserv.Booking, error)
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

const dateFormat = "2006-01-02" // This is Go's way of specifying YYYY-MM-DD

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

	checkInDate, err := time.Parse(dateFormat, req.CheckInDate)
	if err != nil {
		slog.Error("invalid date format", "error", err, "date", req.CheckInDate)
		NewAPIError("invalid_date_format", "invalid date format. Expected YYYY-MM-DD", http.StatusBadRequest).Write(w)
		return
	}

	checkOutDate, err := time.Parse(dateFormat, req.CheckOutDate)
	if err != nil {
		slog.Error("invalid date format", "error", err, "date", req.CheckOutDate)
		NewAPIError("invalid_date_format", "invalid date format. Expected YYYY-MM-DD", http.StatusBadRequest).Write(w)
		return
	}

	booking := reserv.Booking{
		PropertyID:   req.PropertyID,
		GuestID:      req.GuestID,
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

// GetBookingHandler is the handler for getting a booking by id.
func (h *Handler) GetBookingHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		NewAPIError("missing_id", "missing id", http.StatusBadRequest).Write(w)
		return
	}

	affectedRows, booking, err := h.bookingRepo.GetBooking(r.Context(), id)
	if err != nil {
		slog.Error("failed to get booking", "error", err)
		NewAPIError("failed_to_get_booking", "failed to get booking", http.StatusInternalServerError).Write(w)
		return
	}

	if affectedRows == 0 {
		NewAPIError("booking_not_found", "booking not found", http.StatusNotFound).Write(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(booking)
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		NewAPIError("failed_to_encode_response", "failed to encode response", http.StatusInternalServerError).Write(w)
		return
	}
}

// BookingsHandler is the handler for getting all bookings.
func (h *Handler) BookingsHandler(w http.ResponseWriter, r *http.Request) {
	propertyID := r.URL.Query().Get("property_id")
	guestID := r.URL.Query().Get("guest_id")

	bookings, err := h.bookingRepo.Bookings(r.Context(), reserv.BookingFilter{
		PropertyID: propertyID,
		GuestID:    guestID,
	})
	if err != nil {
		slog.Error("failed to get bookings", "error", err)
		NewAPIError("failed_to_get_bookings", "failed to get bookings", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(bookings)
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		NewAPIError("failed_to_encode_response", "failed to encode response", http.StatusInternalServerError).Write(w)
		return
	}
}

func (h *Handler) DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		NewAPIError("missing_id", "missing id", http.StatusBadRequest).Write(w)
		return
	}

	err := h.bookingRepo.DeleteBooking(r.Context(), id)
	if err != nil {
		slog.Error("failed to delete booking", "error", err)
		NewAPIError("failed_to_delete_booking", "failed to delete booking", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
