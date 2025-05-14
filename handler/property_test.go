package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	uid := "user_2x5CiRO5Mf0wBpWO8w469jEJhRq"
	payload := handler.CreatePropertyRequest{
		Title:              "Test Property",
		Description:        "Test Description",
		PricePerNightCents: 10000,
		Currency:           "USD",
		HostID:             uid,
	}
	repo.EXPECT().CreateProperty(gomock.Any(), gomock.Any()).Return(uid, nil)

	jsonBody, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/properties", bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	mux := http.NewServeMux()
	propHandler := handler.NewHandler(repo, nil)
	propHandler.RegisterRoutes(mux)
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

	req := httptest.NewRequest(http.MethodPut, "/properties/"+propertyID, bytes.NewBuffer(jsonBody))
	resp := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	mux := http.NewServeMux()
	handler := handler.NewHandler(repo, nil)
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

	req := httptest.NewRequest(http.MethodDelete, "/properties/"+propertyID, nil)
	resp := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler := handler.NewHandler(repo, nil)
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

	req := httptest.NewRequest(http.MethodGet, "/properties/"+propertyID.String(), nil)
	resp := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler := handler.NewHandler(repo, nil)
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

func TestGetProperty_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockPropertyRepository(ctrl)
	repo.EXPECT().GetProperty(gomock.Any(), gomock.Any()).Return(0, reserv.Property{}, nil)
	propertyID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/properties/"+propertyID, nil)
	resp := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler := handler.NewHandler(repo, nil)
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusNotFound, resp.Code, rBody)
}

func TestGetAmenities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockPropertyRepository(ctrl)
	repo.EXPECT().Amenities(gomock.Any()).Return([]reserv.Amenity{
		{ID: "1", Name: "Amenity 1", CreatedAt: time.Now()},
		{ID: "2", Name: "Amenity 2", CreatedAt: time.Now()},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/amenities", nil)
	resp := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler := handler.NewHandler(repo, nil)
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusOK, resp.Code, rBody)

	var response []reserv.Amenity
	err := json.Unmarshal([]byte(rBody), &response)
	require.NoError(t, err)
	require.Equal(t, 2, len(response))
	require.Equal(t, "Amenity 1", response[0].Name)
	require.Equal(t, "Amenity 2", response[1].Name)
	require.Equal(t, "1", response[0].ID)
	require.Equal(t, "2", response[1].ID)
}

func TestCreatePropertyAmenity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockPropertyRepository(ctrl)
	repo.EXPECT().CreatePropertyAmenities(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	propertyID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPost, "/properties/"+propertyID+"/amenities", bytes.NewBuffer([]byte(`["1", "2"]`)))
	resp := httptest.NewRecorder()

	mux := http.NewServeMux()
	handler := handler.NewHandler(repo, nil)
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	rBody := resp.Body.String()
	require.Equal(t, http.StatusOK, resp.Code, rBody)
}
