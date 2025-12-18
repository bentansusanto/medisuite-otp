-- name: CreateSession :one
INSERT INTO user_sessions (
    user_id,
    ref_token,
    is_blocked,
    client_ip,
    expires_at
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM user_sessions
WHERE ref_token = $1
  AND is_blocked = false
  AND expires_at > NOW()
LIMIT 1;

-- name: GetSessionsByUserId :one
SELECT * FROM user_sessions
WHERE user_id = $1::uuid
  AND is_blocked = false
  AND expires_at > NOW()
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM user_sessions WHERE ref_token = $1;

-- name: BlockSession :exec
UPDATE user_sessions SET is_blocked = true WHERE id = $1::uuid;

-- name: BlockAllUserSessions :exec
UPDATE user_sessions SET is_blocked = true WHERE user_id = $1::uuid;
