package main

import "net/http"

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	} else {
		err := cfg.database.DeleteAllUsers(r.Context())
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't delete users", err)
			return
		}

		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0"))

	}

}
