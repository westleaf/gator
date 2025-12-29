package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/westleaf/gator/internal/database"
)

func parseUnknownTime(s string) (time.Time, error) {
	var layouts = []string{
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02 Jan 2006 15:04",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("no matching layout for %q", s)
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feedToFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("\n=== Fetching feed: %s ===\n", feedToFetch.Name)

	feeds, err := fetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		fmt.Printf("Error fetching feed %s: %v\n", feedToFetch.Name, err)
		// Still mark as fetched so we don't get stuck on broken feeds
		err = s.db.MarkFeedFetched(ctx, feedToFetch.ID)
		if err != nil {
			return err
		}
		return nil
	}

	for _, item := range feeds.Channel.Item {
		pubTime, err := parseUnknownTime(item.PubDate)
		validTime := err == nil
		if err != nil {
			fmt.Printf("Could not parse publish time, using current time: %v\n", err)
			pubTime = time.Now().UTC()
			validTime = true
		}
		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			PublishedAt: sql.NullTime{
				Time:  pubTime,
				Valid: validTime,
			},
			Title: item.Title,
			Url:   item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			FeedID: feedToFetch.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			} else {
				return err
			}
		}
	}

	err = s.db.MarkFeedFetched(ctx, feedToFetch.ID)
	if err != nil {
		return err
	}

	return nil
}
