package handlers

import (
	"net/http"

	"github.com/google/uuid"
	. "github.com/oneelabed/IsraelConflictMonitor/internal/config"
	"github.com/oneelabed/IsraelConflictMonitor/internal/database"
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

func HandlerCheckNewPosts(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	latestIDStr := r.URL.Query().Get("latest_id")
	latestID, err := uuid.Parse(latestIDStr)
	if err != nil {
		RespondWithError(w, 400, "Invalid ID")
		return
	}

	hasNew, err := apiCfg.DB.CheckNewPosts(r.Context(), database.CheckNewPostsParams{
		UserID: user.ID,
		ID:     latestID,
	})

	RespondWithJSON(w, 200, map[string]bool{"hasNew": hasNew})
}
