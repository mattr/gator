-- name: CreateFeed :one
insert into feeds(id, created_at, updated_at, name, url, user_id)
values($1, now(), now(), $2, $3, $4)
returning id, created_at, updated_at, name, url, user_id;

-- name: GetFeeds :many
select id, created_at, updated_at, name, url, user_id
from feeds;