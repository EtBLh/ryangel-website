export type SizeType = 'v-rect' | 'square' | 'fat-v-rect';

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
  product_id: number;
  size_type: SizeType | null;
  quantity: number;
  added_at: string;
  product_name: string;
  unit_price: number;
  stock_quantity: number;
}

export interface Cart {
  items: CartItem[];
  subtotal: number;
  discount: number;
  total: number;
}