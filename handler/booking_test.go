package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		CheckInDate:  "2024-01-01T00:00:00Z",
		CheckOutDate: "2024-01-02T00:00:00Z",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

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
