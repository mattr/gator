# Gator

A command-line RSS feed aggregator.

## Installation

You will need Go (1.23.0+) to install this and Postgres (15+) installed to run this.

Install the app by running:

```bash
go install github.com/mattr/gator
```

Create a `.gatorconfig.json` file in your home directory with the following:

```json
{
  "db_url": "postgres://username@localhost:5432/gator?sslmode=disable"
}
```

replacing `username` with the postgres user to connect to the db and `gator` with the name of the database (if
different).

## Using Gator

Run the aggregator in the background:

```bash
gator agg [time_between_reqs]
```

where `time_between_reqs` is how frequently it should refresh the posts from the feeds. It accepts a time format that
can be parsed in Go, e.g. `10s`, `15m`, `1h`. The aggregator can run in the background.

You can create a new feed by running:

```bash
gator addfeed "Feed name", "https://path-to-feed"
```

When you add a feed, you are automatically subscribed to it.

You can follow an existing feed by running:

```bash
gator follow "https://path-to-feed"
```

or unfollow with:

```bash
gator unfollow "https://path-to-feed"
```

You can see your current following list by running:

```bash
gator following
```

To view the most recent posts from feeds you are following:

```bash
gator browse [limit]
```

which will fetch the `limit` most recent articles for feeds you are following (default: 2).

## Extending the Project

Some options to extend the project:
* Add sorting and filtering options to the browse command
* Add pagination to the browse command
* Add concurrency to the agg command so that it can fetch more frequently
* Add a search command that allows for fuzzy searching of posts
* Add bookmarking or liking posts
* Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
* Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
* Write a service manager that keeps the agg command running in the background and restarts it if it crashes
