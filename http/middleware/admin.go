package middleware

import (
	"net/http"

	db "server/init"
)

// Key type for context values
type (
	adminUserCtx string
)

const (
	// UserIDKey is the key for user ID in the request context
	AdminCtx adminUserCtx = "admin"
)

func CheckAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "CheckAdminMiddleware :- GetUserFromContext Issue ", http.StatusInternalServerError)
			return
		}

		// Fetch Admin User
		_, err := db.Queries.GetAdminUser(r.Context(), userInfo.ID)
		if err != nil {
			http.Error(w, "Not Admin Idiot", http.StatusUnauthorized)
			return
		}

		// Call the next handler with the updated context
		next.ServeHTTP(w, r)
	})
}
