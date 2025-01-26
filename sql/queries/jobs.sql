-- name: CreateRemoteJob :one
INSERT INTO jobs (id, created_at, updated_at, title, company, url,image, description,tag,location, published_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;
--

-- name: GetRemoteJobs :many
SELECT * FROM jobs;
--