-- name: InsertWord :exec
INSERT INTO words (word) VALUES ($1);


-- name: GetByPrefix :one
SELECT word, COUNT(word) FROM words WHERE word LIKE $1 GROUP BY word ORDER BY COUNT(word) DESC LIMIT 1;

-- name: List :many
SELECT word, COUNT(word) FROM words GROUP BY word;
