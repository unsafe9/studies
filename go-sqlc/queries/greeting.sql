-- name: AddGreeting :one
INSERT INTO greeting(content)
VALUES ($1)
RETURNING id;

-- name: GetGreeting :one
SELECT *
FROM greeting
WHERE id = $1
LIMIT 1;
