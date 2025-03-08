-- +goose Up
create table posts (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text not null,
    url text not null unique,
    description text,
    published_at timestamp,
    feed_id uuid not null references feeds
);

-- +goose Down
drop table posts;
