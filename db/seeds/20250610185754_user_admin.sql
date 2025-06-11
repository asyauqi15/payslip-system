-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
INSERT INTO users (username, password_hash, role) VALUES ('admin', crypt('password', gen_salt('bf')), 'admin');
-- +goose StatementEnd