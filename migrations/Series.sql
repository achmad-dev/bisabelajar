CREATE SCHEMA series;

create table if not exists series.series
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

SELECT pg_create_logical_replication_slot('series_replication_slot', 'pgoutput');
CREATE PUBLICATION series_publication FOR TABLE series.series;
