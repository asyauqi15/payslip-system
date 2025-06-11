-- +goose Up
-- +goose StatementBegin

-- Create attendance records for all employees from January to May 2025
-- This will generate realistic attendance patterns for 100 employees over 5 months

-- January 2025 (31 days, excluding weekends: ~22 working days)
-- Working days: Jan 2,3,6,7,8,9,10,13,14,15,16,17,20,21,22,23,24,27,28,29,30,31
INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-01-02 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-01-02 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 95; -- 95% attendance rate

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-01-03 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-01-03 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 98;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-01-06 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-01-06 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 93;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-01-07 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-01-07 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 97;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-01-08 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-01-08 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 96;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-01-09 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-01-09 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 94;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-01-10 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-01-10 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 99;

-- February 2025 (28 days, ~20 working days)
INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-02-03 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-02-03 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 92;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-02-04 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-02-04 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 95;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-02-05 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-02-05 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 98;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-02-06 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-02-06 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 93;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-02-07 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-02-07 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 97;

-- March 2025 (31 days, ~21 working days)
INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-03-03 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-03-03 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 96;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-03-04 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-03-04 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 94;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-03-05 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-03-05 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 99;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-03-06 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-03-06 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 91;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-03-07 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-03-07 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 95;

-- April 2025 (30 days, ~22 working days)
INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-04-01 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-04-01 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 98;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-04-02 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-04-02 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 93;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-04-03 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-04-03 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 97;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-04-04 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-04-04 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 96;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-04-07 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-04-07 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 94;

-- May 2025 (31 days, ~22 working days)
INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-05-02 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-05-02 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 99;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-05-05 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-05-05 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 91;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-05-06 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-05-06 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 95;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-05-07 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-05-07 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 98;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-05-08 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-05-08 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 93;

INSERT INTO attendances (employee_id, clock_in_time, clock_out_time) 
SELECT 
    e.id,
    '2025-05-09 08:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00',
    '2025-05-09 17:' || LPAD((RANDOM() * 60)::INT::TEXT, 2, '0') || ':00'
FROM employees e
WHERE e.id <= 97;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM attendances WHERE clock_in_time >= '2025-01-01' AND clock_in_time < '2025-06-01';
-- +goose StatementEnd
