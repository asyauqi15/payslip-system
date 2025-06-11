-- +goose Up
-- +goose StatementBegin
CREATE TABLE attendance_periods (
    id BIGSERIAL PRIMARY KEY,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attendance_periods_start_date ON attendance_periods(start_date);
CREATE INDEX idx_attendance_periods_end_date ON attendance_periods(end_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS attendance_periods;
-- +goose StatementEnd
