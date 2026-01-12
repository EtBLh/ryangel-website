

import axios from 'axios';

const API_BASE = 'http://localhost:8080/api';

interface ApiEndpoint {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  path: string;
}

export const api: Record<string, ApiEndpoint> = {
  getProducts: {
    method: 'GET',
    path: '/products'
  },
  getProduct: {
    method: 'GET',
    path: '/products/:productId'
  }
};

export const callAPI = async (endpoint: keyof typeof api, params?: Record<string, string | number>) => {
  const { method, path } = api[endpoint];
  let finalPath = path;
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      finalPath = finalPath.replace(`:${key}`, String(value));
    });
  }
  const url = `${API_BASE}${finalPath}`;
  const response = await axios({ method, url });
  return response.data;
};