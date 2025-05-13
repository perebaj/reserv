package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/perebaj/reserv"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Amenities returns all amenities. Obs: As we have a small number of amenities, the pagination is not applied.
func (r *Repository) Amenities(ctx context.Context) ([]reserv.Amenity, error) {
	query := `
		SELECT * FROM amenities
	`

	var amenities []reserv.Amenity
	if err := r.db.SelectContext(ctx, &amenities, query); err != nil {
		return nil, fmt.Errorf("failed to get amenities: %v", err)
	}

	return amenities, nil
}

// GetPropertyAmenities returns the amenities for a property.
func (r *Repository) GetPropertyAmenities(ctx context.Context, propertyID string) ([]reserv.Amenity, error) {
	slog.Info("getting property amenities", "propertyID", propertyID)
	query := `
		SELECT a.id, a.name, a.created_at
		FROM amenities a
		JOIN property_amenities pa ON a.id = pa.amenity_id
		WHERE pa.property_id = $1
		ORDER BY a.name
	`

	rows, err := r.db.QueryContext(ctx, query, propertyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get property amenities: %v", err)
	}
	defer rows.Close()

	var amenities []reserv.Amenity
	for rows.Next() {
		var amenity reserv.Amenity
		if err := rows.Scan(&amenity.ID, &amenity.Name, &amenity.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan property amenities: %v", err)
		}
		amenities = append(amenities, amenity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get property amenities: %v", err)
	}

	return amenities, nil
}

// CreatePropertyAmenities creates the amenities for a property.
func (r *Repository) CreatePropertyAmenities(ctx context.Context, propertyID string, amenities []string) error {
	slog.Info("creating property amenities", "propertyID", propertyID, "amenities", amenities)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO property_amenities (property_id, amenity_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`

	for _, amenity := range amenities {
		if _, err := tx.ExecContext(ctx, query, propertyID, amenity); err != nil {
			return fmt.Errorf("failed to create property amenities: %v", err)
		}
	}

	return tx.Commit()
}

// CreateProperty ...
func (r *Repository) CreateProperty(ctx context.Context, property reserv.Property) (string, error) {
	slog.Info("creating property")
	query := `
		INSERT INTO properties (
			title,
			description,
			price_per_night_cents,
			currency,
			host_id,
			created_at,
			updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id string
	if err := r.db.QueryRowxContext(ctx, query,
		property.Title,
		property.Description,
		property.PricePerNightCents,
		property.Currency,
		property.HostID,
		property.CreatedAt,
		property.UpdatedAt,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("failed to create property: %v", err)
	}

	return id, nil
}

// UpdateProperty ...
func (r *Repository) UpdateProperty(ctx context.Context, property reserv.Property, id string) error {
	slog.Info("updating property", "property_id", id)
	query := `
		UPDATE properties
			SET title = $2,
			description = $3,
			price_per_night_cents = $4,
			currency = $5,
			updated_at = $6
		WHERE id = $1
	`

	if _, err := r.db.ExecContext(ctx, query, id, property.Title, property.Description, property.PricePerNightCents, property.Currency, property.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update property: %v", err)
	}

	return nil
}

// DeleteProperty deletes a property and its amenities.
func (r *Repository) DeleteProperty(ctx context.Context, id string) error {
	slog.Info("deleting property", "id", id)
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		DELETE FROM properties WHERE id = $1
	`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete property: %v", err)
	}

	query = `
		DELETE FROM property_amenities WHERE property_id = $1
	`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete property amenities: %v", err)
	}

	return tx.Commit()
}

// GetProperty returns a property by id. It return the number of rows affected and the property.
func (r *Repository) GetProperty(ctx context.Context, id string) (int, reserv.Property, error) {
	slog.Info("getting property", "id", id)
	query := `
		SELECT * FROM properties WHERE id = $1
	`

	var property reserv.Property
	// if no rows in result error, just return 0 and empty property
	if err := r.db.GetContext(ctx, &property, query, id); err != nil {
		if err == sql.ErrNoRows {
			return 0, reserv.Property{}, nil
		}
		return 0, reserv.Property{}, fmt.Errorf("failed to get property: %v", err)
	}

	return 1, property, nil
}

func (r *Repository) Properties(ctx context.Context) ([]reserv.Property, error) {
	slog.Info("getting properties")
	query := `
		SELECT
        p.id, p.host_id, p.title, p.description,
        p.price_per_night_cents, p.currency, p.created_at, p.updated_at,
        COALESCE(
            json_agg(
                DISTINCT jsonb_build_object(
                    'id', pi.id,
                    'host_id', pi.host_id,
                    'property_id', pi.property_id,
                    'cloudflare_id', pi.cloudflare_id,
                    'filename', pi.filename
                )
            ) FILTER (WHERE pi.id IS NOT NULL), '[]'
        ) AS images,
        COALESCE(
            json_agg(
                DISTINCT jsonb_build_object(
                    'id', a.id,
                    'name', a.name
                )
            ) FILTER (WHERE a.id IS NOT NULL), '[]'
        ) AS amenities
    FROM
        properties p
    LEFT JOIN
        property_images pi ON p.id = pi.property_id
    LEFT JOIN
        property_amenities pa ON p.id = pa.property_id
    LEFT JOIN
        amenities a ON pa.amenity_id = a.id
    GROUP BY
        p.id
    ORDER BY
        p.created_at DESC`

	// PropertyWithJSON is a helper struct the JSON fields that we are aggregating in the query.
	type PropertyWithJSON struct {
		reserv.Property
		ImagesJSON    json.RawMessage `db:"images"`
		AmenitiesJSON json.RawMessage `db:"amenities"`
	}

	var propertiesWithJSON []PropertyWithJSON
	err := r.db.SelectContext(ctx, &propertiesWithJSON, query)
	if err != nil {
		return nil, fmt.Errorf("failed to scan properties: %v", err)
	}

	properties := make([]reserv.Property, len(propertiesWithJSON))
	for i, p := range propertiesWithJSON {
		properties[i] = p.Property
		if err := json.Unmarshal(p.ImagesJSON, &properties[i].Images); err != nil {
			return nil, fmt.Errorf("failed to unmarshal images: %v, raw JSON: %s", err, string(p.ImagesJSON))
		}

		if err := json.Unmarshal(p.AmenitiesJSON, &properties[i].Amenities); err != nil {
			return nil, fmt.Errorf("failed to unmarshal amenities: %v, raw JSON: %s", err, string(p.AmenitiesJSON))
		}
	}
	return properties, nil
}

// CreateImage creates an image and associates it with a property.
func (r *Repository) CreateImage(ctx context.Context, image reserv.PropertyImage) (string, error) {
	slog.Info("creating image")
	query := `
		INSERT INTO property_images (
			property_id,
			host_id,
			cloudflare_id,
			filename,
			created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id string
	if err := r.db.QueryRowxContext(ctx, query,
		image.PropertyID,
		image.HostID,
		image.CloudflareID,
		image.Filename,
		image.CreatedAt,
	).Scan(&id); err != nil {
		return "", fmt.Errorf("failed to create image: %v", err)
	}

	return id, nil
}
