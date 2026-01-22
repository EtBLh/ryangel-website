export type SizeType = 'v-rect' | 'square' | 'fat-v-rect' | 'big-square';

export interface Product {
  product_id: number;
  product_name: string;
  product_description: string;
  product_type: string;
  hashtag: string | null;
  sku: string;
  price: number;
  compare_at_price: number | null;
  cost_price: number | null;
  quantity: number;
  weight: number | null;
  dimensions: string | null;
  is_featured: boolean;
  is_active: boolean;
  available_sizes: SizeType[];
  seo_title: string | null;
  seo_description: string | null;
  tags: string | null;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  images: ProductImage[];
  categories: ProductCategory[];
}

export interface ProductImage {
  image_id: number;
  product_id: number;
  url: string;
  thumbnail_url: string;
  alt_text: string;
  size_type: SizeType | null;
  is_primary: boolean;
  sort_order: number;
  created_at: string;
}

export interface ProductCategory {
  category_id: number;
  category_name: string;
  category_description: string | null;
  is_active: boolean;
}

export interface CartItem {
  cart_item_id: number;
  product_id: number;
  size_type: SizeType | null;
  quantity: number;
  added_at: string;
  product_name: string;
  unit_price: number;
  stock_quantity: number;
  thumbnail_url: string;
}

export interface Cart {
  items: CartItem[];
  subtotal: number;
  discounted_subtotal: number;
  shipping_fee: number;
  discounted_shipping_fee: number;
  discount: number;
  total: number;
}
export interface Order {
  order_id: number;
  order_number: string;
  order_status: string;
  subtotal_amount: number;
  discount_amount: number;
  shipping_amount: number;
  tax_amount: number;
  total_amount: number;
  payment_method: string;
  payment_status: string;
  order_date: string;
  payment_proof?: string | null;
  ebuy_store_name?: string;
}

export interface OrderItem {
  order_item_id: number;
  product_name: string;
  product_sku: string;
  quantity: number;
  unit_price: number;
  total_price: number;
  product_image?: string;
  size_type?: string;
}

export interface OrderWithItems {
    order: Order;
    items: OrderItem[];
}

export interface EbuyStore {
  store_id: string;
  store_name: string;
  type: string;
  office_hours: string;
  address: string;
  address_en: string;
  latitude: number;
  longitude: number;
}

export interface Client {
  client_id: number;
  email: string | null;
  username: string | null;
  phone: string;
  is_active: boolean;
  activated: boolean;
  date_of_birth?: string | null;
}
