-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Migrate cart table
ALTER TABLE cart DROP CONSTRAINT IF EXISTS cart_pkey CASCADE;
ALTER TABLE cart ALTER COLUMN cart_id DROP DEFAULT;
ALTER TABLE cart ALTER COLUMN cart_id SET DATA TYPE UUID USING (uuid_generate_v4());
ALTER TABLE cart ALTER COLUMN cart_id SET DEFAULT uuid_generate_v4();
ALTER TABLE cart ADD PRIMARY KEY (cart_id);

-- Migrate cart_items table
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS cart_items_cart_id_fkey;
ALTER TABLE cart_items ALTER COLUMN cart_id SET DATA TYPE UUID USING (cart_id::text::uuid);
ALTER TABLE cart_items ADD CONSTRAINT cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES cart(cart_id) ON DELETE CASCADE;
