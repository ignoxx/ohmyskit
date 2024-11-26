-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = ?1;

-- name: FindUserByID :one
SELECT * FROM users WHERE id = ?1;

-- name: UpdateUserEmailVerifiedAt :exec
UPDATE users SET email_verified_at = CURRENT_TIMESTAMP WHERE id = ?1;

-- name: UpdateUserFirstLastName :exec
UPDATE users SET first_name = ?1, last_name = ?2 WHERE id = ?3;

-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, first_name, last_name, created_at, updated_at
) VALUES (
    ?1, ?2, ?3, ?4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING *;
