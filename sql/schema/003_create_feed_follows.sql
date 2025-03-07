-- +goose Up
create table feed_follows(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null references users on delete cascade,
    feed_id uuid not null references feeds on delete cascade
);
create unique index feed_follows_user_feed_idx on feed_follows (user_id, feed_id);

-- +goose Down
drop table feed_follows;
