package middleware

import (
	"encoding/json"
	"net/http"
	"os"

	"server/http/response"
)

type SuperAdminBody struct {
	SecretKey string `json:"secret_key"`
}

func SupremeLeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Body
		var reqBody SuperAdminBody

		// Decode the request body into the struct
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			response.RespondeWithError(w, http.StatusUnprocessableEntity, "invalid json")
			return
		}

		// Match "SecretKey" is valid
		if reqBody.SecretKey == os.Getenv("supereme_leader_secret_key") {
			// Call the next handler with the updated context
			next.ServeHTTP(w, r)
			return
		}

		// Stop and Send Error Back
		http.Error(w, "U r not Supreme Leader Idiot", http.StatusTeapot)
	})
}
