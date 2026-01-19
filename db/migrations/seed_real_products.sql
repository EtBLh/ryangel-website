TRUNCATE products CASCADE;

DO $$
DECLARE
    new_id INT;
BEGIN

    -- Product: 環遊世界
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('環遊世界', 'Faichun - 環遊世界', 'faiachun', 'FC-001', 20.00, 100, true, '{v-rect,square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-環遊世界.jpg', '/api/media/products/v-rect-環遊世界-sm.jpg', '環遊世界', 'v-rect', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-環遊世界.jpg', '/api/media/products/square-環遊世界-sm.jpg', '環遊世界', 'square', 1, false);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-環遊世界.jpg', '/api/media/products/big-square-環遊世界-sm.jpg', '環遊世界', 'big-square', 2, false);

    -- Product: 家肥屋潤
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('家肥屋潤', 'Faichun - 家肥屋潤', 'faiachun', 'FC-002', 20.00, 100, true, '{v-rect,square,fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-家肥屋潤.jpg', '/api/media/products/v-rect-家肥屋潤-sm.jpg', '家肥屋潤', 'v-rect', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-家肥屋潤.jpg', '/api/media/products/square-家肥屋潤-sm.jpg', '家肥屋潤', 'square', 1, false);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-家肥屋潤.jpg', '/api/media/products/fat-v-rect-家肥屋潤-sm.jpg', '家肥屋潤', 'fat-v-rect', 2, false);

    -- Product: 家宅平安
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('家宅平安', 'Faichun - 家宅平安', 'faiachun', 'FC-003', 20.00, 100, true, '{v-rect,square,fat-v-rect,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-家宅平安.jpg', '/api/media/products/v-rect-家宅平安-sm.jpg', '家宅平安', 'v-rect', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-家宅平安.jpg', '/api/media/products/square-家宅平安-sm.jpg', '家宅平安', 'square', 1, false);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-家宅平安.jpg', '/api/media/products/fat-v-rect-家宅平安-sm.jpg', '家宅平安', 'fat-v-rect', 2, false);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-家宅平安.jpg', '/api/media/products/big-square-家宅平安-sm.jpg', '家宅平安', 'big-square', 3, false);

    -- Product: 出入平安
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('出入平安', 'Faichun - 出入平安', 'faiachun', 'FC-004', 20.00, 100, true, '{v-rect,square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-出入平安.jpg', '/api/media/products/v-rect-出入平安-sm.jpg', '出入平安', 'v-rect', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-出入平安.jpg', '/api/media/products/square-出入平安-sm.jpg', '出入平安', 'square', 1, false);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-出入平安.jpg', '/api/media/products/big-square-出入平安-sm.jpg', '出入平安', 'big-square', 2, false);

    -- Product: 煩事遠離朕 銀紙當被冚
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('煩事遠離朕 銀紙當被冚', 'Faichun - 煩事遠離朕 銀紙當被冚', 'faiachun', 'FC-005', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-煩事遠離朕 銀紙當被冚.jpg', '/api/media/products/square-煩事遠離朕 銀紙當被冚-sm.jpg', '煩事遠離朕 銀紙當被冚', 'square', 0, true);

    -- Product: 蒸蒸日上
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('蒸蒸日上', 'Faichun - 蒸蒸日上', 'faiachun', 'FC-006', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-蒸蒸日上.jpg', '/api/media/products/square-蒸蒸日上-sm.jpg', '蒸蒸日上', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-蒸蒸日上.jpg', '/api/media/products/big-square-蒸蒸日上-sm.jpg', '蒸蒸日上', 'big-square', 1, false);

    -- Product: 馬上暴富
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('馬上暴富', 'Faichun - 馬上暴富', 'faiachun', 'FC-007', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-馬上暴富.jpg', '/api/media/products/square-馬上暴富-sm.jpg', '馬上暴富', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-馬上暴富.jpg', '/api/media/products/big-square-馬上暴富-sm.jpg', '馬上暴富', 'big-square', 1, false);

    -- Product: 厄運退散
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('厄運退散', 'Faichun - 厄運退散', 'faiachun', 'FC-008', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-厄運退散.jpg', '/api/media/products/square-厄運退散-sm.jpg', '厄運退散', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-厄運退散.jpg', '/api/media/products/big-square-厄運退散-sm.jpg', '厄運退散', 'big-square', 1, false);

    -- Product: 身體Healthy 日日Happy 事事Lucky 大把Money
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('身體Healthy 日日Happy 事事Lucky 大把Money', 'Faichun - 身體Healthy 日日Happy 事事Lucky 大把Money', 'faiachun', 'FC-009', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-身體Healthy 日日Happy 事事Lucky 大把Money.jpg', '/api/media/products/square-身體Healthy 日日Happy 事事Lucky 大把Money-sm.jpg', '身體Healthy 日日Happy 事事Lucky 大把Money', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-身體Healthy 日日Happy 事事Lucky 大把Money.jpg', '/api/media/products/big-square-身體Healthy 日日Happy 事事Lucky 大把Money-sm.jpg', '身體Healthy 日日Happy 事事Lucky 大把Money', 'big-square', 1, false);

    -- Product: 福
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('福', 'Faichun - 福', 'faiachun', 'FC-010', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-福.jpg', '/api/media/products/square-福-sm.jpg', '福', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-福.jpg', '/api/media/products/big-square-福-sm.jpg', '福', 'big-square', 1, false);

    -- Product: 健康
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('健康', 'Faichun - 健康', 'faiachun', 'FC-011', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-健康.jpg', '/api/media/products/square-健康-sm.jpg', '健康', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-健康.jpg', '/api/media/products/big-square-健康-sm.jpg', '健康', 'big-square', 1, false);

    -- Product: 招財
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('招財', 'Faichun - 招財', 'faiachun', 'FC-012', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-招財.jpg', '/api/media/products/square-招財-sm.jpg', '招財', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-招財.jpg', '/api/media/products/big-square-招財-sm.jpg', '招財', 'big-square', 1, false);

    -- Product: 納福
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('納福', 'Faichun - 納福', 'faiachun', 'FC-013', 20.00, 100, true, '{square,big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-納福.jpg', '/api/media/products/square-納福-sm.jpg', '納福', 'square', 0, true);

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-納福.jpg', '/api/media/products/big-square-納福-sm.jpg', '納福', 'big-square', 1, false);

    -- Product: 招財豬堡
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('招財豬堡', 'Faichun - 招財豬堡', 'faiachun', 'FC-014', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-招財豬堡.jpg', '/api/media/products/v-rect-招財豬堡-sm.jpg', '招財豬堡', 'v-rect', 0, true);

    -- Product: 見字抖抖
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('見字抖抖', 'Faichun - 見字抖抖', 'faiachun', 'FC-015', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-見字抖抖.jpg', '/api/media/products/v-rect-見字抖抖-sm.jpg', '見字抖抖', 'v-rect', 0, true);

    -- Product: 升職加薪
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('升職加薪', 'Faichun - 升職加薪', 'faiachun', 'FC-016', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-升職加薪.jpg', '/api/media/products/v-rect-升職加薪-sm.jpg', '升職加薪', 'v-rect', 0, true);

    -- Product: 大吉大利
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('大吉大利', 'Faichun - 大吉大利', 'faiachun', 'FC-017', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-大吉大利.jpg', '/api/media/products/v-rect-大吉大利-sm.jpg', '大吉大利', 'v-rect', 0, true);

    -- Product: 逆風翻盤
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('逆風翻盤', 'Faichun - 逆風翻盤', 'faiachun', 'FC-018', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-逆風翻盤.jpg', '/api/media/products/v-rect-逆風翻盤-sm.jpg', '逆風翻盤', 'v-rect', 0, true);

    -- Product: 情緒穩定
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('情緒穩定', 'Faichun - 情緒穩定', 'faiachun', 'FC-019', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-情緒穩定.jpg', '/api/media/products/v-rect-情緒穩定-sm.jpg', '情緒穩定', 'v-rect', 0, true);

    -- Product: 莫生氣
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('莫生氣', 'Faichun - 莫生氣', 'faiachun', 'FC-020', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-莫生氣.jpg', '/api/media/products/v-rect-莫生氣-sm.jpg', '莫生氣', 'v-rect', 0, true);

    -- Product: 準時收工
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('準時收工', 'Faichun - 準時收工', 'faiachun', 'FC-021', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-準時收工.jpg', '/api/media/products/v-rect-準時收工-sm.jpg', '準時收工', 'v-rect', 0, true);

    -- Product: 蒸神爽利
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('蒸神爽利', 'Faichun - 蒸神爽利', 'faiachun', 'FC-022', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-蒸神爽利.jpg', '/api/media/products/v-rect-蒸神爽利-sm.jpg', '蒸神爽利', 'v-rect', 0, true);

    -- Product: 大把錢洗
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('大把錢洗', 'Faichun - 大把錢洗', 'faiachun', 'FC-023', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-大把錢洗.jpg', '/api/media/products/v-rect-大把錢洗-sm.jpg', '大把錢洗', 'v-rect', 0, true);

    -- Product: 馬上戀愛
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('馬上戀愛', 'Faichun - 馬上戀愛', 'faiachun', 'FC-024', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-馬上戀愛.jpg', '/api/media/products/v-rect-馬上戀愛-sm.jpg', '馬上戀愛', 'v-rect', 0, true);

    -- Product: 心平氣和
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('心平氣和', 'Faichun - 心平氣和', 'faiachun', 'FC-025', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-心平氣和.jpg', '/api/media/products/v-rect-心平氣和-sm.jpg', '心平氣和', 'v-rect', 0, true);

    -- Product: 珠圓肉潤
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('珠圓肉潤', 'Faichun - 珠圓肉潤', 'faiachun', 'FC-026', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-珠圓肉潤.jpg', '/api/media/products/v-rect-珠圓肉潤-sm.jpg', '珠圓肉潤', 'v-rect', 0, true);

    -- Product: 新奇可愛
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('新奇可愛', 'Faichun - 新奇可愛', 'faiachun', 'FC-027', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-新奇可愛.jpg', '/api/media/products/v-rect-新奇可愛-sm.jpg', '新奇可愛', 'v-rect', 0, true);

    -- Product: 身體健康
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('身體健康', 'Faichun - 身體健康', 'faiachun', 'FC-028', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-身體健康.jpg', '/api/media/products/v-rect-身體健康-sm.jpg', '身體健康', 'v-rect', 0, true);

    -- Product: 肌肉暴漲
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('肌肉暴漲', 'Faichun - 肌肉暴漲', 'faiachun', 'FC-029', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-肌肉暴漲.jpg', '/api/media/products/v-rect-肌肉暴漲-sm.jpg', '肌肉暴漲', 'v-rect', 0, true);

    -- Product: 保持運動
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('保持運動', 'Faichun - 保持運動', 'faiachun', 'FC-030', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-保持運動.jpg', '/api/media/products/v-rect-保持運動-sm.jpg', '保持運動', 'v-rect', 0, true);

    -- Product: 財富自由
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('財富自由', 'Faichun - 財富自由', 'faiachun', 'FC-031', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-財富自由.jpg', '/api/media/products/v-rect-財富自由-sm.jpg', '財富自由', 'v-rect', 0, true);

    -- Product: 無焦無慮
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('無焦無慮', 'Faichun - 無焦無慮', 'faiachun', 'FC-032', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-無焦無慮.jpg', '/api/media/products/v-rect-無焦無慮-sm.jpg', '無焦無慮', 'v-rect', 0, true);

    -- Product: 唔L洗做
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('唔L洗做', 'Faichun - 唔L洗做', 'faiachun', 'FC-033', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-唔L洗做.jpg', '/api/media/products/v-rect-唔L洗做-sm.jpg', '唔L洗做', 'v-rect', 0, true);

    -- Product: 心想事成
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('心想事成', 'Faichun - 心想事成', 'faiachun', 'FC-034', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-心想事成.jpg', '/api/media/products/v-rect-心想事成-sm.jpg', '心想事成', 'v-rect', 0, true);

    -- Product: 發過豬頭
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('發過豬頭', 'Faichun - 發過豬頭', 'faiachun', 'FC-035', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-發過豬頭.jpg', '/api/media/products/v-rect-發過豬頭-sm.jpg', '發過豬頭', 'v-rect', 0, true);

    -- Product: 心花怒放
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('心花怒放', 'Faichun - 心花怒放', 'faiachun', 'FC-036', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-心花怒放.jpg', '/api/media/products/v-rect-心花怒放-sm.jpg', '心花怒放', 'v-rect', 0, true);

    -- Product: 不勞而獲
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('不勞而獲', 'Faichun - 不勞而獲', 'faiachun', 'FC-037', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-不勞而獲.jpg', '/api/media/products/v-rect-不勞而獲-sm.jpg', '不勞而獲', 'v-rect', 0, true);

    -- Product: 躺贏
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('躺贏', 'Faichun - 躺贏', 'faiachun', 'FC-038', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-躺贏.jpg', '/api/media/products/v-rect-躺贏-sm.jpg', '躺贏', 'v-rect', 0, true);

    -- Product: 快樂長大
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('快樂長大', 'Faichun - 快樂長大', 'faiachun', 'FC-039', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-快樂長大.jpg', '/api/media/products/v-rect-快樂長大-sm.jpg', '快樂長大', 'v-rect', 0, true);

    -- Product: 一索得寶
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('一索得寶', 'Faichun - 一索得寶', 'faiachun', 'FC-040', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-一索得寶.jpg', '/api/media/products/v-rect-一索得寶-sm.jpg', '一索得寶', 'v-rect', 0, true);

    -- Product: 馬到成功
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('馬到成功', 'Faichun - 馬到成功', 'faiachun', 'FC-041', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-馬到成功.jpg', '/api/media/products/v-rect-馬到成功-sm.jpg', '馬到成功', 'v-rect', 0, true);

    -- Product: 天降大運
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('天降大運', 'Faichun - 天降大運', 'faiachun', 'FC-042', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-天降大運.jpg', '/api/media/products/v-rect-天降大運-sm.jpg', '天降大運', 'v-rect', 0, true);

    -- Product: 無病息災
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('無病息災', 'Faichun - 無病息災', 'faiachun', 'FC-043', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-無病息災.jpg', '/api/media/products/v-rect-無病息災-sm.jpg', '無病息災', 'v-rect', 0, true);

    -- Product: 馬上行大運
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('馬上行大運', 'Faichun - 馬上行大運', 'faiachun', 'FC-044', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-馬上行大運.jpg', '/api/media/products/v-rect-馬上行大運-sm.jpg', '馬上行大運', 'v-rect', 0, true);

    -- Product: 啡嘗快樂
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('啡嘗快樂', 'Faichun - 啡嘗快樂', 'faiachun', 'FC-045', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-啡嘗快樂.jpg', '/api/media/products/v-rect-啡嘗快樂-sm.jpg', '啡嘗快樂', 'v-rect', 0, true);

    -- Product: 快樂加馬
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('快樂加馬', 'Faichun - 快樂加馬', 'faiachun', 'FC-046', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-快樂加馬.jpg', '/api/media/products/v-rect-快樂加馬-sm.jpg', '快樂加馬', 'v-rect', 0, true);

    -- Product: 睡滿滿
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('睡滿滿', 'Faichun - 睡滿滿', 'faiachun', 'FC-047', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-睡滿滿.jpg', '/api/media/products/v-rect-睡滿滿-sm.jpg', '睡滿滿', 'v-rect', 0, true);

    -- Product: 豚豚圓圓
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('豚豚圓圓', 'Faichun - 豚豚圓圓', 'faiachun', 'FC-048', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-豚豚圓圓.jpg', '/api/media/products/v-rect-豚豚圓圓-sm.jpg', '豚豚圓圓', 'v-rect', 0, true);

    -- Product: 財福兩旺
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('財福兩旺', 'Faichun - 財福兩旺', 'faiachun', 'FC-049', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-財福兩旺.jpg', '/api/media/products/v-rect-財福兩旺-sm.jpg', '財福兩旺', 'v-rect', 0, true);

    -- Product: 悠然自得
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('悠然自得', 'Faichun - 悠然自得', 'faiachun', 'FC-050', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-悠然自得.jpg', '/api/media/products/v-rect-悠然自得-sm.jpg', '悠然自得', 'v-rect', 0, true);

    -- Product: 四季平安
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('四季平安', 'Faichun - 四季平安', 'faiachun', 'FC-051', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-四季平安.jpg', '/api/media/products/v-rect-四季平安-sm.jpg', '四季平安', 'v-rect', 0, true);

    -- Product: 幸福滿屋
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('幸福滿屋', 'Faichun - 幸福滿屋', 'faiachun', 'FC-052', 20.00, 100, true, '{v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/v-rect-幸福滿屋.jpg', '/api/media/products/v-rect-幸福滿屋-sm.jpg', '幸福滿屋', 'v-rect', 0, true);

    -- Product: 勿多勞
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('勿多勞', 'Faichun - 勿多勞', 'faiachun', 'FC-053', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-勿多勞.jpg', '/api/media/products/square-勿多勞-sm.jpg', '勿多勞', 'square', 0, true);

    -- Product: 學業進步 精叻醒目
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('學業進步 精叻醒目', 'Faichun - 學業進步 精叻醒目', 'faiachun', 'FC-054', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-學業進步 精叻醒目.jpg', '/api/media/products/square-學業進步 精叻醒目-sm.jpg', '學業進步 精叻醒目', 'square', 0, true);

    -- Product: 磚心上班 開心下班
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('磚心上班 開心下班', 'Faichun - 磚心上班 開心下班', 'faiachun', 'FC-055', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-磚心上班 開心下班.jpg', '/api/media/products/square-磚心上班 開心下班-sm.jpg', '磚心上班 開心下班', 'square', 0, true);

    -- Product: 吉星高照 貴人超吉多
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('吉星高照 貴人超吉多', 'Faichun - 吉星高照 貴人超吉多', 'faiachun', 'FC-056', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-吉星高照 貴人超吉多.jpg', '/api/media/products/square-吉星高照 貴人超吉多-sm.jpg', '吉星高照 貴人超吉多', 'square', 0, true);

    -- Product: 我命由我不由天
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('我命由我不由天', 'Faichun - 我命由我不由天', 'faiachun', 'FC-057', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-我命由我不由天.jpg', '/api/media/products/square-我命由我不由天-sm.jpg', '我命由我不由天', 'square', 0, true);

    -- Product: 心中有佛 淡定老實
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('心中有佛 淡定老實', 'Faichun - 心中有佛 淡定老實', 'faiachun', 'FC-058', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-心中有佛 淡定老實.jpg', '/api/media/products/square-心中有佛 淡定老實-sm.jpg', '心中有佛 淡定老實', 'square', 0, true);

    -- Product: 收工即走 盡忠職守
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('收工即走 盡忠職守', 'Faichun - 收工即走 盡忠職守', 'faiachun', 'FC-059', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-收工即走 盡忠職守.jpg', '/api/media/products/square-收工即走 盡忠職守-sm.jpg', '收工即走 盡忠職守', 'square', 0, true);

    -- Product: 長賺長有 有車有樓
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('長賺長有 有車有樓', 'Faichun - 長賺長有 有車有樓', 'faiachun', 'FC-060', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-長賺長有 有車有樓.jpg', '/api/media/products/square-長賺長有 有車有樓-sm.jpg', '長賺長有 有車有樓', 'square', 0, true);

    -- Product: 幸福加馬 愛神上線
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('幸福加馬 愛神上線', 'Faichun - 幸福加馬 愛神上線', 'faiachun', 'FC-061', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-幸福加馬 愛神上線.jpg', '/api/media/products/square-幸福加馬 愛神上線-sm.jpg', '幸福加馬 愛神上線', 'square', 0, true);

    -- Product: 情緒穩定 保持微笑
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('情緒穩定 保持微笑', 'Faichun - 情緒穩定 保持微笑', 'faiachun', 'FC-062', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-情緒穩定 保持微笑.jpg', '/api/media/products/square-情緒穩定 保持微笑-sm.jpg', '情緒穩定 保持微笑', 'square', 0, true);

    -- Product: 見字抖抖 廢水自由
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('見字抖抖 廢水自由', 'Faichun - 見字抖抖 廢水自由', 'faiachun', 'FC-063', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-見字抖抖 廢水自由.jpg', '/api/media/products/square-見字抖抖 廢水自由-sm.jpg', '見字抖抖 廢水自由', 'square', 0, true);

    -- Product: 天賦異稟 才華爆燈
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('天賦異稟 才華爆燈', 'Faichun - 天賦異稟 才華爆燈', 'faiachun', 'FC-064', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-天賦異稟 才華爆燈.jpg', '/api/media/products/square-天賦異稟 才華爆燈-sm.jpg', '天賦異稟 才華爆燈', 'square', 0, true);

    -- Product: 身體健康 平安喜樂
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('身體健康 平安喜樂', 'Faichun - 身體健康 平安喜樂', 'faiachun', 'FC-065', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-身體健康 平安喜樂.jpg', '/api/media/products/square-身體健康 平安喜樂-sm.jpg', '身體健康 平安喜樂', 'square', 0, true);

    -- Product: 食極唔肥 八舊腹肌
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('食極唔肥 八舊腹肌', 'Faichun - 食極唔肥 八舊腹肌', 'faiachun', 'FC-066', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-食極唔肥 八舊腹肌.jpg', '/api/media/products/square-食極唔肥 八舊腹肌-sm.jpg', '食極唔肥 八舊腹肌', 'square', 0, true);

    -- Product: 百毒不侵
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('百毒不侵', 'Faichun - 百毒不侵', 'faiachun', 'FC-067', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-百毒不侵.jpg', '/api/media/products/square-百毒不侵-sm.jpg', '百毒不侵', 'square', 0, true);

    -- Product: 生意興隆 撈乜都掂
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('生意興隆 撈乜都掂', 'Faichun - 生意興隆 撈乜都掂', 'faiachun', 'FC-068', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-生意興隆 撈乜都掂.jpg', '/api/media/products/square-生意興隆 撈乜都掂-sm.jpg', '生意興隆 撈乜都掂', 'square', 0, true);

    -- Product: 焦慮退散 淡淡定 有錢剩
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('焦慮退散 淡淡定 有錢剩', 'Faichun - 焦慮退散 淡淡定 有錢剩', 'faiachun', 'FC-069', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-焦慮退散 淡淡定 有錢剩.jpg', '/api/media/products/square-焦慮退散 淡淡定 有錢剩-sm.jpg', '焦慮退散 淡淡定 有錢剩', 'square', 0, true);

    -- Product: 無驚無險 又到飯點
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('無驚無險 又到飯點', 'Faichun - 無驚無險 又到飯點', 'faiachun', 'FC-070', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-無驚無險 又到飯點.jpg', '/api/media/products/square-無驚無險 又到飯點-sm.jpg', '無驚無險 又到飯點', 'square', 0, true);

    -- Product: 所願皆成 得償所願
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('所願皆成 得償所願', 'Faichun - 所願皆成 得償所願', 'faiachun', 'FC-071', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-所願皆成 得償所願.jpg', '/api/media/products/square-所願皆成 得償所願-sm.jpg', '所願皆成 得償所願', 'square', 0, true);

    -- Product: 順從自己 Take It Easy
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('順從自己 Take It Easy', 'Faichun - 順從自己 Take It Easy', 'faiachun', 'FC-072', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-順從自己 Take It Easy.jpg', '/api/media/products/square-順從自己 Take It Easy-sm.jpg', '順從自己 Take It Easy', 'square', 0, true);

    -- Product: 掂過碌蔗 起身就HEA
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('掂過碌蔗 起身就HEA', 'Faichun - 掂過碌蔗 起身就HEA', 'faiachun', 'FC-073', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-掂過碌蔗 起身就HEA.jpg', '/api/media/products/square-掂過碌蔗 起身就HEA-sm.jpg', '掂過碌蔗 起身就HEA', 'square', 0, true);

    -- Product: 吃飽飽 睡香香
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('吃飽飽 睡香香', 'Faichun - 吃飽飽 睡香香', 'faiachun', 'FC-074', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-吃飽飽 睡香香.jpg', '/api/media/products/square-吃飽飽 睡香香-sm.jpg', '吃飽飽 睡香香', 'square', 0, true);

    -- Product: 健康成長 日日開心
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('健康成長 日日開心', 'Faichun - 健康成長 日日開心', 'faiachun', 'FC-075', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-健康成長 日日開心.jpg', '/api/media/products/square-健康成長 日日開心-sm.jpg', '健康成長 日日開心', 'square', 0, true);

    -- Product: 牛馬精神 富貴逼人
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('牛馬精神 富貴逼人', 'Faichun - 牛馬精神 富貴逼人', 'faiachun', 'FC-076', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-牛馬精神 富貴逼人.jpg', '/api/media/products/square-牛馬精神 富貴逼人-sm.jpg', '牛馬精神 富貴逼人', 'square', 0, true);

    -- Product: 愛要及時 食要隨時
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('愛要及時 食要隨時', 'Faichun - 愛要及時 食要隨時', 'faiachun', 'FC-077', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-愛要及時 食要隨時.jpg', '/api/media/products/square-愛要及時 食要隨時-sm.jpg', '愛要及時 食要隨時', 'square', 0, true);

    -- Product: 馬上好運
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('馬上好運', 'Faichun - 馬上好運', 'faiachun', 'FC-078', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-馬上好運.jpg', '/api/media/products/square-馬上好運-sm.jpg', '馬上好運', 'square', 0, true);

    -- Product: 無憂無慮 HAKUNA MATATA
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('無憂無慮 HAKUNA MATATA', 'Faichun - 無憂無慮 HAKUNA MATATA', 'faiachun', 'FC-079', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-無憂無慮 HAKUNA MATATA.jpg', '/api/media/products/square-無憂無慮 HAKUNA MATATA-sm.jpg', '無憂無慮 HAKUNA MATATA', 'square', 0, true);

    -- Product: 馬上升職
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('馬上升職', 'Faichun - 馬上升職', 'faiachun', 'FC-080', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-馬上升職.jpg', '/api/media/products/square-馬上升職-sm.jpg', '馬上升職', 'square', 0, true);

    -- Product: 日日是好日 好好生活
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('日日是好日 好好生活', 'Faichun - 日日是好日 好好生活', 'faiachun', 'FC-081', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-日日是好日 好好生活.jpg', '/api/media/products/square-日日是好日 好好生活-sm.jpg', '日日是好日 好好生活', 'square', 0, true);

    -- Product: 咖啡在手 日日富有
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('咖啡在手 日日富有', 'Faichun - 咖啡在手 日日富有', 'faiachun', 'FC-082', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-咖啡在手 日日富有.jpg', '/api/media/products/square-咖啡在手 日日富有-sm.jpg', '咖啡在手 日日富有', 'square', 0, true);

    -- Product: 幸福加馬
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('幸福加馬', 'Faichun - 幸福加馬', 'faiachun', 'FC-083', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-幸福加馬.jpg', '/api/media/products/square-幸福加馬-sm.jpg', '幸福加馬', 'square', 0, true);

    -- Product: 瞓到自然醒
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('瞓到自然醒', 'Faichun - 瞓到自然醒', 'faiachun', 'FC-084', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-瞓到自然醒.jpg', '/api/media/products/square-瞓到自然醒-sm.jpg', '瞓到自然醒', 'square', 0, true);

    -- Product: 幸福到家 相親相愛
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('幸福到家 相親相愛', 'Faichun - 幸福到家 相親相愛', 'faiachun', 'FC-085', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-幸福到家 相親相愛.jpg', '/api/media/products/square-幸福到家 相親相愛-sm.jpg', '幸福到家 相親相愛', 'square', 0, true);

    -- Product: 銀行有錢 心中有愛
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('銀行有錢 心中有愛', 'Faichun - 銀行有錢 心中有愛', 'faiachun', 'FC-086', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-銀行有錢 心中有愛.jpg', '/api/media/products/square-銀行有錢 心中有愛-sm.jpg', '銀行有錢 心中有愛', 'square', 0, true);

    -- Product: 情深深雨濛濛 遇到桃花眼唔矇
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('情深深雨濛濛 遇到桃花眼唔矇', 'Faichun - 情深深雨濛濛 遇到桃花眼唔矇', 'faiachun', 'FC-087', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-情深深雨濛濛 遇到桃花眼唔矇.jpg', '/api/media/products/square-情深深雨濛濛 遇到桃花眼唔矇-sm.jpg', '情深深雨濛濛 遇到桃花眼唔矇', 'square', 0, true);

    -- Product: 食好瞓好 無曬煩惱
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('食好瞓好 無曬煩惱', 'Faichun - 食好瞓好 無曬煩惱', 'faiachun', 'FC-088', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-食好瞓好 無曬煩惱.jpg', '/api/media/products/square-食好瞓好 無曬煩惱-sm.jpg', '食好瞓好 無曬煩惱', 'square', 0, true);

    -- Product: 有吃有睡 人生不累
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('有吃有睡 人生不累', 'Faichun - 有吃有睡 人生不累', 'faiachun', 'FC-089', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-有吃有睡 人生不累.jpg', '/api/media/products/square-有吃有睡 人生不累-sm.jpg', '有吃有睡 人生不累', 'square', 0, true);

    -- Product: 一步一腳印 健康成長
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('一步一腳印 健康成長', 'Faichun - 一步一腳印 健康成長', 'faiachun', 'FC-090', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-一步一腳印 健康成長.jpg', '/api/media/products/square-一步一腳印 健康成長-sm.jpg', '一步一腳印 健康成長', 'square', 0, true);

    -- Product: 披荊斬棘 一家人最緊要齊齊整整
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('披荊斬棘 一家人最緊要齊齊整整', 'Faichun - 披荊斬棘 一家人最緊要齊齊整整', 'faiachun', 'FC-091', 20.00, 100, true, '{square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/square-披荊斬棘 一家人最緊要齊齊整整.jpg', '/api/media/products/square-披荊斬棘 一家人最緊要齊齊整整-sm.jpg', '披荊斬棘 一家人最緊要齊齊整整', 'square', 0, true);

    -- Product: 豬事皆順
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('豬事皆順', 'Faichun - 豬事皆順', 'faiachun', 'FC-092', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-豬事皆順.jpg', '/api/media/products/fat-v-rect-豬事皆順-sm.jpg', '豬事皆順', 'fat-v-rect', 0, true);

    -- Product: 出入平安
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('出入平安', 'Faichun - 出入平安', 'faiachun', 'FC-093', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-出入平安.jpg', '/api/media/products/fat-v-rect-出入平安-sm.jpg', '出入平安', 'fat-v-rect', 0, true);

    -- Product: 開門大吉
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('開門大吉', 'Faichun - 開門大吉', 'faiachun', 'FC-094', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-開門大吉.jpg', '/api/media/products/fat-v-rect-開門大吉-sm.jpg', '開門大吉', 'fat-v-rect', 0, true);

    -- Product: 平安喜樂
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('平安喜樂', 'Faichun - 平安喜樂', 'faiachun', 'FC-095', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-平安喜樂.jpg', '/api/media/products/fat-v-rect-平安喜樂-sm.jpg', '平安喜樂', 'fat-v-rect', 0, true);

    -- Product: 招財進寶
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('招財進寶', 'Faichun - 招財進寶', 'faiachun', 'FC-096', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-招財進寶.jpg', '/api/media/products/fat-v-rect-招財進寶-sm.jpg', '招財進寶', 'fat-v-rect', 0, true);

    -- Product: 諸邪莫近
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('諸邪莫近', 'Faichun - 諸邪莫近', 'faiachun', 'FC-097', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-諸邪莫近.jpg', '/api/media/products/fat-v-rect-諸邪莫近-sm.jpg', '諸邪莫近', 'fat-v-rect', 0, true);

    -- Product: 鵝福臨門
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('鵝福臨門', 'Faichun - 鵝福臨門', 'faiachun', 'FC-098', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-鵝福臨門.jpg', '/api/media/products/fat-v-rect-鵝福臨門-sm.jpg', '鵝福臨門', 'fat-v-rect', 0, true);

    -- Product: 發財暴富
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('發財暴富', 'Faichun - 發財暴富', 'faiachun', 'FC-099', 20.00, 100, true, '{fat-v-rect}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/fat-v-rect-發財暴富.jpg', '/api/media/products/fat-v-rect-發財暴富-sm.jpg', '發財暴富', 'fat-v-rect', 0, true);

    -- Product: 生意興隆 撈乜都掂
    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)
    VALUES ('生意興隆 撈乜都掂', 'Faichun - 生意興隆 撈乜都掂', 'faiachun', 'FC-100', 20.00, 100, true, '{big-square}', NOW())
    RETURNING product_id INTO new_id;

    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)
    VALUES (new_id, '/api/media/products/big-square-生意興隆 撈乜都掂.jpg', '/api/media/products/big-square-生意興隆 撈乜都掂-sm.jpg', '生意興隆 撈乜都掂', 'big-square', 0, true);
END $$;