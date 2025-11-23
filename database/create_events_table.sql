-- SQL to create the events table
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    description TEXT NOT NULL
);