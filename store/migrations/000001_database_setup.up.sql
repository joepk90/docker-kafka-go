CREATE TABLE IF NOT EXISTS kafka_events (
    id BIGINT PRIMARY KEY,
    tag VARCHAR(255),
    stat VARCHAR(255),
    msg VARCHAR(255),
    -- created_at TIMESTAMPTZ NOT NULL
    request_date TIMESTAMPTZ NOT NULL
);