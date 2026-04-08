import axios from 'axios';

const API_URL = '/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Автоматически добавляем JWT токен в каждый запрос
api.interceptors.request.use((config) => {
  const token = auth.getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Если токен протух — редиректим на логин
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      auth.logout();
      window.location.href = '/admin/login';
    }
    return Promise.reject(error);
  }
);

// ─── Публичные запросы ────────────────────────────────────────────────────────

export const productAPI = {
  getAll: (params = {}) => api.get('/public/products', { params }),
  getById: (id) => api.get(`/public/products/${id}`),
  getBySection: (sectionId, params = {}) =>
    api.get('/public/products', { params: { ...params, section_id: sectionId } }),
};

export const sectionAPI = {
  getAll: (activeOnly = true) =>
    api.get('/public/sections', { params: { active_only: activeOnly } }),
  getById: (id) => api.get(`/public/sections/${id}`),
};

export const contactAPI = {
  getAll: (activeOnly = true) =>
    api.get('/public/contacts', { params: { active_only: activeOnly } }),
};

export const customerAPI = {
  getOrCreate: (data) => api.post('/public/customers', data),
};

export const reservationAPI = {
  create: (data) => api.post('/public/reservations', data),
  getAll: (params) => api.get('/admin/reservations', { params }),
  updateStatus: (id, data) => api.put(`/admin/reservations/${id}`, data),
};

export const storeAPI = {
  get: () => api.get('/public/store-info'),
  update: (data) => api.put('/admin/store-info', data),
  uploadStoreImage: (formData, onProgress) =>
    api.post('/admin/store/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress?.(percentCompleted);
        }
      },
    }),
};

// ─── Административные запросы ─────────────────────────────────────────────────

export const adminAPI = {
  createProduct: (product) => api.post('/admin/products', product),
  updateProduct: (id, product) => api.put(`/admin/products/${id}`, product),
  deleteProduct: (id) => api.delete(`/admin/products/${id}`),

  createSection: (section) => api.post('/admin/sections', section),
  updateSection: (id, section) => api.put(`/admin/sections/${id}`, section),
  deleteSection: (id) => api.delete(`/admin/sections/${id}`),

  createContact: (contact) => api.post('/admin/contacts', contact),
  updateContact: (id, contact) => api.put(`/admin/contacts/${id}`, contact),
  deleteContact: (id) => api.delete(`/admin/contacts/${id}`),

  putStoreInfo: (data) => api.put('/admin/store-info', data),

  uploadProductImage: (id, formData, onProgress) =>
    api.post(`/admin/products/${id}/image`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress?.(percentCompleted);
        }
      },
    }),

  getAdmins: () => api.get('/admin/admins'),
  createAdmin: (data) => api.post('/admin/admins', data),
  updateAdmin: (id, data) => api.put(`/admin/admins/${id}`, data),
  deleteAdmin: (id) => api.delete(`/admin/admins/${id}`),
};

// Алиас для обратной совместимости с компонентами, которые импортируют `admin`
export const admin = adminAPI;

// ─── Аутентификация ───────────────────────────────────────────────────────────

export const authAPI = {
  isFirst: () => api.get('/auth/is-first'),
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
};

export const auth = {
  login: async (username, password) => {
    const response = await authAPI.login({ username, password });
    const { token, admin: adminData } = response.data;
    localStorage.setItem('token', token);
    localStorage.setItem('adminData', JSON.stringify(adminData));
    return adminData;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('adminData');
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },

  getToken: () => {
    return localStorage.getItem('token');
  },

  getAdminData: () => {
    try {
      return JSON.parse(localStorage.getItem('adminData')) || null;
    } catch {
      return null;
    }
  },
};

export default api;
