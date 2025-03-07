-- +goose Up
create table users (
    id uuid primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    name text not null unique default ''
);

-- +goose Down
drop table users;