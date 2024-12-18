package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dis012/StreamingServer/auth"
	"github.com/dis012/StreamingServer/internal/database"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type LoginUser struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (a *ApiConfig) RegisterUser(w http.ResponseWriter, r *http.Request) {
	/*
		Method of ApiConfig struct that will handle the registration of a user.
		It accepts user's email and password, and returns a JSON response.
		It saves the user's email and hashed password to the database.
	*/
	type param struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Secret   string `json:"secret"`
	}

	var user param

	isAdmin := false

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashed_password, err := auth.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Secret == a.admin_secret {
		isAdmin = true
	}

	newUser, err := a.dbQueries.RegisterUser(r.Context(), database.RegisterUserParams{
		Email:          user.Email,
		HashedPassword: hashed_password,
		IsAdmin:        isAdmin,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithJson(w, http.StatusCreated, User{
		Id:    int(newUser.ID),
		Email: newUser.Email,
	})
}

func (a *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	/*
		Method of ApiConfig struct that will handle the login of a user.
		It accepts user's email and password, and returns a JSON response.
		It checks if the user's email and password are correct, and returns a JWT token (access token) and refreash token.
	*/
	type param struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginUser param

	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := a.dbQueries.GetUserByEmail(r.Context(), loginUser.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if dbUser.ID == 0 {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	checkIfPasswordValid := auth.CompareSavedAndInputPassword(loginUser.Password, dbUser.HashedPassword)
	if !checkIfPasswordValid {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := a.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		RevokedAt: sql.NullTime{},
		Token:     refreshToken,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accesToken, err := auth.MakeAccessToken(int(dbUser.ID), a.secret, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithJson(w, http.StatusOK, LoginUser{
		Id:           int(dbUser.ID),
		Email:        dbUser.Email,
		RefreshToken: newRefreshToken.Token,
		AccessToken:  accesToken,
	})
}

func (a *ApiConfig) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	/*
		Handler function that will refresh the access token.
		It accepts the refresh token, checks if refresh token exists
		and if its valid and then returns new access token.
	*/
	bearrerToken, err := auth.GetBearrerOfTheToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, err := a.dbQueries.GetRefreshTokenByToken(r.Context(), bearrerToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	if refreshToken.RevokedAt.Valid {
		http.Error(w, "Refresh token revoked", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeAccessToken(int(refreshToken.UserID), a.secret, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithJson(w, http.StatusOK, map[string]string{
		"access_token": token,
	})
}

func (a *ApiConfig) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	/*
		Handler function that will revoke the refresh token.
	*/

	bearrerToken, err := auth.GetBearrerOfTheToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, err := a.dbQueries.GetRefreshTokenByToken(r.Context(), bearrerToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.dbQueries.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithJson(w, http.StatusNoContent, interface{}(nil))
}
