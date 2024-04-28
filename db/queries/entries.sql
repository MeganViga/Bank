-- name: CreateEntry :one
INSERT INTO entries (
 account_id,
 amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntryByID :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- -- name: UpdateAccount :one
-- UPDATE account
--   set balance = $2
-- WHERE id = $1
-- RETURNING *;

-- name: DeleteEntry :one
DELETE FROM entries
WHERE id = $1
RETURNING *;