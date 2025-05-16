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
		CheckInDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CheckOutDate:    time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
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
		CheckInDate:     time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		CheckOutDate:    time.Date(2025, 1, 14, 0, 0, 0, 0, time.UTC),
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
		CheckInDate:     time.Date(2025, 1, 14, 0, 0, 0, 0, time.UTC),
		CheckOutDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		TotalPriceCents: 10000,
		Currency:        "USD",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	id3, err := repo.CreateBooking(context.Background(), want3)
	require.Error(t, err, "booking overlaps")
	require.Empty(t, id3)

	want4 := reserv.Booking{
		PropertyID:   propertyID,
		GuestID:      guestID,
		CheckInDate:  time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		CheckOutDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	id4, err := repo.CreateBooking(context.Background(), want4)
	require.Error(t, err, "booking overlaps")
	require.Empty(t, id4)

	want5 := reserv.Booking{
		PropertyID:   propertyID,
		GuestID:      guestID,
		CheckInDate:  time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		CheckOutDate: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
	}

	id5, err := repo.CreateBooking(context.Background(), want5)
	require.NoError(t, err)
	require.NotEmpty(t, id5)
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
		CheckInDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CheckOutDate:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
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
		CheckInDate:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CheckOutDate: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
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
