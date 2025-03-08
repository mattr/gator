package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mattr/gator/internal/database"
)

// handlerLogin provides a handler function for selecting a user as the current user.
// If the username provided in the cmd.args is not present in the database, returns an error,
// otherwise sets the configuration to use the specified user.
//
// Invoked with the 'login' argument
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login handler expects a single argument (username)")
	}
	user, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("Logged in as " + s.config.CurrentUserName)
	return nil
}

// handlerRegister registers a new user in the database with the given username.
// The username must be unique; returns an error if another user with the same name already exists.
//
// Invoked with the register argument
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 || len(cmd.args[0]) == 0 {
		return errors.New("register handler expects a single argument (username)")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), Name: cmd.args[0]})
	if err != nil {
		return err
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("Created user ", user.Name)
	fmt.Printf("%v\n", user)
	return nil
}

// handlerReset provides a reset function to remove all users (and by cascade, all feeds and follows).
//
// Invoked with the reset argument
func handlerReset(s *state, cmd command) error {
	return s.db.DeleteAllUsers(context.Background())
}

// handlerUsers lists all the users currently registered in the system
//
// Invoked with the users argument
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		current := ""
		if user.Name == s.config.CurrentUserName {
			current = "(current)"
		}
		fmt.Printf("* %s %s\n", user.Name, current)
	}
	return nil
}

// handlerFeed fetches a feed from a URL and stores it's configuration in the database.
//
// Invoked with the agg argument
func handlerFeed(s *state, cmd command) error {
	//if len(cmd.args) == 0 {
	//	return errors.New("feed handler expects a single argument (url)")
	//}
	//feedURL := cmd.args[0]
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", *feed)
	return nil
}

// handlerFeeds lists all feeds currently stored in the database
//
// Invoked with the feeds argument
func handlerFeeds(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		user, err := getUser(users, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

// handlerAddFeed adds a new feed to the database. The current user is stored as the
// creator.
//
// Invoked with the addfeed argument.
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("add feed handler expects two arguments (name and url)")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feedParams := database.CreateFeedParams{ID: uuid.New(), Name: name, Url: url, UserID: user.ID}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	followParams := database.CreateFeedFollowParams{ID: uuid.New(), UserID: user.ID, FeedID: feed.ID}
	_, err = s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", feed)
	return nil
}

// handlerFeedFollow follows a feed specified by URL for the current user.
//
// Invoked with the follow argument
func handlerFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("follow handler expects a single argument (url)")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{ID: uuid.New(), UserID: user.ID, FeedID: feed.ID}
	follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("%s is following '%s'\n", follow.UserName, follow.FeedName)
	return nil
}

// handlerFeedFollowing lists all the feeds that the current user is following
//
// Invoked with the following argument.
func handlerFeedFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s is following:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("* %s '%s'\n", feed.Name, feed.Url)
	}
	return nil
}

// handlerFeedUnfollow removes a feed from the current user's follows
//
// Invoked with the unfollow argument.
func handlerFeedUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("unfollow handler expects a single argument (url)")
	}

	params := database.DeleteFeedFollowParams{UserID: user.ID, Url: cmd.args[0]}
	return s.db.DeleteFeedFollow(context.Background(), params)
}
