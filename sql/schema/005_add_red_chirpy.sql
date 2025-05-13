-- +goose Up
ALTER TABLE users ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER Thashed_passwordABLE users DROP COLUMN is_chirpy_red;
