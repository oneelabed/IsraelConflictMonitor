package handlers

import (
	"net/http"

	. "github.com/oneelabed/IsraelConflictMonitor/internal/config"
)

func HandlerGetDiversePosts(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	// We'll grab the top 50 diverse posts for the initial load
	posts, err := apiCfg.DB.GetDiversePosts(r.Context(), 50)
	if err != nil {
		RespondWithError(w, 400, "Couldn't get posts")
		return
	}

	// Use the Row converter we wrote earlier to include icons/names
	RespondWithJSON(w, 200, DBDiverseRowsToPosts(posts))
}
