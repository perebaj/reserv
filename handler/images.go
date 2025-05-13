package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/google/uuid"
	"github.com/perebaj/reserv"
)

type Image struct {
	Id       string   `json:"id"`
	FileName string   `json:"filename"`
	Variants []string `json:"variants"`
}

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

	propertyID := mForm.Value["property_id"][0]
	hostID := mForm.Value["host_id"][0]

	if propertyID == "" || hostID == "" {
		slog.Error("property_id and host_id are required")
		NewAPIError("property_id and host_id are required", "property_id and host_id are required", http.StatusBadRequest).Write(w)
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
		defer file.Close()
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
		imgByte, err := json.Marshal(Image{
			Id:       img.ID,
			FileName: img.Filename,
			Variants: img.Variants,
		})
		if err != nil {
			slog.Error("failed to marshal image", "error", err.Error())
			NewAPIError("failed to marshal image", "failed to marshal image", http.StatusInternalServerError).Write(w)
			return
		}
		_, _ = w.Write(imgByte)
	}

	_, err = h.repo.CreateImage(r.Context(), reserv.PropertyImage{
		PropertyID:   uuid.MustParse(propertyID),
		HostID:       uuid.MustParse(hostID),
		CloudflareID: uuid.MustParse(cloudflareID),
		Filename:     filename,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		slog.Error("failed to create image", "error", err.Error())
		NewAPIError("failed to create image", "failed to create image", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
