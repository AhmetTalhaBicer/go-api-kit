-- sqlc query file for domain2.
-- After editing, run `make sqlc` to regenerate db/sqlc/.

-- name: ListDomain2 :many
SELECT id, name, created_at, updated_at
FROM domain2
ORDER BY id;

-- name: GetDomain2 :one
SELECT id, name, created_at, updated_at
FROM domain2
WHERE id = $1
LIMIT 1;

-- name: CreateDomain2 :one
INSERT INTO domain2 (name, created_at, updated_at)
VALUES ($1, NOW(), NOW())
RETURNING id, name, created_at, updated_at;

-- name: UpdateDomain2 :one
UPDATE domain2
SET name       = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, created_at, updated_at;

-- name: DeleteDomain2 :exec
DELETE FROM domain2
WHERE id = $1;
