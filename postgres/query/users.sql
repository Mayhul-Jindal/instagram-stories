-- name: CreateUser :one
INSERT INTO users (
  email, hashed_password
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetUsersEmails :many
select email from users
where id != $1
limit $2
offset $3;

-- name: GetUserById :one
select * from users
where id = $1
limit 1;

-- name: GetUserByEmail :one
select * from users
where email = $1;

-- name: GetUserByIdAndEmail :one
select * from users
where id = $1 and email = $2
limit 1;