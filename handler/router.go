package handler

import (
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

// Handler is responsable to gather all important implementations to inject into the handler.
// TODO(@perebaj): As we have a small number of handlers, we can keep them here.
type Handler struct {
	repo        PropertyRepository
	bookingRepo BookingRepository
	CloudFlare  CloudFlareAPI
}

// NewHandler creates a new handler
func NewHandler(repo PropertyRepository, cloudFlare CloudFlareAPI, bookingRepo BookingRepository) *Handler {
	return &Handler{repo: repo, CloudFlare: cloudFlare, bookingRepo: bookingRepo}
}

// RegisterRoutes registers all property routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/properties", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetProperties(w, r)
		case http.MethodPost:
			h.CreateProperty(w, r)
		case http.MethodOptions:
			// Handle preflight requests
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/properties/{id}", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetProperty(w, r)
		case http.MethodPut:
			h.UpdateProperty(w, r)
		case http.MethodDelete:
			h.DeleteProperty(w, r)
		case http.MethodOptions:
			// Handle preflight requests
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/images", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.handlerPostImage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/images/{id}", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			h.handlerDeleteImage(w, r)
		}
	})))

	mux.Handle("/properties/{id}/amenities", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.PostAmenity(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.HandleFunc("/amenities", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetAmenities(w, r)
		}
	})

	mux.Handle("/bookings", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateBookingHandler(w, r)
		case http.MethodGet:
			h.BookingsHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/bookings/{id}", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			h.DeleteBookingHandler(w, r)
		case http.MethodGet:
			h.GetBookingHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/protected", clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(protectedHandler)))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("protected route")
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slog.Info("protected route", "claims", claims)
	slog.Info("protected route", "userID", claims.Subject)
	w.WriteHeader(http.StatusOK)
}
