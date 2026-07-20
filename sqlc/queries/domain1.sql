-- sqlc query file for domain1.
-- After editing, run `make sqlc` to regenerate db/sqlc/.

-- name: ListDomain1 :many
SELECT id, name, created_at, updated_at
FROM domain1
ORDER BY id;

-- name: GetDomain1 :one
SELECT id, name, created_at, updated_at
FROM domain1
WHERE id = $1
LIMIT 1;

-- name: CreateDomain1 :one
INSERT INTO domain1 (name, created_at, updated_at)
VALUES ($1, NOW(), NOW())
RETURNING id, name, created_at, updated_at;

-- name: UpdateDomain1 :one
UPDATE domain1
SET name       = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, created_at, updated_at;

-- name: DeleteDomain1 :exec
DELETE FROM domain1
WHERE id = $1;
