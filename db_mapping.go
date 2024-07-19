package main

import (
	"time"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
	"github.com/google/uuid"
)


type Feed struct {
	ID            uuid.UUID    `json:"id"`
	Name          string       `json:"name"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Url           string       `json:"url"`
	UserID        uuid.UUID    `json:"user_id"`
	LastFetchedAt *time.Time   `json:"last_fetched_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID: feedFollow.ID,
		FeedID: feedFollow.FeedID,
		UserID: feedFollow.UserID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
	}
}

func databaseFeedToFeed(feed database.Feed) Feed {
	var timePtr *time.Time
	if feed.LastFetchedAt.Valid {
		timePtr = &feed.LastFetchedAt.Time
	}
	return Feed{
		ID: feed.ID,
		Url: feed.Url,
		Name: feed.Name,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		UserID: feed.UserID,
		LastFetchedAt: timePtr,
	}

}

func databaseUserToUser(user database.User) User {
	return User {
		ID: user.ID,
		Name: user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ApiKey: user.ApiKey,
	}
}
