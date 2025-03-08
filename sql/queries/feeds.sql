-- name: CreateFeed :one
insert into feeds(id, created_at, updated_at, name, url, user_id)
values ($1, now(), now(), $2, $3, $4)
returning id, created_at, updated_at, name, url, user_id;

-- name: GetFeeds :many
select id, created_at, updated_at, name, url, user_id
from feeds;

-- name: GetFeedByURL :one
select id, created_at, updated_at, name, url, user_id
from feeds
where url = $1;

-- name: GetFeedsForUser :many
select feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, feeds.user_id
from feeds
         inner join feed_follows on feed_follows.feed_id = feeds.id
         inner join users on users.id = feed_follows.user_id
where users.id = $1;
