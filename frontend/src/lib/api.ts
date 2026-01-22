import { logout } from '../store/authSlice';
import { clearCart } from '../store/cartSlice';
import { adminLogout } from '../store/adminAuthSlice';

import axios from 'axios';
import { store } from '../store';
import { toast } from 'sonner';

const API_BASE = (import.meta.env.VITE_API_ROOT || 'https://ryangel.com/api');

interface ApiEndpoint {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  path: string;
  requiresCartId?: boolean;
  requiresAuth?: boolean;
  requiresAdminAuth?: boolean;
}

export const api: Record<string, ApiEndpoint> = {
  // Admin Endpoints
  adminLogin: {
    method: 'POST',
    path: '/admin/login'
  },
  adminMe: {
    method: 'GET',
    path: '/admin/me',
    requiresAdminAuth: true
  },
  adminLogout: {
    method: 'POST',
    path: '/admin/logout',
    requiresAdminAuth: true
  },
  adminGetStats: {
    method: 'GET',
    path: '/admin/orders/stats',
    requiresAdminAuth: true
  },
  adminGetOrders: {
    method: 'GET',
    path: '/admin/orders',
    requiresAdminAuth: true
  },
  adminGetOrderItems: {
    method: 'GET',
    path: '/admin/orders/:orderId/items',
    requiresAdminAuth: true
  },
  adminUpdateOrderStatus: {
    method: 'PATCH',
    path: '/admin/orders/:orderId/status',
    requiresAdminAuth: true
  },

  getProducts: {
    method: 'GET',
    path: '/products'
  },
  getProduct: {
    method: 'GET',
    path: '/products/:productId'
  },
  addToCart: {
    method: 'POST',
    path: '/cart/items',
    requiresCartId: true,
    requiresAuth: true
  },
  getCart: {
    method: 'GET',
    path: '/cart',
    requiresCartId: true,
    requiresAuth: true
  },
  updateCartItem: {
    method: 'PATCH',
    path: '/cart/items/:cartItemId',
    requiresAuth: true
  },
  removeCartItem: {
    method: 'DELETE',
    path: '/cart/items/:cartItemId',
    requiresAuth: true
  },
  getEbuyStores: {
    method: 'GET',
    path: '/ebuystores'
  },
  createOrder: {
    method: 'POST',
    path: '/orders',
    requiresAuth: true,
    requiresCartId: true
  },
  // Auth endpoints
  clientRegister: {
    method: 'POST',
    path: '/clients/register'
  },
  clientLogin: {
    method: 'POST',
    path: '/clients/login'
  },
  verifyOTP: {
    method: 'POST',
    path: '/clients/verify-otp'
  },
  getOrders: {
    method: 'GET',
    path: '/orders',
    requiresAuth: true
  },
  resendOTP: {
    method: 'POST',
    path: '/clients/resend-otp'
  },
  clientMe: {
    method: 'GET',
    path: '/clients/me',
    requiresAuth: true
  },
  updateClient: {
    method: 'PATCH',
    path: '/clients/me',
    requiresAuth: true
  },
  clientLogout: {
    method: 'POST',
    path: '/clients/logout',
    requiresAuth: true
  }
};

export const callAPI = async (
  endpoint: keyof typeof api,
  params?: Record<string, string | number>,
  data?: any
) => {
  const { method, path, requiresCartId } = api[endpoint];
  const endpointConfig = api[endpoint] as ApiEndpoint & { requiresAuth?: boolean; requiresAdminAuth?: boolean };
  let finalPath = path;
  const queryParams: Record<string, any> = {};

  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (finalPath.includes(`:${key}`)) {
        finalPath = finalPath.replace(`:${key}`, String(value));
      } else {
        queryParams[key] = value;
      }
    });
  }

  const headers: Record<string, string> = {};

  const state = store.getState();

  // Add X-Cart-ID header if required and cartId exists
  if (requiresCartId) {
    const cartId = state.cart.cartId;
    if (cartId) {
      headers['X-Cart-ID'] = String(cartId);
    }
  }

  // Add Authorization header if required
  if (endpointConfig.requiresAuth) {
    const token = state.auth.token;
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  // Add Admin Authorization header if required
  if (endpointConfig.requiresAdminAuth) {
    // @ts-ignore
    const token = state.adminAuth.token;
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  const url = `${API_BASE}${finalPath}`;
  try {
    const response = await axios({ 
      method, 
      url, 
      data, 
      headers,
      params: queryParams 
    });
    return response.data;
  } catch (error: any) {
    if (error.response?.status === 401) {
      if (endpointConfig.requiresAdminAuth) {
        store.dispatch(adminLogout());
      } else {
        store.dispatch(logout());
      }
    }

    const errData = error.response?.data;
    if (error.response?.status === 404 && errData?.error?.code === 'CART_NOT_FOUND') {
      store.dispatch(clearCart());
      toast.error('Cart does not exist');
    }

    throw error;
  }
};