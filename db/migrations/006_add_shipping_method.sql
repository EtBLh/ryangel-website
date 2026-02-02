-- Add shipping_method enum and field to orders table
CREATE TYPE shipping_method_enum AS ENUM ('ebuy_store', 'direct_address');

ALTER TABLE orders ADD COLUMN shipping_method shipping_method_enum DEFAULT 'ebuy_store';

-- Update existing orders to have ebuy_store as default method
UPDATE orders SET shipping_method = 'ebuy_store' WHERE ebuy_store_id IS NOT NULL;

-- Drop old constraint that required either shipping_address_id or ebuy_store_id
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_shipping_address_check;

-- Add new constraint: for ebuy_store method, ebuy_store_id is required
-- For direct_address, neither ebuy_store_id nor shipping_address_id should be set
ALTER TABLE orders ADD CONSTRAINT orders_shipping_method_check CHECK (
    (shipping_method = 'ebuy_store' AND ebuy_store_id IS NOT NULL AND shipping_address_id IS NULL) OR
    (shipping_method = 'direct_address' AND ebuy_store_id IS NULL AND shipping_address_id IS NULL)
);

-- Add new free shipping discount for direct address orders over $120
INSERT INTO discounts (
    discount_name, discount_code, discount_type,
    discount_value, minimum_order_amount, is_auto_apply,
    start_date, end_date, applies_to, is_active
) VALUES (
    'Free Shipping for Direct Address ($120+)', 'AUTO_FREE_SHIPPING_DIRECT', 'free_shipping',
    0, 120.00, TRUE,
    NOW(), NOW() + INTERVAL '10 years', 'all_products', TRUE
);
