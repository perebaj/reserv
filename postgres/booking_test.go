//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/perebaj/reserv"
	"github.com/perebaj/reserv/postgres"
	"github.com/stretchr/testify/require"
)

func TestCreateBooking(t *testing.T) {
	db := OpenDB(t)
	defer db.Close()

	repo := postgres.NewRepository(db)

	guestID := uuid.New().String()

	property := reserv.Property{
		HostID:             guestID,
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(context.Background(), property)
	require.NoError(t, err)

	want := reserv.Booking{
		PropertyID:      propertyID,
		GuestID:         guestID,
		CheckInDate:     time.Now().AddDate(2025, 1, 1).Format("2006-01-02"),
		CheckOutDate:    time.Now().AddDate(2025, 1, 4).Format("2006-01-02"),
		TotalPriceCents: 10000,
		Currency:        "USD",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	id, err := repo.CreateBooking(context.Background(), want)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	want2 := reserv.Booking{
		PropertyID:      propertyID,
		GuestID:         guestID,
		CheckInDate:     time.Now().AddDate(2025, 1, 5).Format("2006-01-02"),
		CheckOutDate:    time.Now().AddDate(2025, 1, 14).Format("2006-01-02"),
		TotalPriceCents: 10000,
		Currency:        "USD",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	id2, err := repo.CreateBooking(context.Background(), want2)
	require.NoError(t, err)
	require.NotEmpty(t, id2)

	want3 := reserv.Booking{
		PropertyID:      propertyID,
		GuestID:         guestID,
		CheckInDate:     time.Now().AddDate(2025, 1, 14).Format("2006-01-02"),
		CheckOutDate:    time.Now().AddDate(2025, 1, 15).Format("2006-01-02"),
		TotalPriceCents: 10000,
		Currency:        "USD",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	id3, err := repo.CreateBooking(context.Background(), want3)
	require.NoError(t, err)
	require.NotEmpty(t, id3)

	// var got reserv.Booking
	// err = db.GetContext(context.Background(), &got, "SELECT * FROM bookings WHERE id = $1", id)
	// require.NoError(t, err)
	// require.Equal(t, got.ID, id)
	// require.Equal(t, got.PropertyID, want.PropertyID)
	// require.Equal(t, got.GuestID, want.GuestID)
	// require.Equal(t, got.CheckInDate, want.CheckInDate)
	// require.Equal(t, got.CheckOutDate, want.CheckOutDate)
	// require.Equal(t, got.TotalPriceCents, want.TotalPriceCents)
	// require.Equal(t, got.Currency, want.Currency)
}

func TestGetBooking(t *testing.T) {
	db := OpenDB(t)
	defer db.Close()

	repo := postgres.NewRepository(db)

	guestID := uuid.New().String()

	property := reserv.Property{
		HostID:             guestID,
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(context.Background(), property)
	require.NoError(t, err)

	booking := reserv.Booking{
		PropertyID:      propertyID,
		GuestID:         guestID,
		CheckInDate:     time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		CheckOutDate:    time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
		TotalPriceCents: 10000,
		Currency:        "USD",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	id, err := repo.CreateBooking(context.Background(), booking)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	affected, got, err := repo.GetBooking(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, affected, 1)
	require.Equal(t, got.ID, id)
	require.Equal(t, got.PropertyID, booking.PropertyID)
	require.Equal(t, got.GuestID, booking.GuestID)
	require.Equal(t, got.CheckInDate, booking.CheckInDate)
	require.Equal(t, got.CheckOutDate, booking.CheckOutDate)
}

func TestDeleteBooking(t *testing.T) {
	db := OpenDB(t)
	defer db.Close()

	repo := postgres.NewRepository(db)

	guestID := uuid.New().String()

	property := reserv.Property{
		HostID:             guestID,
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	propertyID, err := repo.CreateProperty(context.Background(), property)
	require.NoError(t, err)

	booking := reserv.Booking{
		PropertyID:   propertyID,
		GuestID:      guestID,
		CheckInDate:  time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		CheckOutDate: time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
	}

	id, err := repo.CreateBooking(context.Background(), booking)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	err = repo.DeleteBooking(context.Background(), id)
	require.NoError(t, err)

	affected, _, err := repo.GetBooking(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, affected, 0)
}
