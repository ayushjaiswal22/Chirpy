-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    $2,
    $2,
    $3,
    $4
)
RETURNING *;
-- name: GetUserByToken :one
SELECT user_id FROM refresh_tokens
WHERE token=$1 AND expires_at>$2 AND revoked_at IS NULL;
-- name: RevokeToken :exec
UPDATE refresh_tokens 
SET revoked_at=$2, updated_at=$2
WHERE token=$1;
-- name: DeletToken :exec
DELETE FROM refresh_tokens 
WHERE token=$1;
