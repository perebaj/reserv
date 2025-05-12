package postgres

import (
	"context"
	"database/sql"

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
		return nil, err
	}

	return amenities, nil
}

// GetPropertyAmenities returns the amenities for a property.
func (r *Repository) GetPropertyAmenities(ctx context.Context, propertyID string) ([]reserv.Amenity, error) {
	query := `
		SELECT a.id, a.name, a.created_at
		FROM amenities a
		JOIN property_amenities pa ON a.id = pa.amenity_id
		WHERE pa.property_id = $1
		ORDER BY a.name
	`

	rows, err := r.db.QueryContext(ctx, query, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amenities []reserv.Amenity
	for rows.Next() {
		var amenity reserv.Amenity
		if err := rows.Scan(&amenity.ID, &amenity.Name, &amenity.CreatedAt); err != nil {
			return nil, err
		}
		amenities = append(amenities, amenity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return amenities, nil
}

// CreatePropertyAmenities creates the amenities for a property.
func (r *Repository) CreatePropertyAmenities(ctx context.Context, propertyID string, amenities []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
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
			return err
		}
	}

	return tx.Commit()
}

// CreateProperty ...
func (r *Repository) CreateProperty(ctx context.Context, property reserv.Property) (string, error) {
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
		return "", err
	}

	return id, nil
}

// UpdateProperty ...
func (r *Repository) UpdateProperty(ctx context.Context, property reserv.Property, id string) error {
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
		return err
	}

	return nil
}

// DeleteProperty deletes a property and its amenities.
func (r *Repository) DeleteProperty(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		DELETE FROM properties WHERE id = $1
	`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	query = `
		DELETE FROM property_amenities WHERE property_id = $1
	`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return tx.Commit()
}

// GetProperty returns a property by id. It return the number of rows affected and the property.
func (r *Repository) GetProperty(ctx context.Context, id string) (int, reserv.Property, error) {
	query := `
		SELECT * FROM properties WHERE id = $1
	`

	var property reserv.Property
	// if no rows in result error, just return 0 and empty property
	if err := r.db.GetContext(ctx, &property, query, id); err != nil {
		if err == sql.ErrNoRows {
			return 0, reserv.Property{}, nil
		}
		return 0, reserv.Property{}, err
	}

	return 1, property, nil
}
