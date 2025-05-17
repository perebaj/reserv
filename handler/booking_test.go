package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/perebaj/reserv"
	"github.com/perebaj/reserv/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateBookingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockBookingRepo.EXPECT().CreateBooking(gomock.Any(), gomock.Any()).Return("123", nil)

	handler := NewHandler(nil, nil, mockBookingRepo)

	requestBody := CreateBooking{
		PropertyID:   "123",
		GuestID:      "456",
		CheckInDate:  "2024-01-01",
		CheckOutDate: "2024-01-02",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	ctx := clerk.ContextWithSessionClaims(req.Context(), &clerk.SessionClaims{
		RegisteredClaims: clerk.RegisteredClaims{
			Subject: "456",
		},
	})
	req = req.WithContext(ctx)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusCreated, resp.Code, rBody)

	var response map[string]string
	err = json.Unmarshal([]byte(rBody), &response)
	require.NoError(t, err)
	require.Equal(t, "123", response["id"])
}

func TestDeleteBookingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockBookingRepo.EXPECT().DeleteBooking(gomock.Any(), gomock.Any()).Return(nil)

	handler := NewHandler(nil, nil, mockBookingRepo)

	req := httptest.NewRequest(http.MethodDelete, "/bookings/123", nil)

	resp := httptest.NewRecorder()
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNoContent, resp.Code)
}

func TestGetBookingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockBookingRepo.EXPECT().GetBooking(gomock.Any(), gomock.Any()).Return(1, reserv.Booking{
		ID: "123",
	}, nil)

	handler := NewHandler(nil, nil, mockBookingRepo)

	req := httptest.NewRequest(http.MethodGet, "/bookings/123", nil)

	resp := httptest.NewRecorder()
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)
}

func TestBookingsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookingRepo := mock.NewMockBookingRepository(ctrl)
	mockBookingRepo.EXPECT().Bookings(gomock.Any(), gomock.Any()).Return([]reserv.Booking{
		{
			ID: "123",
		},
	}, nil)

	handler := NewHandler(nil, nil, mockBookingRepo)

	req := httptest.NewRequest(http.MethodGet, "/bookings?property_id=123&guest_id=456", nil)
	req.Header.Set("Authorization", "Bearer test_token")

	ctx := clerk.ContextWithSessionClaims(req.Context(), &clerk.SessionClaims{
		RegisteredClaims: clerk.RegisteredClaims{
			Subject: "456",
		},
	})
	req = req.WithContext(ctx)

	resp := httptest.NewRecorder()
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)
}
