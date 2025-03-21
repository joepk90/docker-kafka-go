-- name: AddEventRecord :exec
INSERT INTO kafka_events (id, tag, stat, msg)
VALUES ($1, $2, $3, $4);

-- name: GetEventRecord :one
SELECT id
FROM kafka_events
WHERE id = $1;
