-- +goose Up
-- +goose StatementBegin

INSERT INTO users (username, password_hash, role, status) VALUES
('employee001', crypt('password123', gen_salt('bf')), 'default', 1),
('employee002', crypt('password123', gen_salt('bf')), 'default', 1),
('employee003', crypt('password123', gen_salt('bf')), 'default', 1),
('employee004', crypt('password123', gen_salt('bf')), 'default', 1),
('employee005', crypt('password123', gen_salt('bf')), 'default', 1),
('employee006', crypt('password123', gen_salt('bf')), 'default', 1),
('employee007', crypt('password123', gen_salt('bf')), 'default', 1),
('employee008', crypt('password123', gen_salt('bf')), 'default', 1),
('employee009', crypt('password123', gen_salt('bf')), 'default', 1),
('employee010', crypt('password123', gen_salt('bf')), 'default', 1),
('employee011', crypt('password123', gen_salt('bf')), 'default', 1),
('employee012', crypt('password123', gen_salt('bf')), 'default', 1),
('employee013', crypt('password123', gen_salt('bf')), 'default', 1),
('employee014', crypt('password123', gen_salt('bf')), 'default', 1),
('employee015', crypt('password123', gen_salt('bf')), 'default', 1),
('employee016', crypt('password123', gen_salt('bf')), 'default', 1),
('employee017', crypt('password123', gen_salt('bf')), 'default', 1),
('employee018', crypt('password123', gen_salt('bf')), 'default', 1),
('employee019', crypt('password123', gen_salt('bf')), 'default', 1),
('employee020', crypt('password123', gen_salt('bf')), 'default', 1),
('employee021', crypt('password123', gen_salt('bf')), 'default', 1),
('employee022', crypt('password123', gen_salt('bf')), 'default', 1),
('employee023', crypt('password123', gen_salt('bf')), 'default', 1),
('employee024', crypt('password123', gen_salt('bf')), 'default', 1),
('employee025', crypt('password123', gen_salt('bf')), 'default', 1),
('employee026', crypt('password123', gen_salt('bf')), 'default', 1),
('employee027', crypt('password123', gen_salt('bf')), 'default', 1),
('employee028', crypt('password123', gen_salt('bf')), 'default', 1),
('employee029', crypt('password123', gen_salt('bf')), 'default', 1),
('employee030', crypt('password123', gen_salt('bf')), 'default', 1),
('employee031', crypt('password123', gen_salt('bf')), 'default', 1),
('employee032', crypt('password123', gen_salt('bf')), 'default', 1),
('employee033', crypt('password123', gen_salt('bf')), 'default', 1),
('employee034', crypt('password123', gen_salt('bf')), 'default', 1),
('employee035', crypt('password123', gen_salt('bf')), 'default', 1),
('employee036', crypt('password123', gen_salt('bf')), 'default', 1),
('employee037', crypt('password123', gen_salt('bf')), 'default', 1),
('employee038', crypt('password123', gen_salt('bf')), 'default', 1),
('employee039', crypt('password123', gen_salt('bf')), 'default', 1),
('employee040', crypt('password123', gen_salt('bf')), 'default', 1),
('employee041', crypt('password123', gen_salt('bf')), 'default', 1),
('employee042', crypt('password123', gen_salt('bf')), 'default', 1),
('employee043', crypt('password123', gen_salt('bf')), 'default', 1),
('employee044', crypt('password123', gen_salt('bf')), 'default', 1),
('employee045', crypt('password123', gen_salt('bf')), 'default', 1),
('employee046', crypt('password123', gen_salt('bf')), 'default', 1),
('employee047', crypt('password123', gen_salt('bf')), 'default', 1),
('employee048', crypt('password123', gen_salt('bf')), 'default', 1),
('employee049', crypt('password123', gen_salt('bf')), 'default', 1),
('employee050', crypt('password123', gen_salt('bf')), 'default', 1),
('employee051', crypt('password123', gen_salt('bf')), 'default', 1),
('employee052', crypt('password123', gen_salt('bf')), 'default', 1),
('employee053', crypt('password123', gen_salt('bf')), 'default', 1),
('employee054', crypt('password123', gen_salt('bf')), 'default', 1),
('employee055', crypt('password123', gen_salt('bf')), 'default', 1),
('employee056', crypt('password123', gen_salt('bf')), 'default', 1),
('employee057', crypt('password123', gen_salt('bf')), 'default', 1),
('employee058', crypt('password123', gen_salt('bf')), 'default', 1),
('employee059', crypt('password123', gen_salt('bf')), 'default', 1),
('employee060', crypt('password123', gen_salt('bf')), 'default', 1),
('employee061', crypt('password123', gen_salt('bf')), 'default', 1),
('employee062', crypt('password123', gen_salt('bf')), 'default', 1),
('employee063', crypt('password123', gen_salt('bf')), 'default', 1),
('employee064', crypt('password123', gen_salt('bf')), 'default', 1),
('employee065', crypt('password123', gen_salt('bf')), 'default', 1),
('employee066', crypt('password123', gen_salt('bf')), 'default', 1),
('employee067', crypt('password123', gen_salt('bf')), 'default', 1),
('employee068', crypt('password123', gen_salt('bf')), 'default', 1),
('employee069', crypt('password123', gen_salt('bf')), 'default', 1),
('employee070', crypt('password123', gen_salt('bf')), 'default', 1),
('employee071', crypt('password123', gen_salt('bf')), 'default', 1),
('employee072', crypt('password123', gen_salt('bf')), 'default', 1),
('employee073', crypt('password123', gen_salt('bf')), 'default', 1),
('employee074', crypt('password123', gen_salt('bf')), 'default', 1),
('employee075', crypt('password123', gen_salt('bf')), 'default', 1),
('employee076', crypt('password123', gen_salt('bf')), 'default', 1),
('employee077', crypt('password123', gen_salt('bf')), 'default', 1),
('employee078', crypt('password123', gen_salt('bf')), 'default', 1),
('employee079', crypt('password123', gen_salt('bf')), 'default', 1),
('employee080', crypt('password123', gen_salt('bf')), 'default', 1),
('employee081', crypt('password123', gen_salt('bf')), 'default', 1),
('employee082', crypt('password123', gen_salt('bf')), 'default', 1),
('employee083', crypt('password123', gen_salt('bf')), 'default', 1),
('employee084', crypt('password123', gen_salt('bf')), 'default', 1),
('employee085', crypt('password123', gen_salt('bf')), 'default', 1),
('employee086', crypt('password123', gen_salt('bf')), 'default', 1),
('employee087', crypt('password123', gen_salt('bf')), 'default', 1),
('employee088', crypt('password123', gen_salt('bf')), 'default', 1),
('employee089', crypt('password123', gen_salt('bf')), 'default', 1),
('employee090', crypt('password123', gen_salt('bf')), 'default', 1),
('employee091', crypt('password123', gen_salt('bf')), 'default', 1),
('employee092', crypt('password123', gen_salt('bf')), 'default', 1),
('employee093', crypt('password123', gen_salt('bf')), 'default', 1),
('employee094', crypt('password123', gen_salt('bf')), 'default', 1),
('employee095', crypt('password123', gen_salt('bf')), 'default', 1),
('employee096', crypt('password123', gen_salt('bf')), 'default', 1),
('employee097', crypt('password123', gen_salt('bf')), 'default', 1),
('employee098', crypt('password123', gen_salt('bf')), 'default', 1),
('employee099', crypt('password123', gen_salt('bf')), 'default', 1),
('employee100', crypt('password123', gen_salt('bf')), 'default', 1);

-- Create employee records with varying salaries (5M to 15M IDR)
INSERT INTO employees (user_id, base_salary)
SELECT u.id, 
       CASE 
           WHEN u.username LIKE 'employee0%' AND CAST(SUBSTRING(u.username, 9) AS INTEGER) <= 20 THEN 5000000  
           WHEN u.username LIKE 'employee0%' AND CAST(SUBSTRING(u.username, 9) AS INTEGER) <= 60 THEN 8000000  
           WHEN u.username LIKE 'employee0%' AND CAST(SUBSTRING(u.username, 9) AS INTEGER) <= 90 THEN 12000000 
           ELSE 15000000
       END
FROM users u 
WHERE u.username LIKE 'employee%' 
ORDER BY u.username;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM employees WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'employee%');
DELETE FROM users WHERE username LIKE 'employee%';
-- +goose StatementEnd
