package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	. "github.com/oneelabed/IsraelConflictMonitor/internal/config"
	"github.com/oneelabed/IsraelConflictMonitor/internal/database"
)

func HandlerCreateUser(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Username string `json:"username"`
		Password string `json:"password_hash"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameter{}

	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, 500, "Couldn't hash password")
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Username:     params.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// "23505" is the code for unique_violation (username already exists)
			if pqErr.Code == "23505" {
				RespondWithError(w, http.StatusConflict, "Username already taken")
				return
			}
		}
		RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	RespondWithJSON(w, 201, DBUserToUser(user))
}

func HandlerLogin(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Invalid request payload")
		return
	}

	user, err := apiCfg.DB.GetUserByUsername(r.Context(), params.Username)
	if err != nil {
		RespondWithError(w, 401, "Invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))
	if err != nil {
		RespondWithError(w, 401, "Invalid username or password")
		return
	}

	RespondWithJSON(w, 200, DBUserToUser(user))
}

func HandlerGetUserByAPI(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJSON(w, 200, DBUserToUser(user))
}

func HandlerGetPostsForUser(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Failed to get posts: %v", err))
		return
	}

	RespondWithJSON(w, 200, DBPostRowsToPosts(posts))
}

func HandlerSearchPosts(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	queryParam := r.URL.Query().Get("q")
	param2 := sql.NullString{
		String: queryParam,
		Valid:  true,
	}

	if queryParam == "" {
		RespondWithError(w, 400, "Search query 'q' is required")
		return
	}

	posts, err := apiCfg.DB.SearchPostsForUser(r.Context(), database.SearchPostsForUserParams{
		UserID:  user.ID,
		Column2: param2,
	})
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Search failed: %v", err))
		return
	}

	RespondWithJSON(w, 200, DBSearchRowsToPosts(posts))
}

func HandlerGetAllUsers(apiCfg *ApiConfig, w http.ResponseWriter, r *http.Request, user database.User) {
	users, err := apiCfg.DB.GetAllUsers(r.Context())
	if err != nil {
		RespondWithError(w, 400, "Could not fetch users")
		return
	}

	RespondWithJSON(w, 200, DBUsersToUsers(users))
}
