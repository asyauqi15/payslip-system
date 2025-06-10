-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
INSERT INTO users (email, password_hash, role) VALUES ('admin@example.com', crypt('plaintext_password', gen_salt('bf')), 'admin');
-- +goose StatementEnd