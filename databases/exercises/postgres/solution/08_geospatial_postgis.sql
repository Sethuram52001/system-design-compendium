CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS places(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    location geography(Point, 4326) NOT NULL
);

CREATE INDEX IF NOT EXISTS places_location_gist_idx ON places USING GIST(location);

SELECT name
FROM places
WHERE ST_DWithin(
  location,
  ST_MakePoint(-73.9857, 40.7484)::geography,
  5000
);