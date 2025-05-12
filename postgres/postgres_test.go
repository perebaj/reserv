//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/perebaj/reserv"
	"github.com/perebaj/reserv/postgres"
	"github.com/stretchr/testify/require"
)

// OpenDB create a new database for testing and return a connection to it.
// Why: For testing, we need a new database for each test to avoid side effects.
// So Opendb creates a new database with a random suffix, and after the test, it drops the database.
func OpenDB(t *testing.T) *sqlx.DB {
	t.Helper()

	cfg := postgres.Config{
		URL:             os.Getenv("POSTGRES_URL"),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	if cfg.URL == "" {
		t.Skip("POSTGRES_URL is not set")
	}

	db, err := sql.Open("postgres", cfg.URL)
	require.NoError(t, err, "error connecting to Postgres: %v", err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	require.NoError(t, err, "error pinging Postgres: %v", err)

	// create a new database with random suffix
	postgresURL, err := url.Parse(cfg.URL)
	require.NoError(t, err, "error parsing Postgres connection URL: %v", err)

	database := strings.TrimLeft(postgresURL.Path, "/")
	randSuffix := fmt.Sprintf("%x", time.Now().UnixNano())

	database = fmt.Sprintf("%s-%x", database, randSuffix)
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, database))
	require.NoError(t, err, "error creating database for test: %v", err)

	postgresURL.Path = "/" + database
	cfg.URL = postgresURL.String()
	testDB, err := postgres.OpenDB(cfg)
	require.NoError(t, err, "error opening test database: %v", err)

	err = postgres.Migrate(testDB.DB)
	require.NoError(t, err, "error running migrations: %v", err)

	// after run the tests, drop the database
	t.Cleanup(func() {
		defer func() {
			_ = testDB.Close()
		}()

		defer func() {
			_ = db.Close()
		}()
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE "%s" WITH (FORCE);`, database))
		require.NoError(t, err, "error dropping database for test: %v", err)
	})

	return testDB
}

// Test a ping to the database
func TestPing(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	err := db.Ping()
	require.NoError(t, err)
}

func TestGetPropertyAmenities(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()
	ctx := context.Background()
	repo := postgres.NewRepository(db)

	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "123e4567-e89b-12d3-a456-426614174000",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)

	amenities := []string{"wifi", "pool"}
	err = repo.CreatePropertyAmenities(ctx, propertyID, amenities)
	require.NoError(t, err)

	propertyAmenities, err := repo.GetPropertyAmenities(ctx, propertyID)
	require.NoError(t, err)
	for _, amenity := range propertyAmenities {
		require.Contains(t, amenities, amenity.ID)
		require.NotNil(t, amenity.CreatedAt)
		// Not zero because we are inserting the created_at on the db layer.
		require.NotZero(t, amenity.CreatedAt)
	}
}

func TestCreatePropertyAmenities(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := postgres.NewRepository(db)
	ctx := context.Background()

	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "123e4567-e89b-12d3-a456-426614174000",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)

	amenities := []string{"wifi", "pool"}
	err = repo.CreatePropertyAmenities(ctx, propertyID, amenities)
	require.NoError(t, err)

	var propertyAmenities []reserv.PropertyAmenity
	err = db.SelectContext(ctx, &propertyAmenities, "SELECT * FROM property_amenities WHERE property_id = $1", propertyID)
	require.NoError(t, err)

	for _, amenity := range propertyAmenities {
		require.Contains(t, amenities, amenity.AmenityID)
		require.Equal(t, propertyID, amenity.PropertyID)
		require.NotNil(t, amenity.CreatedAt)
		// Not zero because we are inserting the created_at on the db layer.
		require.NotZero(t, amenity.CreatedAt)
	}

	// Test what happens when we try to insert the same amenity twice.
	err = repo.CreatePropertyAmenities(ctx, propertyID, amenities)
	require.NoError(t, err)
}

func TestCreateProperty(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := postgres.NewRepository(db)
	ctx := context.Background()

	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "123e4567-e89b-12d3-a456-426614174000",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)
	require.NotEmpty(t, propertyID)

	var createdProperty reserv.Property
	err = db.GetContext(ctx, &createdProperty, "SELECT * FROM properties WHERE id = $1", propertyID)
	require.NoError(t, err)
	require.Equal(t, property.Title, createdProperty.Title)
	require.Equal(t, property.Description, createdProperty.Description)
	require.Equal(t, property.PricePerNightCents, createdProperty.PricePerNightCents)
	require.Equal(t, property.Currency, createdProperty.Currency)
	require.Equal(t, property.HostID, createdProperty.HostID)
	require.NotNil(t, createdProperty.CreatedAt)
	require.NotNil(t, createdProperty.UpdatedAt)
}

func TestUpdateProperty(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := postgres.NewRepository(db)
	ctx := context.Background()

	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "123e4567-e89b-12d3-a456-426614174000",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)
	require.NotEmpty(t, propertyID)

	property.Title = "Updated Property"
	property.PricePerNightCents = 15000
	property.UpdatedAt = property.UpdatedAt.Add(time.Minute)
	property.Description = "Updated Description"
	property.Currency = "BRL"

	err = repo.UpdateProperty(ctx, property, propertyID)
	require.NoError(t, err)

	var updatedProperty reserv.Property
	err = db.GetContext(ctx, &updatedProperty, "SELECT * FROM properties WHERE id = $1", propertyID)
	require.NoError(t, err)
	require.Equal(t, property.Title, updatedProperty.Title)
	require.Equal(t, property.Description, updatedProperty.Description)
	require.Equal(t, property.PricePerNightCents, updatedProperty.PricePerNightCents)
	require.Equal(t, property.Currency, updatedProperty.Currency)
	require.Equal(t, property.HostID, updatedProperty.HostID)
	require.NotNil(t, updatedProperty.CreatedAt)
	require.NotNil(t, updatedProperty.UpdatedAt)
}

func TestDeleteProperty(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := postgres.NewRepository(db)
	ctx := context.Background()

	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "123e4567-e89b-12d3-a456-426614174000",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)
	require.NotEmpty(t, propertyID)

	amenities := []string{"wifi", "pool"}
	err = repo.CreatePropertyAmenities(ctx, propertyID, amenities)
	require.NoError(t, err)

	err = repo.DeleteProperty(ctx, propertyID)
	require.NoError(t, err)

	var deletedProperty reserv.Property
	affected, deletedProperty, err := repo.GetProperty(ctx, propertyID)
	require.NoError(t, err)
	require.Equal(t, 0, affected)
	require.Empty(t, deletedProperty)

	gotAmenities, err := repo.GetPropertyAmenities(ctx, propertyID)
	require.NoError(t, err)
	require.Len(t, gotAmenities, 0)
}
