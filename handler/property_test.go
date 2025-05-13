package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/perebaj/reserv"
	"github.com/perebaj/reserv/handler"
	"github.com/perebaj/reserv/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateProperty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockPropertyRepository(ctrl)

	uid := uuid.New().String()
	payload := handler.CreatePropertyRequest{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             uuid.MustParse(uid),
	}
	repo.EXPECT().CreateProperty(gomock.Any(), gomock.Any()).Return(uid, nil)

	jsonBody, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/properties", bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	mux := http.NewServeMux()
	handler := handler.NewPropertyHandler(repo)
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusCreated, resp.Code, rBody)

	var response map[string]string
	err = json.Unmarshal([]byte(rBody), &response)
	require.NoError(t, err)
	require.Equal(t, uid, response["id"])
}

func TestUpdateProperty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockPropertyRepository(ctrl)

	payload := handler.UpdatePropertyRequest{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
	}

	propertyID := uuid.New().String()
	repo.EXPECT().UpdateProperty(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	jsonBody, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/properties?id="+propertyID, bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	mux := http.NewServeMux()
	handler := handler.NewPropertyHandler(repo)
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusNoContent, resp.Code, rBody)
}

func TestDeleteProperty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockPropertyRepository(ctrl)

	propertyID := uuid.New().String()
	repo.EXPECT().DeleteProperty(gomock.Any(), gomock.Any()).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/properties?id="+propertyID, nil)
	resp := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler := handler.NewPropertyHandler(repo)
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusNoContent, resp.Code, rBody)
}

func TestGetProperty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockPropertyRepository(ctrl)

	propertyID := uuid.New()
	repo.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Return(1, reserv.Property{ID: propertyID, Title: "Test Property", Description: "Test Description", PricePerNightCents: 10000, Currency: "USD"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/properties?id="+propertyID.String(), nil)
	resp := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler := handler.NewPropertyHandler(repo)
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusOK, resp.Code, rBody)

	var response reserv.Property
	err := json.Unmarshal([]byte(rBody), &response)
	require.NoError(t, err)
	require.Equal(t, propertyID, response.ID)
	require.Equal(t, "Test Property", response.Title)
	require.Equal(t, "Test Description", response.Description)
	require.Equal(t, int64(10000), response.PricePerNightCents)
	require.Equal(t, "USD", response.Currency)
}
