CREATE TABLE IF NOT EXISTS kafka_events (
    id BIGSERIAL PRIMARY KEY,
    event_id BiGINT,
    event_tag VARCHAR(255),
    event_stat VARCHAR(255),
    event_msg VARCHAR(255),
    -- created_at TIMESTAMPTZ NOT NULL
    request_date TIMESTAMPTZ NOT NULL
);