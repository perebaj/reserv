package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/google/uuid"
	"github.com/perebaj/reserv"
)

//go:generate mockgen -source images.go -destination ../mock/images.go -package mock

// CloudFlareAPI is the interface for the CloudFlare API.
type CloudFlareAPI interface {
	UploadImage(ctx context.Context, container *cloudflare.ResourceContainer, params cloudflare.UploadImageParams) (cloudflare.Image, error)
}

// Image is the response body for getting an image.
type Image struct {
	ID       string   `json:"id"`
	FileName string   `json:"filename"`
	Variants []string `json:"variants"`
}

// CreatePropertyImage is the request body for creating a property image.
type CreatePropertyImage struct {
	PropertyID string `json:"property_id"`
	HostID     string `json:"host_id"`
}

func (h *Handler) handlerPostImage(w http.ResponseWriter, r *http.Request) {
	slog.Info("handlePostImage")

	err := r.ParseMultipartForm(32 << 20) // 32MB is the maximum size of a file we can upload
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err.Error())
		NewAPIError("failed to parse multipart form", "failed to parse multipart form", http.StatusInternalServerError).Write(w)
		return
	}
	mForm := r.MultipartForm

	propertyID, ok := mForm.Value["property_id"]
	if !ok {
		slog.Error("property_id is required")
		NewAPIError("property_id is required", "property_id is required", http.StatusBadRequest).Write(w)
		return
	}
	hostID, ok := mForm.Value["host_id"]
	if !ok {
		slog.Error("host_id is required")
		NewAPIError("host_id is required", "host_id is required", http.StatusBadRequest).Write(w)
		return
	}

	propertyIDStr := propertyID[0]
	hostIDStr := hostID[0]

	if propertyIDStr == "" || hostIDStr == "" {
		slog.Error("property_id and host_id are required and must be valid UUIDs")
		NewAPIError("property_id and host_id are required and must be valid UUIDs", "property_id and host_id are required and must be valid UUIDs", http.StatusBadRequest).Write(w)
		return
	}

	var cloudflareID string
	var filename string
	for k := range mForm.File {
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			slog.Error("failed to get image from form", "error", err.Error())
			NewAPIError("failed to get image from form", "failed to get image from form", http.StatusInternalServerError).Write(w)
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				slog.Error("failed to close file", "error", err.Error())
			}
		}()
		// TODO(@perebaj): Remove this variable from here. Just doing that because injecting the ACCOUNT_ID using the SDK is not working.
		accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
		if accountID == "" {
			slog.Error("CLOUDFLARE_ACCOUNT_ID is not set")
			NewAPIError("CLOUDFLARE_ACCOUNT_ID is not set", "CLOUDFLARE_ACCOUNT_ID is not set", http.StatusInternalServerError).Write(w)
			return
		}
		img, err := h.CloudFlare.UploadImage(r.Context(), &cloudflare.ResourceContainer{
			Identifier: accountID,
			Level:      cloudflare.AccountRouteLevel,
		}, cloudflare.UploadImageParams{
			File: file,
			Name: fileHeader.Filename,
		})

		filename = img.Filename
		cloudflareID = img.ID

		if err != nil {
			slog.Error("failed to upload image", "error", err.Error())
			NewAPIError("failed to upload image", "failed to upload image", http.StatusInternalServerError).Write(w)
			return
		}
		slog.Info("uploaded image", "id", img.ID, "filename", img.Filename)
	}

	id, err := h.repo.CreateImage(r.Context(), reserv.PropertyImage{
		PropertyID:   uuid.MustParse(propertyIDStr),
		HostID:       "user_2x5CiRO5Mf0wBpWO8w469jEJhRq",
		CloudflareID: uuid.MustParse(cloudflareID),
		Filename:     filename,
		CreatedAt:    time.Now(),
	})
	slog.Info("image created on postgres", "id", id)
	if err != nil {
		slog.Error("failed to create image", "error", err.Error())
		NewAPIError("failed to create image", "failed to create image", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
