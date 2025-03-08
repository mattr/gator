-- name: CreatePost :one
insert into posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
values ($1, now(), now(), $2, $3, $4, $5, $6)
returning id, created_at, updated_at, title, url, description, published_at, feed_id;

-- name: GetPostsForUser :many
select posts.id,
       posts.created_at,
       posts.updated_at,
       posts.title,
       posts.url,
       posts.description,
       posts.published_at,
       posts.feed_id,
       feeds.name as feed_name
from posts
         inner join feeds on feeds.id = posts.feed_id
where posts.feed_id in (select feed_follows.feed_id from feed_follows where feed_follows.user_id = $1)
order by published_at desc
limit $2;
