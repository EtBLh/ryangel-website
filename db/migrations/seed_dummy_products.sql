DO $$
DECLARE
    i INT;
    new_product_id INT;
BEGIN
    FOR i IN 1..20 LOOP
        INSERT INTO products (
            product_name, 
            product_description, 
            product_type, 
            sku, 
            price, 
            quantity, 
            is_active, 
            available_sizes,
            created_at
        ) VALUES (
            'Test Faichun ' || i,
            'This is a test product for infinite scrolling',
            'faiachun',
            'TEST-FC-' || i,
            88.00,
            100,
            true,
            '{v-rect}',
            NOW()
        ) RETURNING product_id INTO new_product_id;

        INSERT INTO product_images (
            product_id,
            image_path,
            thumbnail_path,
            alt_text,
            size_type,
            sort_order,
            is_primary
        ) VALUES (
            new_product_id,
            '/api/media/products/s001.jpg',
            '/api/media/products/s001-sm.jpg',
            'Test Faichun ' || i,
            'v-rect',
            0,
            true
        );
    END LOOP;
END $$;
