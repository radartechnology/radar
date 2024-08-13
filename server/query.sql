-- name: GetToken :one
SELECT *
FROM
    tokens
WHERE
    token = $1
LIMIT 1;
