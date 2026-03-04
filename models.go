package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/oneelabed/RSSAggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func DBUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
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
		UserID:    dbFeed.UserID,
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
