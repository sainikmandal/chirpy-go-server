-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    email text not null unique
);

-- +goose Down
DROP TABLE users;