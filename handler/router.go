package handler

import (
	"net/http"
)

// Handle is responsable to gather all important implementations to inject into the handler.
// TODO(@perebaj): As we have a small number of handlers, we can keep them here.
type Handler struct {
	repo       PropertyRepository
	CloudFlare CloudFlareAPI
}

func NewHandler(repo PropertyRepository, cloudFlare CloudFlareAPI) *Handler {
	return &Handler{repo: repo, CloudFlare: cloudFlare}
}

// RegisterRoutes registers all property routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/properties", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Has("id") {
				h.GetProperty(w, r)
			} else {
				h.GetProperties(w, r)
			}
		case http.MethodPost:
			h.CreateProperty(w, r)
		case http.MethodPut:
			if r.URL.Query().Has("id") {
				h.UpdateProperty(w, r)
			} else {
				http.Error(w, "Missing property ID", http.StatusBadRequest)
			}
		case http.MethodDelete:
			if r.URL.Query().Has("id") {
				h.DeleteProperty(w, r)
			} else {
				http.Error(w, "Missing property ID", http.StatusBadRequest)
			}
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
}
