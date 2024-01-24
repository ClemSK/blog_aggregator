package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ClemSK/blog_aggregator/internal/auth"
	"github.com/ClemSK/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	respondWithJson(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not get user: %v", err))
		return
	}

	respondWithJson(w, 200, databaseUserToUser(user))
}
