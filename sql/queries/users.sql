-- name: CreateUser :one
insert into users (id, created_at, updated_at, name)
values ($1, now(), now(), $2)
returning id, created_at, updated_at, name;

-- name: GetUserByName :one
select id, created_at, updated_at, name
from users
where name = $1;

-- noinspection SqlWithoutWhere
-- name: DeleteAllUsers :exec
delete from users;

-- name: GetUsers :many
select id, created_at, updated_at, name
from users;