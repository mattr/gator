// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
insert into users (id, created_at, updated_at, name)
values ($1, now(), now(), $2)
returning id, created_at, updated_at, name
`

type CreateUserParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
delete from users
`

// noinspection SqlWithoutWhere
func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const getUserByName = `-- name: GetUserByName :one
select id, created_at, updated_at, name
from users
where name = $1
`

func (q *Queries) GetUserByName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}
