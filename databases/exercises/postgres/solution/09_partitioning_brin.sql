CREATE TABLE IF NOT EXISTS events(
    id BIGSERIAL,
    created_at TIMESTAMPTZ NOT NULL,
    payload JSONB NOT NULL
) PARTITION BY RANGE (created_at);

-- child partition tables
CREATE TABLE IF NOT EXISTS events_2026_05
PARTITION OF events 
FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');


CREATE INDEX IF NOT EXISTS events_2026_05_created_at_brin_idx
ON events_2026_05
USING BRIN(created_at);

CREATE INDEX IF NOT EXISTS events_2026_06_created_at_brin_idx
ON events_2026_06
USING BRIN(created_at);