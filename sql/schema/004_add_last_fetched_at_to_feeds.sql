-- +goose Up
alter table feeds
add column last_fetched_at timestamp;