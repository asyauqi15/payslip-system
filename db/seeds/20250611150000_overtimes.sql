-- +goose Up
-- +goose StatementBegin

-- Create overtime records for employees from January to May 2025
-- This will generate realistic overtime patterns with varying frequencies and durations

-- January 2025 Overtime Records
-- Project deadline scenarios and end-of-month activities
INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-01-15 18:00:00+07:00',
    '2025-01-15 20:30:00+07:00',
    'Year-end system maintenance and data migration'
FROM employees e
WHERE e.id <= 25; -- Senior employees handling critical tasks

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-01-22 17:30:00+07:00',
    '2025-01-22 20:30:00+07:00',
    'Monthly report preparation and client presentation'
FROM employees e
WHERE e.id BETWEEN 10 AND 40;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-01-29 18:00:00+07:00',
    '2025-01-29 19:45:00+07:00',
    'End of month inventory and financial closing'
FROM employees e
WHERE e.id BETWEEN 30 AND 60;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-01-31 17:30:00+07:00',
    '2025-01-31 20:30:00+07:00',
    'Critical system deployment and monitoring'
FROM employees e
WHERE e.id BETWEEN 5 AND 20;

-- February 2025 Overtime Records
-- Valentine project and quarterly preparations
INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-02-10 18:30:00+07:00',
    '2025-02-10 21:30:00+07:00',
    'Special Valentine campaign system updates'
FROM employees e
WHERE e.id BETWEEN 15 AND 45;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-02-18 17:45:00+07:00',
    '2025-02-18 20:00:00+07:00',
    'Database optimization and performance tuning'
FROM employees e
WHERE e.id BETWEEN 25 AND 55;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-02-25 18:00:00+07:00',
    '2025-02-25 19:30:00+07:00',
    'Monthly security audit and compliance check'
FROM employees e
WHERE e.id BETWEEN 35 AND 65;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-02-28 17:30:00+07:00',
    '2025-02-28 20:30:00+07:00',
    'Quarter-end financial reconciliation and reporting'
FROM employees e
WHERE e.id BETWEEN 1 AND 30;

-- March 2025 Overtime Records
-- Quarter-end activities and spring projects
INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-03-05 18:15:00+07:00',
    '2025-03-05 20:45:00+07:00',
    'New feature development and testing'
FROM employees e
WHERE e.id BETWEEN 20 AND 50;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-03-12 17:30:00+07:00',
    '2025-03-12 19:15:00+07:00',
    'Client demo preparation and system setup'
FROM employees e
WHERE e.id BETWEEN 40 AND 70;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-03-20 18:00:00+07:00',
    '2025-03-20 21:00:00+07:00',
    'Major system upgrade and migration'
FROM employees e
WHERE e.id BETWEEN 10 AND 35;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-03-27 17:45:00+07:00',
    '2025-03-27 20:15:00+07:00',
    'End of quarter audit and documentation'
FROM employees e
WHERE e.id BETWEEN 50 AND 80;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-03-31 18:30:00+07:00',
    '2025-03-31 21:30:00+07:00',
    'Q1 closing activities and data backup'
FROM employees e
WHERE e.id BETWEEN 5 AND 25;

-- April 2025 Overtime Records
-- Spring season projects and mid-year preparations
INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-04-08 18:00:00+07:00',
    '2025-04-08 20:30:00+07:00',
    'Spring product launch preparation'
FROM employees e
WHERE e.id BETWEEN 30 AND 60;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-04-15 17:30:00+07:00',
    '2025-04-15 19:45:00+07:00',
    'Mid-month performance optimization'
FROM employees e
WHERE e.id BETWEEN 15 AND 45;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-04-22 18:15:00+07:00',
    '2025-04-22 21:15:00+07:00',
    'Customer support system enhancement'
FROM employees e
WHERE e.id BETWEEN 35 AND 65;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-04-29 17:45:00+07:00',
    '2025-04-29 20:00:00+07:00',
    'Monthly infrastructure maintenance'
FROM employees e
WHERE e.id BETWEEN 25 AND 55;

-- May 2025 Overtime Records
-- Labor day preparation and mid-year planning
INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-05-06 18:30:00+07:00',
    '2025-05-06 21:30:00+07:00',
    'Post-holiday system recovery and updates'
FROM employees e
WHERE e.id BETWEEN 20 AND 50;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-05-13 17:30:00+07:00',
    '2025-05-13 19:30:00+07:00',
    'Mid-year planning and strategy session'
FROM employees e
WHERE e.id BETWEEN 40 AND 70;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-05-20 18:00:00+07:00',
    '2025-05-20 20:45:00+07:00',
    'Quality assurance and testing cycle'
FROM employees e
WHERE e.id BETWEEN 10 AND 40;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-05-27 17:45:00+07:00',
    '2025-05-27 20:45:00+07:00',
    'End of month processing and reconciliation'
FROM employees e
WHERE e.id BETWEEN 30 AND 60;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-05-30 18:15:00+07:00',
    '2025-05-30 19:45:00+07:00',
    'Memorial day system maintenance'
FROM employees e
WHERE e.id BETWEEN 50 AND 80;

-- Additional weekend overtime for critical projects
INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-01-11 09:00:00+07:00',
    '2025-01-11 13:00:00+07:00',
    'Weekend emergency maintenance'
FROM employees e
WHERE e.id BETWEEN 1 AND 15;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-02-15 10:00:00+07:00',
    '2025-02-15 14:30:00+07:00',
    'Weekend project deployment'
FROM employees e
WHERE e.id BETWEEN 20 AND 35;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-03-22 08:30:00+07:00',
    '2025-03-22 12:00:00+07:00',
    'Weekend data migration and backup'
FROM employees e
WHERE e.id BETWEEN 45 AND 60;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-04-12 09:15:00+07:00',
    '2025-04-12 13:45:00+07:00',
    'Weekend security patch installation'
FROM employees e
WHERE e.id BETWEEN 25 AND 40;

INSERT INTO overtimes (employee_id, start_at, end_at, description) 
SELECT 
    e.id,
    '2025-05-17 10:30:00+07:00',
    '2025-05-17 15:00:00+07:00',
    'Weekend performance testing'
FROM employees e
WHERE e.id BETWEEN 35 AND 50;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM overtimes WHERE start_at >= '2025-01-01' AND start_at < '2025-06-01';
-- +goose StatementEnd
