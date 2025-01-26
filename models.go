package main

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/manzil-infinity180/golang-webrss/internal/database"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	result := make([]Feed, len(feeds))
	for i, feed := range feeds {
		result[i] = databaseFeedToFeed(feed)
	}
	return result
}

type FeedFollows struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollows {
	return FeedFollows{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	}
}

func databaseFeedsFollowsToFeedsFollows(feedsFollows []database.FeedFollow) []FeedFollows {
	result := make([]FeedFollows, len(feedsFollows))
	for i, feedsFollow := range feedsFollows {
		result[i] = databaseFeedFollowToFeedFollow(feedsFollow)
	}
	return result
}

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func databasePostToPost(post database.Post) Post {
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: nullStringToStringPtr(post.Description),
		PublishedAt: nullTimeToTimePtr(post.PublishedAt),
		FeedID:      post.FeedID,
	}
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func nullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func databasePostsToPosts(posts []database.Post) []Post {
	result := make([]Post, len(posts))
	for i, post := range posts {
		result[i] = databasePostToPost(post)
	}
	return result
}

type Job struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Company     string     `json:"company"`
	Url         string     `json:"url"`
	Image       string     `json:"image"`
	Description *string    `json:"description"`
	Tag         *string    `json:"tag"`
	Location    string     `json:"location"`
	PublishedAt *time.Time `json:"published_at"`
}

func databaseJobToJob(post database.Job) Job {
	return Job{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Company:     post.Company,
		Image:       post.Image,
		Location:    post.Location,
		Url:         post.Url,
		Description: nullStringToStringPtr(post.Description),
		Tag:         nullStringToStringPtr(post.Tag),
		PublishedAt: nullTimeToTimePtr(post.PublishedAt),
	}
}

func databaseJobsToJobs(jobs []database.Job) []Job {
	result := make([]Job, len(jobs))
	for i, job := range jobs {
		result[i] = databaseJobToJob(job)
	}
	return result
}
