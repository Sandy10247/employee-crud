package userhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"server/http/helper"
	"server/http/middleware"
	"server/http/response"
	"server/sql/database"

	db "server/init"
)

func HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type userReqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	var reqBody userReqBody

	// Decode the request body into the struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		response.RespondeWithError(w, 400, "invalid json")
		return
	}

	hashed, err := helper.HashPassword(reqBody.Password)
	if err != nil {
		response.RespondeWithError(w, 400, "internal server err")
		return
	}

	user, err := db.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email:        reqBody.Email,
		PasswordHash: hashed,
		Username:     reqBody.Username,
	})
	if err != nil {
		response.RespondeWithError(w, 400, fmt.Sprintf("Couldnot create user %v", err))
		return
	}

	token, err := helper.CreateToken(int64(user.ID), user.Email, user.Username)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SetJWTToken(w, "jwt", token)

	response.RespondeWithJSON(w, 200, dbuserToUser(user))
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type userReqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var reqBody userReqBody

	// Decode the request body into the struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		response.RespondeWithError(w, 400, "invalid json")
		return
	}

	user, err := db.Queries.GetUserByName(r.Context(), reqBody.Username)
	if err != nil {
		response.RespondeWithError(w, 400, err.Error())
		return
	}

	if !helper.CheckPasswordHash(reqBody.Password, user.PasswordHash) {
		response.RespondeWithError(w, 400, "invalid credentials")
		return
	}

	log.Println("user id", user.ID, "user email", user.Email, "user username", user.Username)

	token, err := helper.CreateToken(int64(user.ID), user.Email, user.Username)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SetJWTToken(w, "jwt", token)

	response.RespondeWithJSON(w, 200, dbuserToUser(user))
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	helper.UnsetJWTToken(w, "jwt")
	response.RespondeWithJSON(w, 200, map[string]interface{}{
		"Status":  true,
		"message": "Logged out successfully",
	})
}

func CheckStatus(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		response.RespondeWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	// response.RespondeWithJSON(w, 200, dbuserToUser(user))
	response := map[string]interface{}{
		"status":   true,
		"message":  "User info retrieved successfully",
		"email":    userInfo.Email,
		"username": userInfo.Username,
		"id":       userInfo.ID,
	}

	json.NewEncoder(w).Encode(response)
}
