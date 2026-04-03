-- name: CreateUser :one
INSERT INTO users (name, lastname, email, password, created_at, status, is_admin)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, lastname = $3, email = $4, status = $5, is_admin = $6
WHERE id = $1;

-- name: FindUserById :one
SELECT id, name, lastname, email, created_at, status, is_admin
FROM users
WHERE id = $1;

-- name: FindUserByEmail :one
SELECT id, name, lastname, email, created_at, status, is_admin
FROM users
WHERE email = $1;

-- name: GetPassword :one 
SELECT id, password
FROM users  
WHERE email = $1;
