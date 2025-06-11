-- +goose Up
-- +goose StatementBegin

-- Create reimbursement records for employees from January to May 2025
-- This will generate realistic business expense reimbursements with varying amounts and purposes

-- January 2025 Reimbursements
-- Business travel and client meetings
INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    250000, -- 250k IDR for transportation
    '2025-01-08 10:30:00+07:00',
    'Transportation costs for client meeting at downtown office'
FROM employees e
WHERE e.id BETWEEN 10 AND 25;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    500000, -- 500k IDR for hotel
    '2025-01-15 14:20:00+07:00',
    'Hotel accommodation for 2-day business trip to Surabaya'
FROM employees e
WHERE e.id BETWEEN 30 AND 40;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    150000, -- 150k IDR for meals
    '2025-01-22 12:15:00+07:00',
    'Business lunch with potential clients'
FROM employees e
WHERE e.id BETWEEN 50 AND 70;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    300000, -- 300k IDR for internet
    '2025-01-29 16:45:00+07:00',
    'Internet allowance for work from home setup'
FROM employees e
WHERE e.id BETWEEN 1 AND 20;

-- February 2025 Reimbursements
-- Training and professional development
INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    750000, -- 750k IDR for training
    '2025-02-05 09:00:00+07:00',
    'Professional certification training course fee'
FROM employees e
WHERE e.id BETWEEN 15 AND 35;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    200000, -- 200k IDR for transportation
    '2025-02-12 11:30:00+07:00',
    'Transportation for vendor meeting in Bandung'
FROM employees e
WHERE e.id BETWEEN 45 AND 60;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    400000, -- 400k IDR for equipment
    '2025-02-18 13:20:00+07:00',
    'Office supplies and equipment for remote work'
FROM employees e
WHERE e.id BETWEEN 25 AND 45;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    180000, -- 180k IDR for meals
    '2025-02-25 19:15:00+07:00',
    'Team dinner after successful project completion'
FROM employees e
WHERE e.id BETWEEN 65 AND 85;

-- March 2025 Reimbursements
-- Quarter-end activities and conferences
INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    1200000, -- 1.2M IDR for conference
    '2025-03-07 08:45:00+07:00',
    'Technology conference registration and materials'
FROM employees e
WHERE e.id BETWEEN 5 AND 15;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    350000, -- 350k IDR for transportation
    '2025-03-14 15:30:00+07:00',
    'Flight tickets for business meeting in Medan'
FROM employees e
WHERE e.id BETWEEN 35 AND 50;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    220000, -- 220k IDR for communication
    '2025-03-21 10:10:00+07:00',
    'Mobile phone bill for business communication'
FROM employees e
WHERE e.id BETWEEN 55 AND 75;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    450000, -- 450k IDR for accommodation
    '2025-03-28 17:00:00+07:00',
    'Hotel and meals for 3-day client workshop'
FROM employees e
WHERE e.id BETWEEN 20 AND 35;

-- April 2025 Reimbursements
-- Spring activities and team building
INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    300000, -- 300k IDR for team building
    '2025-04-04 12:30:00+07:00',
    'Team building activity expenses and transportation'
FROM employees e
WHERE e.id BETWEEN 40 AND 60;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    250000, -- 250k IDR for parking
    '2025-04-11 14:45:00+07:00',
    'Parking fees for multiple client visits this month'
FROM employees e
WHERE e.id BETWEEN 70 AND 90;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    600000, -- 600k IDR for medical
    '2025-04-18 16:20:00+07:00',
    'Medical check-up and health screening reimbursement'
FROM employees e
WHERE e.id BETWEEN 10 AND 30;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    380000, -- 380k IDR for software
    '2025-04-25 11:00:00+07:00',
    'Professional software license and tools subscription'
FROM employees e
WHERE e.id BETWEEN 50 AND 70;

-- May 2025 Reimbursements
-- Mid-year preparation and maintenance
INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    420000, -- 420k IDR for maintenance
    '2025-05-02 09:30:00+07:00',
    'Vehicle maintenance for business travel'
FROM employees e
WHERE e.id BETWEEN 25 AND 45;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    280000, -- 280k IDR for stationery
    '2025-05-09 13:45:00+07:00',
    'Office stationery and printing costs'
FROM employees e
WHERE e.id BETWEEN 60 AND 80;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    500000, -- 500k IDR for books
    '2025-05-16 15:15:00+07:00',
    'Professional development books and learning materials'
FROM employees e
WHERE e.id BETWEEN 15 AND 35;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    320000, -- 320k IDR for utilities
    '2025-05-23 10:50:00+07:00',
    'Electricity bill increase due to work from home'
FROM employees e
WHERE e.id BETWEEN 85 AND 100;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    650000, -- 650k IDR for emergency
    '2025-05-30 18:30:00+07:00',
    'Emergency travel expenses for urgent client support'
FROM employees e
WHERE e.id BETWEEN 5 AND 20;

-- Additional monthly allowances and recurring expenses
INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    200000, -- 200k IDR monthly allowance
    '2025-01-31 23:59:00+07:00',
    'Monthly transportation allowance'
FROM employees e
WHERE e.id BETWEEN 1 AND 50;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    200000, -- 200k IDR monthly allowance
    '2025-02-28 23:59:00+07:00',
    'Monthly transportation allowance'
FROM employees e
WHERE e.id BETWEEN 51 AND 100;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    200000, -- 200k IDR monthly allowance
    '2025-03-31 23:59:00+07:00',
    'Monthly transportation allowance'
FROM employees e
WHERE e.id BETWEEN 1 AND 50;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    200000, -- 200k IDR monthly allowance
    '2025-04-30 23:59:00+07:00',
    'Monthly transportation allowance'
FROM employees e
WHERE e.id BETWEEN 51 AND 100;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    200000, -- 200k IDR monthly allowance
    '2025-05-31 23:59:00+07:00',
    'Monthly transportation allowance'
FROM employees e
WHERE e.id BETWEEN 1 AND 50;

-- Special project reimbursements
INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    800000, -- 800k IDR for special project
    '2025-02-14 16:00:00+07:00',
    'Valentine campaign project additional expenses'
FROM employees e
WHERE e.id BETWEEN 20 AND 30;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    1000000, -- 1M IDR for major project
    '2025-03-31 20:00:00+07:00',
    'Q1 closing project completion bonus reimbursement'
FROM employees e
WHERE e.id BETWEEN 5 AND 15;

INSERT INTO reimbursements (employee_id, amount, date, description) 
SELECT 
    e.id,
    450000, -- 450k IDR for weekend work
    '2025-04-12 14:00:00+07:00',
    'Weekend work meal and transportation compensation'
FROM employees e
WHERE e.id BETWEEN 25 AND 40;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM reimbursements WHERE date >= '2025-01-01' AND date < '2025-06-01';
-- +goose StatementEnd
