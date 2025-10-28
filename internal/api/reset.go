package api

import (
	"net/http"
)

func (h *Handler) ResetHits(w http.ResponseWriter, r *http.Request) {
	if h.Config.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	} else {
		err := h.Config.Database.DeleteAllUsers(r.Context())
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Couldn't delete users", err)
			return
		}

		h.Config.FileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0"))

	}

}
