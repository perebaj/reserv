package reserv

import (
	"time"

	"github.com/google/uuid"
)

// Amenity represents a property amenity
type Amenity struct {
	// ID is the unique identifier for the amenity. Required.
	ID string `json:"id" db:"id"`
	// Name is the name of the amenity. Required.
	Name string `json:"name" db:"name"`
	// CreatedAt is the timestamp when the amenity was created. Required.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// PropertyFilter represents a filter for properties
type PropertyFilter struct {
	// HostID is the unique identifier for the host of the property.
	HostID string
}

// Property represents a property listing
type Property struct {
	// ID is the unique identifier for the property. It is generated by the database. Optinal.
	ID uuid.UUID `json:"id" db:"id"`
	// HostID is the unique identifier for the host of the property. Required.
	HostID string `json:"host_id" db:"host_id"`
	// Title is the title of the property. Required.
	Title string `json:"title" db:"title"`
	// Description is the description of the property. Required.
	Description string `json:"description" db:"description"`
	// PricePerNightCents is the price per night in cents. Required.
	PricePerNightCents int64 `json:"price_per_night_cents" db:"price_per_night_cents"`
	// Currency is the currency of the property. Example: "USD", "BRL". Required.
	Currency string `json:"currency" db:"currency"`
	// Amenities is the list of amenities for the property. Example: ["wifi", "pool"].
	Amenities []Amenity `json:"amenities" db:"-"`
	// Images is the list of images for the property.
	Images []PropertyImage `json:"images" db:"-"`
	// CreatedAt is the timestamp when the property was created. Optional.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// UpdatedAt is the timestamp when the property was updated. Required.
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// PropertyImage represents an image for a property.
type PropertyImage struct {
	ID uuid.UUID `json:"id" db:"id"`
	// HostID is the unique identifier for the host of the property. Required.
	HostID string `json:"host_id" db:"host_id"`
	// CreatedAt is the timestamp when the property image was created. Required.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	// PropertyID is the unique identifier for the property. Required.
	PropertyID uuid.UUID `json:"property_id" db:"property_id"`
	// CloudflareID is the unique identifier for the image in Cloudflare. Required.
	CloudflareID uuid.UUID `json:"cloudflare_id" db:"cloudflare_id"`
	// Filename is the filename of the image. Required.
	Filename string `json:"filename" db:"filename"`
}

// PropertyAmenity represents the junction between properties and amenities
type PropertyAmenity struct {
	// PropertyID is the unique identifier for the property. Required.
	PropertyID uuid.UUID `json:"property_id" db:"property_id"`
	// AmenityID is the unique identifier for the amenity. Required.
	AmenityID string `json:"amenity_id" db:"amenity_id"`
	// CreatedAt is the timestamp when the property amenity was created. Required.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
