-- +goose Up
-- +goose StatementBegin
CREATE TABLE overtimes (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL,
    start_at TIMESTAMP WITH TIME ZONE NOT NULL,
    end_at TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX idx_overtimes_employee_id ON overtimes(employee_id);
CREATE INDEX idx_overtimes_start_at ON overtimes(start_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS overtimes;
-- +goose StatementEnd
