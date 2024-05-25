CREATE SCHEMA shorts;

create table if not exists shorts.shorts
(
    id SERIAL PRIMARY KEY,
    series_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    duration INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)

SELECT pg_create_logical_replication_slot('shorts_replication_slot', 'pgoutput');
CREATE PUBLICATION shorts_publication FOR TABLE shorts.shorts;
