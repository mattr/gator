-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
        values ($1, now(), now(), $2, $3)
        returning id, created_at, updated_at, user_id, feed_id)
select inserted_feed_follow.*,
       feeds.name as feed_name,
       users.name as user_name
from inserted_feed_follow
         inner join feeds on inserted_feed_follow.feed_id = feeds.id
         inner join users on inserted_feed_follow.user_id = users.id;

-- name: DeleteFeedFollow :exec
delete
from feed_follows
where feed_follows.user_id = $1
  and feed_follows.feed_id = (select feeds.id from feeds where feeds.url = $2);