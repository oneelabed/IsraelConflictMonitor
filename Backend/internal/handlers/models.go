package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/oneelabed/IsraelConflictMonitor/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"name"`
	Password  string    `json:"password_hash"`
	APIKey    string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	IconUrl   string    `json:"icon_url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
	FeedName    string    `json:"feed_name"`
	FeedIcon    string    `json:"feed_icon"`
}

func DBUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Username:  dbUser.Username,
		Password:  dbUser.PasswordHash,
		APIKey:    dbUser.ApiKey,
	}
}

func DBFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		IconUrl:   dbFeed.IconUrl,
	}
}

func DBFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}

	for _, feed := range dbFeeds {
		feeds = append(feeds, DBFeedToFeed(feed))
	}

	return feeds
}

func DBFollowToFollow(dbFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        dbFollow.ID,
		CreatedAt: dbFollow.CreatedAt,
		UpdatedAt: dbFollow.UpdatedAt,
		UserID:    dbFollow.UserID,
		FeedID:    dbFollow.FeedID,
	}
}

func DBFollowsToFollows(dbFollows []database.FeedFollow) []FeedFollow {
	follows := []FeedFollow{}

	for _, follow := range dbFollows {
		follows = append(follows, DBFollowToFollow(follow))
	}

	return follows
}

func DBPostToPost(dbPost database.Post) Post {
	var desc *string
	if dbPost.Description.Valid {
		desc = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: desc,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}
}

func DBPostsToPosts(dbPosts []database.Post) []Post {
	posts := []Post{}

	for _, post := range dbPosts {
		posts = append(posts, DBPostToPost(post))
	}

	return posts
}

// New helper for the JOIN query results
func DBPostRowToPost(dbPost database.GetPostsForUserRow) Post {
	var desc *string
	if dbPost.Description.Valid {
		desc = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: desc,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
		FeedName:    dbPost.FeedName,
		FeedIcon:    dbPost.FeedIcon,
	}
}

func DBPostRowsToPosts(dbPosts []database.GetPostsForUserRow) []Post {
	posts := []Post{}
	for _, post := range dbPosts {
		posts = append(posts, DBPostRowToPost(post))
	}
	return posts
}

// New helper for the JOIN query results
func DBDiverseRowToPost(dbPost database.GetDiversePostsRow) Post {
	var desc *string
	if dbPost.Description.Valid {
		desc = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: desc,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
		FeedName:    dbPost.FeedName,
		FeedIcon:    dbPost.FeedIcon,
	}
}

func DBDiverseRowsToPosts(dbPosts []database.GetDiversePostsRow) []Post {
	posts := []Post{}
	for _, post := range dbPosts {
		posts = append(posts, DBDiverseRowToPost(post))
	}
	return posts
}

// Converter for the Search Query Rows
func DBSearchRowToPost(dbPost database.SearchPostsForUserRow) Post {
	var desc *string
	if dbPost.Description.Valid {
		desc = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: desc,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
		FeedName:    dbPost.FeedName,
		FeedIcon:    dbPost.FeedIcon,
	}
}

// Slice converter for Search Results
func DBSearchRowsToPosts(dbPosts []database.SearchPostsForUserRow) []Post {
	posts := make([]Post, len(dbPosts))
	for i, post := range dbPosts {
		posts[i] = DBSearchRowToPost(post)
	}
	return posts
}
