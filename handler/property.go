// Package handler contains all the handlers for the application.
package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/perebaj/reserv"
)

//go:generate mockgen -source property.go -destination ../mock/property.go -package mock

// PropertyRepository is the repository for the property. Gathers all the methods to interact with the property.
type PropertyRepository interface {
	// Property methods
	// CreateProperty creates a new property
	CreateProperty(ctx context.Context, property reserv.Property) (string, error)
	// UpdateProperty updates an existing property
	UpdateProperty(ctx context.Context, property reserv.Property, id string) error
	// DeleteProperty deletes a property
	DeleteProperty(ctx context.Context, id string) error
	// GetProperty gets a property by id
	GetProperty(ctx context.Context, id string) (int, reserv.Property, error)
	// Properties gets all properties with sub-resources. Not contains pagination yet.
	// TODO: Add pagination
	Properties(ctx context.Context, filter reserv.PropertyFilter) ([]reserv.Property, error)
	// GetPropertyAmenities gets the amenities for a property
	GetPropertyAmenities(ctx context.Context, propertyID string) ([]reserv.Amenity, error)
	// CreatePropertyAmenities creates amenities for a property
	CreatePropertyAmenities(ctx context.Context, propertyID string, amenities []string) error

	// Images methods
	// CreateImage creates an image for a property
	CreateImage(ctx context.Context, image reserv.PropertyImage) (string, error)

	// Amenities methods
	Amenities(ctx context.Context) ([]reserv.Amenity, error)
}

// CreatePropertyRequest represents the request body for creating a property
type CreatePropertyRequest struct {
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	PricePerNightCents int64    `json:"price_per_night_cents"`
	Currency           string   `json:"currency"`
	HostID             string   `json:"host_id"`
	Amenities          []string `json:"amenities"`
}

// CreateProperty creates a new property
func (h *Handler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var req CreatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", "error", err)
		NewAPIError("invalid_request_body", "invalid request body", http.StatusBadRequest).Write(w)
		return
	}

	if req.Title == "" || req.Description == "" || req.PricePerNightCents == 0 || req.Currency == "" || req.HostID == "" {
		slog.Error("missing required fields", "request", req)
		NewAPIError("missing_required_fields", "missing required fields", http.StatusBadRequest).Write(w)
		return
	}

	now := time.Now()
	property := reserv.Property{
		Title:              req.Title,
		Description:        req.Description,
		PricePerNightCents: req.PricePerNightCents,
		Currency:           req.Currency,
		HostID:             req.HostID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	id, err := h.repo.CreateProperty(r.Context(), property)
	if err != nil {
		slog.Error("failed to create property", "error", err)
		NewAPIError("create_property_error", "failed to create property", http.StatusInternalServerError).Write(w)
		return
	}

	if len(req.Amenities) > 0 {
		if err := h.repo.CreatePropertyAmenities(r.Context(), id, req.Amenities); err != nil {
			slog.Error("failed to create property amenities", "error", err)
			NewAPIError("create_property_amenities_error", "failed to create property amenities", http.StatusInternalServerError).Write(w)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"id": id})
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		NewAPIError("encode_response_error", "failed to encode response", http.StatusInternalServerError).Write(w)
		return
	}
}

// UpdatePropertyRequest represents the request body for updating a property
type UpdatePropertyRequest struct {
	Title              string `json:"title"`
	Description        string `json:"description"`
	PricePerNightCents int64  `json:"price_per_night_cents"`
	Currency           string `json:"currency"`
}

// UpdateProperty updates an existing property
func (h *Handler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	propertyID := r.PathValue("id")
	if propertyID == "" {
		NewAPIError("missing_property_id", "missing property id", http.StatusBadRequest).Write(w)
		return
	}

	var req UpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("failed to decode request body", "error", err)
		NewAPIError("invalid_request_body", "invalid request body", http.StatusBadRequest).Write(w)
		return
	}

	// Validate required fields
	if req.Title == "" || req.Description == "" || req.PricePerNightCents == 0 || req.Currency == "" {
		NewAPIError("missing_required_fields", "missing required fields", http.StatusBadRequest).Write(w)
		return
	}

	property := reserv.Property{
		Title:              req.Title,
		Description:        req.Description,
		PricePerNightCents: req.PricePerNightCents,
		Currency:           req.Currency,
		UpdatedAt:          time.Now(),
	}

	if err := h.repo.UpdateProperty(r.Context(), property, propertyID); err != nil {
		slog.Error("failed to update property", "error", err)
		NewAPIError("update_property_error", "failed to update property", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteProperty deletes a property
func (h *Handler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	propertyID := r.PathValue("id")
	if propertyID == "" {
		NewAPIError("missing_property_id", "missing property id", http.StatusBadRequest).Write(w)
		return
	}

	if err := h.repo.DeleteProperty(r.Context(), propertyID); err != nil {
		slog.Error("failed to delete property", "error", err)
		NewAPIError("delete_property_error", "failed to delete property", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetProperty gets a property by id
func (h *Handler) GetProperty(w http.ResponseWriter, r *http.Request) {
	propertyID := r.PathValue("id")
	if propertyID == "" {
		NewAPIError("missing_property_id", "missing property id", http.StatusBadRequest).Write(w)
		return
	}

	affected, property, err := h.repo.GetProperty(r.Context(), propertyID)
	if err != nil {
		slog.Error("failed to get property", "error", err)
		NewAPIError("get_property_error", "failed to get property", http.StatusInternalServerError).Write(w)
		return
	}
	if affected == 0 {
		NewAPIError("property_not_found", "property not found", http.StatusNotFound).Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(property)
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		NewAPIError("encode_response_error", "failed to encode response", http.StatusInternalServerError).Write(w)
		return
	}
}

// GetProperties gets all properties
func (h *Handler) GetProperties(w http.ResponseWriter, r *http.Request) {
	hostID := r.URL.Query().Get("host_id")
	properties, err := h.repo.Properties(r.Context(), reserv.PropertyFilter{HostID: hostID})
	if err != nil {
		slog.Error("failed to get properties", "error", err)
		NewAPIError("get_properties_error", "failed to get properties", http.StatusInternalServerError).Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(properties)
	if err != nil {
		slog.Error("failed to encode response", "error", err)
		NewAPIError("encode_response_error", "failed to encode response", http.StatusInternalServerError).Write(w)
		return
	}
}

// PostAmenity creates amenities for a property
func (h *Handler) PostAmenity(w http.ResponseWriter, r *http.Request) {
	propertyID := r.PathValue("id")
	if propertyID == "" {
		NewAPIError("missing_property_id", "missing property id", http.StatusBadRequest).Write(w)
		return
	}

	var amenties []string
	if err := json.NewDecoder(r.Body).Decode(&amenties); err != nil {
		slog.Error("failed to decode request body", "error", err)
		NewAPIError("invalid_request_body", "invalid request body", http.StatusBadRequest).Write(w)
		return
	}

	if len(amenties) == 0 {
		NewAPIError("missing_amenity_ids", "missing amenity ids", http.StatusBadRequest).Write(w)
		return
	}

	if err := h.repo.CreatePropertyAmenities(r.Context(), propertyID, amenties); err != nil {
		slog.Error("failed to create property amenities", "error", err)
		NewAPIError("create_property_amenities_error", "failed to create property amenities", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetAmenities gets all amenities
func (h *Handler) GetAmenities(w http.ResponseWriter, r *http.Request) {
	amenities, err := h.repo.Amenities(r.Context())
	if err != nil {
		slog.Error("failed to get amenities", "error", err)
		NewAPIError("get_amenities_error", "failed to get amenities", http.StatusInternalServerError).Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	amenitiesBytes, err := json.Marshal(amenities)
	if err != nil {
		slog.Error("failed to marshal amenities", "error", err)
		NewAPIError("marshal_amenities_error", "failed to marshal amenities", http.StatusInternalServerError).Write(w)
		return
	}

	_, _ = w.Write(amenitiesBytes)
}
