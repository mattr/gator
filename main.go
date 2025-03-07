package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/mattr/gator/internal/config"
	"github.com/mattr/gator/internal/database"
	"html"
	"io"
	"log"
	"net/http"
	"os"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	available map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.available[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f := c.available[cmd.name]
	if f == nil {
		return fmt.Errorf("command %q not found", cmd.name)
	}
	return f(s, cmd)
}

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

func handlerReset(s *state, cmd command) error {
	return s.db.DeleteAllUsers(context.Background())
}

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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("add feed handler expects two arguments (name and url)")
	}

	user, err := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return err
	}

	name := cmd.args[0]
	url := cmd.args[1]

	params := database.CreateFeedParams{ID: uuid.New(), Name: name, Url: url, UserID: user.ID}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", feed)
	return nil
}

func getUser(users []database.User, id uuid.UUID) (*database.User, error) {
	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

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

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := &state{
		config: &cfg,
		db:     database.New(db),
	}

	c := commands{available: make(map[string]func(*state, command) error, 3)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerFeed)
	c.register("addfeed", handlerAddFeed)
	c.register("feeds", handlerFeeds)

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatal("Insufficient arguments provided")
	}
	cmd := command{name: userArgs[1], args: userArgs[2:]}
	err = c.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
