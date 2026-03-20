import axios from 'axios';

const API_URL = '/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Публичные запросы
export const productAPI = {
  // Получить все товары
  getAll: (params = {}) => api.get('/public/products', { params }),
  
  // Получить товар по ID
  getById: (id) => api.get(`/public/products/${id}`),
  
  // Получить товары раздела
  getBySection: (sectionId, params = {}) => 
    api.get('/public/products', { params: { ...params, section_id: sectionId } }),
};

export const sectionAPI = {
  // Получить все разделы
  getAll: (activeOnly = true) => 
    api.get('/public/sections', { params: { active_only: activeOnly } }),
  
  // Получить раздел по ID
  getById: (id) => api.get(`/public/sections/${id}`),
};

export const contactAPI = {
  // Получить все контакты
  getAll: (activeOnly = true) => 
    api.get('/public/contacts', { params: { active_only: activeOnly } }),
};

export const storeAPI = {
    get: () => api.get('/public/store-info'),
    update: (data) => adminAPI.putStoreInfo(data),
};

// Административные запросы
export const adminAPI = {
  // Товары
  createProduct: (product) => api.post('/admin/products', product),
  updateProduct: (id, product) => api.put(`/admin/products/${id}`, product),
  deleteProduct: (id) => api.delete(`/admin/products/${id}`),
  
  // Разделы
  createSection: (section) => api.post('/admin/sections', section),
  updateSection: (id, section) => api.put(`/admin/sections/${id}`, section),
  deleteSection: (id) => api.delete(`/admin/sections/${id}`),
  
  // Контакты
  createContact: (contact) => api.post('/admin/contacts', contact),
  updateContact: (id, contact) => api.put(`/admin/contacts/${id}`, contact),
  deleteContact: (id) => api.delete(`/admin/contacts/${id}`),

  putStoreInfo: (data) => api.put('/admin/store-info', data),
};

// Аутентификация (простая - только для админки)
export const auth = {
  login: (username, password) => {
    // Простая проверка в коде
    if (username === 'admin' && password === 'admin') {
      localStorage.setItem('isAuthenticated', 'true');
      return Promise.resolve({ success: true });
    }
    return Promise.reject(new Error('Неверный логин или пароль'));
  },
  
  logout: () => {
    localStorage.removeItem('isAuthenticated');
  },
  
  isAuthenticated: () => {
    return localStorage.getItem('isAuthenticated') === 'true';
  },
};

// Экспорт с проверкой аутентификации для админки
export const admin = {
  createProduct: async (product) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.createProduct(product);
  },
  
  updateProduct: async (id, product) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.updateProduct(id, product);
  },
  
  deleteProduct: async (id) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.deleteProduct(id);
  },
  
  // Аналогично для разделов и контактов
  createSection: async (section) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.createSection(section);
  },
  
  updateSection: async (id, section) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.updateSection(id, section);
  },
  
  deleteSection: async (id) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.deleteSection(id);
  },
  
  createContact: async (contact) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.createContact(contact);
  },
  
  updateContact: async (id, contact) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.updateContact(id, contact);
  },
  
  deleteContact: async (id) => {
    if (!auth.isAuthenticated()) {
      throw new Error('Требуется авторизация');
    }
    return adminAPI.deleteContact(id);
  },



  uploadProductImage: (id, formData, onProgress) => {
    return api.post(`/admin/products/${id}/image`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress?.(percentCompleted);
        }
      }
    });
}
};

export default api;