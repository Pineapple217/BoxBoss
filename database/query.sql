-- name: GetRepo :one
SELECT * FROM repositories
WHERE id = ? LIMIT 1;

-- name: ListRepos :many
SELECT * FROM repositories
ORDER BY name;

-- -- name: CreateAuthor :one
-- INSERT INTO authors (
--   name, bio
-- ) VALUES (
--   ?, ?
-- )
-- RETURNING *;

-- -- name: UpdateAuthor :exec
-- UPDATE authors
-- set name = ?,
-- bio = ?
-- WHERE id = ?;

-- name: DeleteRepo :exec
DELETE FROM repositories
WHERE id = ?;