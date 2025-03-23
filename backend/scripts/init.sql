
CREATE TABLE IF NOT EXISTS test_users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at DATE DEFAULT CURRENT_DATE,
    updated_at DATE DEFAULT CURRENT_DATE
);

CREATE TABLE IF NOT EXISTS test_products (
    id SERIAL PRIMARY KEY,
    product_reference VARCHAR(100) NOT NULL UNIQUE,
    product_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    product_category VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock_location VARCHAR(255) NOT NULL,
    supplier VARCHAR(255) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    created_at DATE DEFAULT CURRENT_DATE,
    updated_at DATE DEFAULT CURRENT_DATE
);

-- Insert initial data (optional)
INSERT INTO test_users (name, email, password) 
VALUES ('admin', 'admin@example.com', 'hashedpassword')
ON CONFLICT (email) DO NOTHING;

INSERT INTO test_products (product_reference, product_name, created_at, status, product_category, price, stock_location, supplier, quantity) 
VALUES 
    ('PROD-202503-001', 'Laptop', '2025-03-01', 'Available', 'Electronics', 1200.00, 'Warehouse A', 'TechCorp', 10),
    ('PROD-202503-002', 'Phone', '2025-03-02', 'Available', 'Electronics', 800.00, 'Warehouse B', 'GadgetCo', 15)
ON CONFLICT (product_reference) DO NOTHING;