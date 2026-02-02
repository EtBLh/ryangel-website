CREATE TABLE sys_config (
    config_key VARCHAR(50) PRIMARY KEY,
    config_value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Seed initial values
INSERT INTO sys_config (config_key, config_value, description) VALUES 
('shipping_fee_ebuy', '5.0', 'Shipping fee for Macau Ebuy pickup'),
('shipping_fee_sf', '40.0', 'Shipping fee for SF Express (Direct Address)');
