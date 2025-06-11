-- +goose Up
-- +goose StatementBegin
CREATE TABLE reimbursements (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL,
    amount BIGINT NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

CREATE INDEX idx_reimbursements_employee_id ON reimbursements(employee_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reimbursements;
-- +goose StatementEnd
