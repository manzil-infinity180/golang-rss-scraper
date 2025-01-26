package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetJobs(w http.ResponseWriter, r *http.Request) {

	jobs, err := cfg.DB.GetRemoteJobs(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get jobs for user")
		return
	}
	respondWithJson(w, http.StatusOK, databaseJobsToJobs(jobs))
}
