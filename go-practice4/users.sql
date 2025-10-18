CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    balance DECIMAL(10,2) DEFAULT 0.00
);

INSERT INTO users (name, email, balance) VALUES
('Alice Johnson', 'alice@example.com', 2000.00),
('Bob Smith', 'bob@example.com', 600.00),
('Charlie Brown', 'charlie@example.com', 850.00),
('Diana Prince', 'diana@example.com', 1100.00)
ON CONFLICT (email) DO NOTHING;