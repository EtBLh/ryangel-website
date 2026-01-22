import json

def generate_sql(json_path, output_path):
    with open(json_path, 'r', encoding='utf-8') as f:
        products = json.load(f)

    with open(output_path, 'w', encoding='utf-8') as f:
        f.write("TRUNCATE products CASCADE;\n\n")
        f.write("DO $$\n")
        f.write("DECLARE\n")
        f.write("    new_id INT;\n")
        f.write("BEGIN\n\n")

        for i, product in enumerate(products):
            name = product['name'].replace("'", "''")
            sizes = product['size']
            sizes_str = '{' + ','.join(sizes) + '}'
            sku = f"FC-{i+1:03d}"
            
            f.write(f"    -- Product: {name}\n")
            f.write(f"    INSERT INTO products (product_name, product_description, product_type, sku, price, quantity, is_active, available_sizes, created_at)\n")
            f.write(f"    VALUES ('{name}', 'Faichun - {name}', 'faiachun', '{sku}', 20.00, 100, true, '{sizes_str}', NOW())\n")
            f.write(f"    RETURNING product_id INTO new_id;\n\n")

            for j, size in enumerate(sizes):
                image_path = f"/api/media/products/{size}-{name}.jpg"
                thumbnail_path = f"/api/media/products/{size}-{name}-sm.jpg"
                is_primary = "true" if j == 0 else "false"
                
                f.write(f"    INSERT INTO product_images (product_id, image_path, thumbnail_path, alt_text, size_type, sort_order, is_primary)\n")
                f.write(f"    VALUES (new_id, '{image_path}', '{thumbnail_path}', '{name}', '{size}', {j}, {is_primary});\n\n")

        f.write("END $$;\n")

if __name__ == "__main__":
    generate_sql('products.json', 'db/migrations/seed_from_json.sql')
