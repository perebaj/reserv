package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

type Image struct {
	Id       string   `json:"id"`
	FileName string   `json:"filename"`
	Variants []string `json:"variants"`
}

func (h *Handler) handlerPostImage(w http.ResponseWriter, r *http.Request) {
	slog.Info("handlePostImage")
	err := r.ParseMultipartForm(32 << 20) // 32MB is the maximum size of a file we can upload
	if err != nil {
		slog.Error("failed to parse multipart form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mForm := r.MultipartForm
	for k := range mForm.File {
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			slog.Error("failed to get image from form", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
			Metadata: map[string]interface{}{
				"upload": "api",
			},
		})
		if err != nil {
			slog.Error("failed to upload image", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(imgByte)
	}
}
