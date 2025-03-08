package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mattr/gator/internal/database"
	"html"
	"io"
	"net/http"
)

// fetchFeed fetches the feed from the provided URL and creates a new RSSFeed object
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "Gator")
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}
	feed := &RSSFeed{}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(data, feed)
	if err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return feed, nil
}

// getUser is a filter function which takes a slice of users and an ID, and returns
// the first user that matches.
func getUser(users []database.User, id uuid.UUID) (*database.User, error) {
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}
