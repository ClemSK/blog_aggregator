package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ClemSK/blog_aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create feed: %v", err))
		return
	}

	respondWithJson(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not get feed: %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerDeleteFeed(w http.ResponseWriter, r *http.Request) {
	feedIDString := chi.URLParam(r, "feedID")
	feedID, err := uuid.Parse(feedIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not parse feed id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeed(r.Context(), feedID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not delete feed : %v", err))
		return
	}

	respondWithJson(w, http.StatusOK, struct{}{})
}
