-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;
-- name: GetHashedPassword :one
SELECT * FROM users
WHERE email=$1;
-- name: UpdateUser :one
UPDATE users
SET email=$2, hashed_password=$3, updated_at=$4
WHERE id=$1
RETURNING *;
-- name: SetChirpyRed :exec
UPDATE users
SET is_chirpy_red=true
WHERE id=$1;
-- name: DeleteUsers :exec
DELETE FROM users;
