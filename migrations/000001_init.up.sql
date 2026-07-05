CREATE SCHEMA IF NOT EXISTS todoapp;

CREATE TABLE IF NOT EXISTS todoapp.users (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    full_name VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15) CHECK(
        phone_number ~ '^\+\d{9,14}$'
    )
);

CREATE TABLE IF NOT EXISTS todoapp.tasks (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL DEFAULT 1,
    title VARCHAR(100) NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    description VARCHAR(1000) CHECK (char_length(description) BETWEEN 1 AND 1000),
    completed BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,

    CHECK(
        (completed = FALSE AND completed_at IS NULL) OR
        (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
    ),

    author_user_id INTEGER NOT NULL REFERENCES todoapp.users(id)
);
