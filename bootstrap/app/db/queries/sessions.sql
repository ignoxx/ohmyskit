-- name: CreateSession :one
INSERT INTO sessions (
    user_id, token, expires_at, created_at, updated_at
) VALUES (
    ?1, ?2, ?3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: DeleteSessionByToken :exec
DELETE FROM sessions WHERE token = ?1;

-- name: FindSessionByTokenAndExpiration :one
-- and also return the associated user
SELECT sessions.*, users.id as user_id, users.email as user_email
FROM sessions
JOIN users ON sessions.user_id = users.id
WHERE token = ?1 AND expires_at > CURRENT_TIMESTAMP;
