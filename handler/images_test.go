package handler

import (
	"bytes"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cloudflare/cloudflare-go"
	"github.com/perebaj/reserv/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_handlerPostImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// @TODO(@perebaj): Its not cool to have a environment variable in the tests.
	// Doing this because injecting the ACCOUNT_ID using the SDK is not working.
	_ = os.Setenv("CLOUDFLARE_ACCOUNT_ID", "123")

	cloudFlareMock := mock.NewMockCloudFlareAPI(ctrl)
	cloudFlareMock.EXPECT().UploadImage(gomock.Any(), gomock.Any(), gomock.Any()).Return(cloudflare.Image{
		ID:       "8a28a876-66e0-4fd6-abe1-cf3b8f3a2ab0",
		Filename: "image.jpg",
	}, nil)
	repoMock := mock.NewMockPropertyRepository(ctrl)
	repoMock.EXPECT().CreateImage(gomock.Any(), gomock.Any()).Return("61e3ecbd-9ea5-4b8e-994e-ebd00f77ec73", nil)
	handler := &Handler{
		repo:       repoMock,
		CloudFlare: cloudFlareMock,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("property_id", "eead76e2-7b39-440a-8bf8-ba78be330994")
	require.NoError(t, err)

	err = writer.WriteField("host_id", "7ad31150-c2f8-41bd-94c1-08663939d742")
	require.NoError(t, err)

	part, err := writer.CreateFormFile("file", "image.jpg")
	require.NoError(t, err)

	f, err := os.Open("testdata/image.png")
	require.NoError(t, err)
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("failed to close file", "error", err)
		}
	}()

	_, err = io.Copy(part, f)
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/images", body)
	resp := httptest.NewRecorder()
	req.Header.Set("accept", "*/*")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.ServeHTTP(resp, req)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp.Code, string(bodyBytes))
}
