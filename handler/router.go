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
		if r.URL.Query().Has("id") {
			switch r.Method {
			case http.MethodGet:
				h.GetProperty(w, r)
			case http.MethodPut:
				h.UpdateProperty(w, r)
			case http.MethodDelete:
				h.DeleteProperty(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.GetProperties(w, r)
		case http.MethodPost:
			h.CreateProperty(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/images", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			h.handlerPostImage(w, r)
		}
	})
}
