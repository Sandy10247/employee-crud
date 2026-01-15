package adminhandler

import (
	"errors"
	"fmt"
	"net/http"

	"server/http/middleware"
	"server/http/response"

	db "server/init"

	"github.com/jackc/pgx/v5"
)

// Make the user Admin, if they are so Remove then as Admin
func MakeBreak(w http.ResponseWriter, r *http.Request) {
	// Extract UserInfo from context
	userInfo, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		response.RespondeWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	// check whether user is Admin
	_, err := db.Queries.GetAdminUser(r.Context(), userInfo.ID)
	// User is not Admin
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		_, err := db.Queries.CreateAdminUser(r.Context(), userInfo.ID)
		if err != nil {
			response.RespondeWithError(w, http.StatusBadRequest, "Admin User Issue")
			return
		}

		response.RespondeWithJSON(w, http.StatusCreated, "New Admin Added")
		return
	}
	if err != nil {
		response.RespondeWithError(w, http.StatusBadRequest, fmt.Sprintf("Admin User Issue %v", err))
		return
	}

	// user exists, Remove user
	_, err = db.Queries.DeleteAdminUser(r.Context(), userInfo.ID)
	if err != nil {
		response.RespondeWithError(w, http.StatusBadRequest, "Admin User Issue")
		return
	}

	response.RespondeWithJSON(w, http.StatusOK, "Removed Existing User")
}
