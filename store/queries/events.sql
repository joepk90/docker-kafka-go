-- name: AddEventRecord :one
INSERT INTO kafka_events (event_id, event_tag, event_stat, event_msg, request_date)
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
RETURNING id;

-- name: GetEventRecord :one
SELECT id
FROM kafka_events
WHERE id = $1;
