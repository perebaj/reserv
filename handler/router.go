package handler

import (
	"net/http"
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
	mux.HandleFunc("/properties", func(w http.ResponseWriter, r *http.Request) {
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
	})

	mux.HandleFunc("/properties/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.handlerPostImage(w, r)
		}
	})

	mux.HandleFunc("/properties/{id}/amenities", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.PostAmenity(w, r)
		}
	})

	mux.HandleFunc("/amenities", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetAmenities(w, r)
		}
	})

	mux.HandleFunc("/bookings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateBookingHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			h.DeleteBookingHandler(w, r)
		case http.MethodGet:
			h.GetBookingHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
