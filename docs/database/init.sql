-- 1. Таблица филиалов
CREATE TABLE branches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Таблица пользователей (Клиенты, Менеджеры, Суперадмины)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'client',
    telegram_id BIGINT UNIQUE,
    branch_id INTEGER REFERENCES branches(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- CHECK-проверка для ролей
    CONSTRAINT check_user_role CHECK (role IN ('client', 'manager', 'super_admin'))
);

-- 3. Таблица автомобилей
CREATE TABLE cars (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    license_plate VARCHAR(50) NOT NULL UNIQUE,
    brand VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Таблица заказов-нарядов
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    branch_id INTEGER NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    car_id INTEGER NOT NULL REFERENCES cars(id) ON DELETE RESTRICT,
    service_type VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'new',
    preferred_date DATE NOT NULL,
    preferred_time TIME NOT NULL,
    cost NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    client_confirmed BOOLEAN DEFAULT FALSE
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- CHECK-проверка для статусов заказа
    CONSTRAINT check_order_status CHECK (status IN ('new', 'confirmed', 'cancel_requested', 'canceled', 'in_progress', 'ready', 'completed'))
);
