package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

// func publicRoute(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(`{"access": "public"}`))
// }

// func protectedRoute(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := clerk.SessionClaimsFromContext(r.Context())
// 	if !ok {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte(`{"access": "unauthorized"}`))
// 		return
// 	}
// 	fmt.Fprintf(w, `{"user_id": "%s"}`, claims.Subject)
// }

type createUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// createUserResponse wraps the clerk.User to be used on this application.
type createUserResponse struct {
	clerk.User
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("error decoding request", "error", err)
		NewAPIError("bad_request", "invalid request", http.StatusBadRequest).Write(w)
		return
	}

	skipPasswordRequirement := true
	createOrganizationEnabled := true
	usr, err := user.Create(r.Context(), &user.CreateParams{
		EmailAddresses:            &[]string{req.Email},
		FirstName:                 &req.FirstName,
		LastName:                  &req.LastName,
		SkipPasswordRequirement:   &skipPasswordRequirement,
		CreateOrganizationEnabled: &createOrganizationEnabled,
	})

	if apiErr, ok := err.(*clerk.APIErrorResponse); ok {
		// create a function to check if the error is an email_address_already_exists error
		if isEmailAddressAlreadyExists(apiErr) {
			slog.Error("email already in use", "error", err)
			NewAPIError("email_already_in_use", "email already in use", http.StatusBadRequest).Write(w)
			return
		}
	}

	if err != nil {
		slog.Error("error creating user", "error", err)
		NewAPIError("internal_server_error", "error creating user", http.StatusInternalServerError).Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	var respUser createUserResponse
	respUser.User = *usr
	json.NewEncoder(w).Encode(respUser)
}

func isEmailAddressAlreadyExists(apiErr *clerk.APIErrorResponse) bool {
	for _, err := range apiErr.Errors {
		if err.Code == "form_identifier_exists" {
			return true
		}
	}
	return false
}

// Router centralize all routes for the API and returns a ServeMux to be used in the server initialization.
func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /user", createUser)
	return mux
}
