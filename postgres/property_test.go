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

	"github.com/google/uuid"
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

func TestAmenities(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := postgres.NewRepository(db)
	ctx := context.Background()

	amenities, err := repo.Amenities(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, amenities)
	for _, amenity := range amenities {
		require.NotEmpty(t, amenity.ID)
		require.NotEmpty(t, amenity.Name)
		require.NotNil(t, amenity.CreatedAt)
		require.NotZero(t, amenity.CreatedAt)
	}
}

func TestGetPropertyAmenities(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()
	ctx := context.Background()
	repo := postgres.NewRepository(db)
	hostID := "user_2x5CiRO5Mf0wBpWO8w469jEJhRq"
	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             hostID,
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
		HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
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

	propertyIDUUID, err := uuid.Parse(propertyID)
	require.NoError(t, err)

	for _, amenity := range propertyAmenities {
		require.Contains(t, amenities, amenity.AmenityID)
		require.Equal(t, propertyIDUUID, amenity.PropertyID)
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
		HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
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
		HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
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
		HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)
	require.NotEmpty(t, propertyID)

	amenities := []string{"wifi", "pool"}
	err = repo.CreatePropertyAmenities(ctx, propertyID, amenities)
	require.NoError(t, err)

	image := reserv.PropertyImage{
		PropertyID:   uuid.MustParse(propertyID),
		HostID:       "2c02e000-42f6-4587-8244-a290421b9c4f",
		CloudflareID: uuid.MustParse("2e195545-8278-41a8-9d01-3c423ec71263"),
		Filename:     "test.jpg",
		CreatedAt:    time.Now(),
	}

	booking := reserv.Booking{
		PropertyID:      propertyID,
		GuestID:         "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CheckInDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CheckOutDate:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Currency:        "USD",
		TotalPriceCents: 10000,
	}

	bookingID, err := repo.CreateBooking(ctx, booking)
	require.NoError(t, err)
	require.NotEmpty(t, bookingID)

	_, err = repo.CreateImage(ctx, image)
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

	affected, _, err = repo.GetBooking(ctx, bookingID)
	require.NoError(t, err)
	require.Equal(t, 0, affected)
}

func TestProperties(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	ctx := context.Background()
	repo := postgres.NewRepository(db)

	properties, err := repo.Properties(ctx, reserv.PropertyFilter{})
	require.NoError(t, err)
	require.Len(t, properties, 0)

	var propertyIDs []string
	for i := 0; i < 3; i++ {
		property := reserv.Property{
			Title:              fmt.Sprintf("Test Property %d", i),
			Description:        "Test Description",
			HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			PricePerNightCents: 10000,
			Currency:           "USD",
		}
		propertyID, err := repo.CreateProperty(ctx, property)
		require.NoError(t, err)
		propertyIDs = append(propertyIDs, propertyID)
	}

	amenities := []string{"wifi", "pool"}
	for _, propertyID := range propertyIDs {
		err = repo.CreatePropertyAmenities(ctx, propertyID, amenities)
		require.NoError(t, err)
	}

	image := reserv.PropertyImage{
		PropertyID:   uuid.MustParse(propertyIDs[0]),
		HostID:       "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CloudflareID: uuid.MustParse("2e195545-8278-41a8-9d01-3c423ec71263"),
		Filename:     "test.jpg",
		CreatedAt:    time.Now(),
	}

	image2 := reserv.PropertyImage{
		PropertyID:   uuid.MustParse(propertyIDs[0]),
		HostID:       "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CloudflareID: uuid.MustParse("8cbd89cc-a87d-4cd7-9a86-3453dae882d8"),
		Filename:     "test2.jpg",
		CreatedAt:    time.Now(),
	}

	image3 := reserv.PropertyImage{
		PropertyID:   uuid.MustParse(propertyIDs[1]),
		HostID:       "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CloudflareID: uuid.MustParse("8cbd89cc-a87d-4cd7-9a86-3453dae882d8"),
		Filename:     "test2.jpg",
		CreatedAt:    time.Now(),
	}

	_, err = repo.CreateImage(ctx, image)
	require.NoError(t, err)

	_, err = repo.CreateImage(ctx, image2)
	require.NoError(t, err)

	_, err = repo.CreateImage(ctx, image3)
	require.NoError(t, err)

	properties, err = repo.Properties(ctx, reserv.PropertyFilter{})
	require.NoError(t, err)
	require.Len(t, properties, 3)

	for _, property := range properties {
		require.NotEmpty(t, property.ID)
		require.NotEmpty(t, property.HostID)
		require.NotEmpty(t, property.Title)
		require.NotEmpty(t, property.Description)
		require.NotEmpty(t, property.Currency)
		require.NotEmpty(t, property.CreatedAt)
		require.NotEmpty(t, property.UpdatedAt)
		require.NotZero(t, property.PricePerNightCents)
		require.Equal(t, int64(10000), property.PricePerNightCents)
		require.Equal(t, "USD", property.Currency)
		require.NotEmpty(t, property.Amenities)
		require.Len(t, property.Amenities, 2)
		for _, amenity := range property.Amenities {
			require.NotEmpty(t, amenity.ID)
			require.NotEmpty(t, amenity.Name)
		}

		if property.ID == uuid.MustParse(propertyIDs[0]) {
			require.NotEmpty(t, property.Images)
			require.Len(t, property.Images, 2)
			for _, image := range property.Images {
				require.NotEmpty(t, image.ID)
				require.NotEmpty(t, image.Filename)
			}
		} else if property.ID == uuid.MustParse(propertyIDs[1]) {
			require.NotEmpty(t, property.Images)
			require.Len(t, property.Images, 1)
			for _, image := range property.Images {
				require.NotEmpty(t, image.ID)
				require.NotEmpty(t, image.Filename)
			}
		} else {
			require.Empty(t, property.Images)
		}
	}
}

func TestPropertiesWithHostIDFilter(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := postgres.NewRepository(db)
	ctx := context.Background()

	properties, err := repo.Properties(ctx, reserv.PropertyFilter{HostID: "user_2x5CiRO5Mf0wBpWO8w469jEJhRq"})
	require.NoError(t, err)
	require.Len(t, properties, 0)

	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
	}

	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)
	require.NotEmpty(t, propertyID)

	properties, err = repo.Properties(ctx, reserv.PropertyFilter{HostID: "user_2x5CiRO5Mf0wBpWO8w469jEJhRq"})
	require.NoError(t, err)
	require.Len(t, properties, 1)

	property2 := reserv.Property{
		Title:              "Test Property 2",
		Description:        "Test Description 2",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
	}

	propertyID2, err := repo.CreateProperty(ctx, property2)
	require.NoError(t, err)
	require.NotEmpty(t, propertyID2)

	properties, err = repo.Properties(ctx, reserv.PropertyFilter{HostID: "user_2x5CiRO5Mf0wBpWO8w469jEJhRq"})
	require.NoError(t, err)
	require.Len(t, properties, 2)
}

func TestCreateImage(t *testing.T) {
	db := OpenDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := postgres.NewRepository(db)
	ctx := context.Background()

	// create a property
	property := reserv.Property{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	propertyID, err := repo.CreateProperty(ctx, property)
	require.NoError(t, err)

	image := reserv.PropertyImage{
		PropertyID:   uuid.MustParse(propertyID),
		HostID:       "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CloudflareID: uuid.MustParse("2e195545-8278-41a8-9d01-3c423ec71263"),
		Filename:     "test.jpg",
		CreatedAt:    time.Now(),
	}

	imageID, err := repo.CreateImage(ctx, image)
	require.NoError(t, err)
	var createdImage reserv.PropertyImage
	err = db.GetContext(ctx, &createdImage, "SELECT * FROM property_images WHERE id = $1", imageID)
	require.NoError(t, err)
	require.Equal(t, image.PropertyID, createdImage.PropertyID)
	require.Equal(t, image.CloudflareID, createdImage.CloudflareID)
	require.Equal(t, image.Filename, createdImage.Filename)
	require.Equal(t, image.HostID, createdImage.HostID)
	require.NotNil(t, createdImage.CreatedAt)
	require.NotZero(t, createdImage.CreatedAt)
}
