-- +goose Up
-- +goose StatementBegin
CREATE TABLE payslips (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL,
    payroll_id BIGINT NOT NULL,
    base_salary BIGINT NOT NULL,
    attendance_count INT NOT NULL DEFAULT 0,
    total_working_days INT NOT NULL,
    prorated_salary BIGINT NOT NULL,
    overtime_total_hours INT NOT NULL DEFAULT 0,
    overtime_total_amount BIGINT NOT NULL DEFAULT 0,
    reimbursement_total BIGINT NOT NULL,
    total_take_home BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    FOREIGN KEY (payroll_id) REFERENCES payrolls(id) ON DELETE CASCADE
);

CREATE INDEX idx_payslips_employee_id ON payslips(employee_id);
CREATE INDEX idx_payslips_payroll_id ON payslips(payroll_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payslips;
-- +goose StatementEnd
