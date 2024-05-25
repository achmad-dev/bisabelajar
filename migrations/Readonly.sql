CREATE SCHEMA readonlydb;

create table if not exists readonlydb.series
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

create table if not exists readonlydb.shorts
(
    id SERIAL PRIMARY KEY,
    series_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    duration INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES readonlydb.series(id)
);

-- For shorts database
CREATE SUBSCRIPTION shorts_subscription
  CONNECTION 'dbname=shorts host=<shorts_host> user=<username> password=<password>'
  PUBLICATION shorts_publication;

-- For series database
CREATE SUBSCRIPTION series_subscription
  CONNECTION 'dbname=series host=<series_host> user=<username> password=<password>'
  PUBLICATION series_publication;