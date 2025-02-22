package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/manzil-infinity180/golang-webrss/internal/database"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respondWithJson(w, http.StatusOK, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, http.StatusOK, databaseUserToUser(user))
}
