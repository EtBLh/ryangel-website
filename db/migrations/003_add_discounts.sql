ALTER TABLE discounts ADD COLUMN is_auto_apply BOOLEAN DEFAULT FALSE;

-- Seed B3G1 Faichun
INSERT INTO discounts (
    discount_name, discount_code, discount_type, 
    buy_quantity, get_quantity, product_type_restriction, 
    is_auto_apply, start_date, end_date, applies_to, is_active
) VALUES (
    'Buy 3 Get 1 Free Faichun', 'AUTO_B3G1_FAICHUN', 'bxgy',
    3, 1, 'faiachun', 
    TRUE, NOW(), NOW() + INTERVAL '10 years', 'specific_categories', TRUE
);

-- Seed Free Shipping
INSERT INTO discounts (
    discount_name, discount_code, discount_type,
    discount_value, is_auto_apply, start_date, end_date, applies_to, is_active
) VALUES (
    'Free Shipping', 'AUTO_FREE_SHIPPING', 'free_shipping',
    0, TRUE, NOW(), NOW() + INTERVAL '10 years', 'all_products', TRUE
);
