-- +goose Up
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    company TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    image TEXT NOT NULL,
    description TEXT,
    tag TEXT,
    location TEXT NOT NULL,
    published_at TIMESTAMP
);

-- +goose Down
DROP TABLE jobs;