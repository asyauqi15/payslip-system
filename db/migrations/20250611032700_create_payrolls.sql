-- +goose Up
-- +goose StatementBegin
CREATE TABLE payrolls (
    id BIGSERIAL PRIMARY KEY,
    attendance_period_id BIGINT NOT NULL UNIQUE,
    employees_count BIGINT NOT NULL,
    total_reimbursement BIGINT NOT NULL,
    total_overtime BIGINT NOT NULL,
    total_payroll BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (attendance_period_id) REFERENCES attendance_periods(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_payrolls_attendance_period_id ON payrolls(attendance_period_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payrolls;
-- +goose StatementEnd
