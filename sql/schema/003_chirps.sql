-- +goose Up
CREATE TABLE chirps (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    body text not null,
    user_id uuid not null REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;