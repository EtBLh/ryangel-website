--------------------------------------------
-- schema for RyAngel e-commerce platform --
-- in PostgreSQL format                   --
--------------------------------------------

-- Admin table
CREATE TABLE admin (
    admin_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    token VARCHAR(255),
    token_expiry TIMESTAMP
);

-- Client/Customer table
CREATE TABLE client (
    client_id SERIAL PRIMARY KEY,
    email VARCHAR(100) UNIQUE,
    username VARCHAR(50) UNIQUE,
    phone VARCHAR(20) NOT NULL UNIQUE,
    date_of_birth DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    token VARCHAR(255),
    token_expiry TIMESTAMP,
    otp_code VARCHAR(6),
    otp_code_expiry TIMESTAMP
);

-- Client addresses (ADD THIS MISSING TABLE)
CREATE TABLE client_address (
    address_id SERIAL PRIMARY KEY,
    client_id INT NOT NULL REFERENCES client(client_id) ON DELETE CASCADE,
    address_type VARCHAR(20) DEFAULT 'home',
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_client_address_default ON client_address (client_id)
WHERE is_default;

-- Ebuy stores (local pickup locations)
DROP TABLE IF EXISTS ebuy_store CASCADE;
CREATE TABLE ebuy_store (
    store_id VARCHAR(255) PRIMARY KEY,
    store_name VARCHAR(255) NOT NULL,
    type VARCHAR(50),
    office_hours VARCHAR(255),
    address TEXT,
    address_en TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create custom ENUM types first
CREATE TYPE product_type_enum AS ENUM ('faiachun', 'bag');
CREATE TYPE size_type_enum AS ENUM ('v-rect', 'square', 'fat-v-rect', 'big-square');
CREATE TYPE discount_type_enum AS ENUM ('percentage', 'fixed_amount', 'free_shipping', 'bxgy');
CREATE TYPE product_restriction_enum AS ENUM ('faiachun', 'bag', 'all');
CREATE TYPE applies_to_enum AS ENUM ('all_products', 'specific_products', 'specific_categories', 'first_order', 'bxgy_products');
CREATE TYPE customer_eligibility_enum AS ENUM ('all_customers', 'new_customers', 'specific_customers', 'vip_customers');
CREATE TYPE order_status_enum AS ENUM ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded');
CREATE TYPE payment_method_enum AS ENUM ('mpay', 'boc', 'bank_transfer');
CREATE TYPE payment_status_enum AS ENUM ('pending', 'paid', 'failed', 'refunded');
CREATE TYPE payment_proof_status_enum AS ENUM ('submitted', 'approved', 'rejected');

-- Products table
CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_type product_type_enum NOT NULL,
    hashtag VARCHAR(100),
    sku VARCHAR(100) UNIQUE NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    compare_at_price DECIMAL(10,2),
    cost_price DECIMAL(10,2),
    quantity INT DEFAULT 0,
    weight DECIMAL(8,2),
    dimensions JSONB,
    available_sizes size_type_enum[] DEFAULT '{}',
    is_featured BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    seo_title VARCHAR(255),
    seo_description TEXT,
    tags JSONB,
    created_by INT REFERENCES admin(admin_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Product categories
CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    category_name VARCHAR(100) NOT NULL,
    category_description TEXT,
    parent_category_id INT REFERENCES categories(category_id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Product-category mapping
CREATE TABLE product_categories (
    product_id INT REFERENCES products(product_id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(category_id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

-- Product images (paths served by API server)
CREATE TABLE product_images (
    image_id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    image_path VARCHAR(255) NOT NULL,
    alt_text VARCHAR(255),
    size_type size_type_enum,
    sort_order INT DEFAULT 0,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_product_images_primary ON product_images (product_id, size_type)
WHERE is_primary;

-- Discount table
CREATE TABLE discounts (
    discount_id SERIAL PRIMARY KEY,
    discount_code VARCHAR(50) UNIQUE,
    discount_name VARCHAR(100) NOT NULL,
    discount_type discount_type_enum NOT NULL,
    discount_value DECIMAL(10,2) DEFAULT NULL,
    
    -- BXGY specific fields
    buy_quantity INT DEFAULT NULL,
    get_quantity INT DEFAULT NULL,
    free_product_id INT REFERENCES products(product_id),
    applies_to_same_product BOOLEAN DEFAULT TRUE,
    
    -- Restrictions
    minimum_order_amount DECIMAL(10,2) DEFAULT 0.00,
    maximum_discount_amount DECIMAL(10,2) DEFAULT NULL,
    product_type_restriction product_restriction_enum DEFAULT 'all',
    
    -- Timing
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    
    -- Usage limits
    usage_limit INT DEFAULT NULL,
    used_count INT DEFAULT 0,
    usage_per_customer INT DEFAULT 1,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Application scope
    applies_to applies_to_enum NOT NULL,
    customer_eligibility customer_eligibility_enum DEFAULT 'all_customers',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Discount products mapping
CREATE TABLE discount_products (
    id SERIAL PRIMARY KEY,
    discount_id INT NOT NULL REFERENCES discounts(discount_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE
);

-- Discount categories mapping
CREATE TABLE discount_categories (
    id SERIAL PRIMARY KEY,
    discount_id INT NOT NULL REFERENCES discounts(discount_id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE
);

-- Orders table
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    client_id INT NOT NULL REFERENCES client(client_id),
    order_status order_status_enum DEFAULT 'pending',
    
    -- Pricing
    subtotal_amount DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    shipping_amount DECIMAL(10,2) DEFAULT 0.00,
    tax_amount DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL,
    
    -- Discount applied
    discount_id INT REFERENCES discounts(discount_id),
    discount_code VARCHAR(50),
    
    -- Shipping address
    -- Either a client address (for home delivery) or an ebuy store (local pickup)
    shipping_address_id INT REFERENCES client_address(address_id),
    ebuy_store_id VARCHAR(255) REFERENCES ebuy_store(store_id),
    
    -- Payment info
    payment_method payment_method_enum NOT NULL,
    payment_status payment_status_enum DEFAULT 'pending',
    payment_reference VARCHAR(100),
    
    -- Tracking
    tracking_number VARCHAR(100),
    shipping_carrier VARCHAR(50),
    
    -- Timestamps
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP NULL,
    shipped_at TIMESTAMP NULL,
    delivered_at TIMESTAMP NULL,
    cancelled_at TIMESTAMP NULL,
    
    -- Notes
    customer_notes TEXT,
    admin_notes TEXT
    ,
    CONSTRAINT orders_shipping_address_check CHECK (
        (
            shipping_address_id IS NOT NULL AND ebuy_store_id IS NULL
        ) OR (
            shipping_address_id IS NULL AND ebuy_store_id IS NOT NULL
        )
    )
);

-- Order items table
CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(product_id),
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    total_price DECIMAL(10,2) NOT NULL,
    
    -- Product snapshot
    product_name VARCHAR(255) NOT NULL,
    product_type product_type_enum NOT NULL,
    product_sku VARCHAR(100) NOT NULL,
    size_type size_type_enum,
    
    -- BXGY tracking
    is_free_item BOOLEAN DEFAULT FALSE,
    parent_discount_id INT REFERENCES discounts(discount_id)
);

-- Shopping cart table
CREATE TABLE cart (
    cart_id SERIAL PRIMARY KEY,
    client_id INT REFERENCES client(client_id) ON DELETE CASCADE,
    discount_id INT REFERENCES discounts(discount_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_items (
    cart_item_id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES cart(cart_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(product_id),
    size_type size_type_enum,
    quantity INT NOT NULL,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (cart_id, product_id, size_type)
);

-- Payment proof uploads (manual payment confirmation)
CREATE TABLE payment_proofs (
    proof_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    client_id INT NOT NULL REFERENCES client(client_id) ON DELETE CASCADE,
    payment_method payment_method_enum NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    transaction_reference VARCHAR(100),
    proof_path VARCHAR(255) NOT NULL,
    notes TEXT,
    status payment_proof_status_enum DEFAULT 'submitted',
    reviewed_by INT REFERENCES admin(admin_id),
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create updated_at triggers (since PostgreSQL doesn't have ON UPDATE CURRENT_TIMESTAMP)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to all tables that need updated_at
CREATE TRIGGER update_admin_updated_at BEFORE UPDATE ON admin FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_client_updated_at BEFORE UPDATE ON client FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_product_images_updated_at BEFORE UPDATE ON product_images FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_discounts_updated_at BEFORE UPDATE ON discounts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_cart_updated_at BEFORE UPDATE ON cart FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_client_address_updated_at BEFORE UPDATE ON client_address FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_payment_proofs_updated_at BEFORE UPDATE ON payment_proofs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();