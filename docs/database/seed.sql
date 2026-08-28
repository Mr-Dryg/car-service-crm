INSERT INTO branches (name, address, phone) VALUES
('Северный филиал', 'ул. Полярная, д. 10', '+74951112233'),
('Южный филиал', 'ул. Солнечная, д. 20', '+74954445566');

INSERT INTO users (name, phone, email, role, branch_id) VALUES
('Иван Иванов', '+79991111111', 'ivan@example.com', 'client', 1),
('Пётр Петров', '+79992222222', 'petr@example.com', 'client', 1),
('Анна Сидорова', '+79993333333', 'anna@example.com', 'client', 2);

INSERT INTO cars (user_id, license_plate, brand, model) VALUES
((SELECT id FROM users WHERE email = 'ivan@example.com'), 'А111АА77', 'Toyota', 'Camry'),
((SELECT id FROM users WHERE email = 'petr@example.com'), 'В222ВВ77', 'BMW', 'X5'),
((SELECT id FROM users WHERE email = 'anna@example.com'), 'С333СС77', 'Kia', 'Rio');

INSERT INTO orders (branch_id, car_id, service_type, status, preferred_date, preferred_time, price) VALUES
(
    1, 
    (SELECT id FROM cars WHERE license_plate = 'А111АА77'), 
    'Замена масла', 
    'new', 
    '2026-09-01', 
    '10:00:00', 
    4500.00
);

INSERT INTO orders (branch_id, car_id, service_type, status, preferred_date, preferred_time, price) VALUES
(
    1, 
    (SELECT id FROM cars WHERE license_plate = 'В222ВВ77'), 
    'Диагностика подвески', 
    'confirmed', 
    '2026-09-01', 
    '12:00:00', 
    2000.00
);

INSERT INTO orders (branch_id, car_id, service_type, status, preferred_date, preferred_time, price) VALUES
(
    2, 
    (SELECT id FROM cars WHERE license_plate = 'С333СС77'), 
    'Шиномонтаж', 
    'new', 
    '2026-09-02', 
    '15:30:00', 
    3200.00
);
