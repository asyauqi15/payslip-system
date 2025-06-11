-- +goose Up
-- +goose StatementBegin

-- Create attendance periods for 1-month ranges
-- Starting from January 2024 to December 2025 (24 months total)
INSERT INTO attendance_periods (start_date, end_date) VALUES

('2025-01-01', '2025-01-31'),
('2025-02-01', '2025-02-28'),
('2025-03-01', '2025-03-31'),
('2025-04-01', '2025-04-30'),
('2025-05-01', '2025-05-31'),
('2025-06-01', '2025-06-30');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM attendance_periods WHERE start_date >= '2025-01-01' AND end_date <= '2025-06-30';
-- +goose StatementEnd
