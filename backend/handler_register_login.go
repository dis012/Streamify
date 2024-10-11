package main

import (
	"encoding/json"
	"net/http"

	"github.com/dis012/StreamingServer/auth"
	"github.com/dis012/StreamingServer/internal/database"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
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
	}

	var user param

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

	newUser, err := a.dbQueries.RegisterUser(r.Context(), database.RegisterUserParams{
		Email:          user.Email,
		HashedPassword: hashed_password,
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
