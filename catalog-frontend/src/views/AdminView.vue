<template>
  <div class="admin-view">
    <!-- Заголовок -->
    <div class="admin-header">
      <h1>👑 Панель администратора</h1>
      <button @click="logout" class="logout-btn">Выйти</button>
    </div>

    <!-- Вкладки -->
    <div class="admin-tabs">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        @click="switchTab(tab.id)"
        :class="['tab-btn', { active: activeTab === tab.id }]"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Контент вкладок -->
    <div class="admin-content">
      <!-- Вкладка товаров -->
      <div v-if="activeTab === 'products'" class="tab-content">
        <div class="section-header">
          <h2>Управление товарами</h2>
          <button @click="openModal('product')" class="add-btn">
            ➕ Добавить товар
          </button>
        </div>

        <div v-if="loading.products" class="loading">
          <div class="spinner"></div>
          <p>Загрузка товаров...</p>
        </div>

        <div v-else-if="products.length === 0" class="empty-state">
          <p>Товаров пока нет</p>
          <button @click="openModal('product')" class="add-first-btn">
            Добавить первый товар
          </button>
        </div>

        <div v-else class="products-table">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Название</th>
                <th>Цена</th>
                <th>Раздел</th>
                <th>Статус</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="product in products" :key="product.id">
                <td>{{ product.id }}</td>
                <td>
                  <div class="product-info-cell">
                    <img 
                      :src="product.image_url || 'https://via.placeholder.com/40x40?text=Нет+фото'" 
                      :alt="product.name"
                      class="product-thumb"
                      @error="handleImageError"
                    />
                    <span>{{ product.name }}</span>
                  </div>
                </td>
                <td>{{ formatPrice(product.price) }} ₽</td>
                <td>{{ getSectionName(product.section_id) }}</td>
                <td>
                  <span :class="['status-badge', product.is_active ? 'active' : 'inactive']">
                    {{ product.is_active ? 'Активен' : 'Неактивен' }}
                  </span>
                </td>
                <td>
                  <div class="action-buttons">
                    <button @click="editItem('product', product)" class="action-btn edit">
                      ✏️
                    </button>
                    <button @click="deleteItem('product', product.id)" class="action-btn delete">
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Вкладка разделов -->
      <div v-else-if="activeTab === 'sections'" class="tab-content">
        <div class="section-header">
          <h2>Управление разделами</h2>
          <button @click="openModal('section')" class="add-btn">
            ➕ Добавить раздел
          </button>
        </div>

        <div v-if="loading.sections" class="loading">
          <div class="spinner"></div>
          <p>Загрузка разделов...</p>
        </div>

        <div v-else-if="sections.length === 0" class="empty-state">
          <p>Разделов пока нет</p>
          <button @click="openModal('section')" class="add-first-btn">
            Добавить первый раздел
          </button>
        </div>

        <div v-else class="sections-grid">
          <div v-for="section in sections" :key="section.id" class="section-card">
            <div class="section-card-header">
              <h3>{{ section.name }}</h3>
              <span class="product-count">{{ section.product_count }} товаров</span>
            </div>
            <p class="section-description">
              {{ section.description || 'Нет описания' }}
            </p>
            <div class="section-card-footer">
              <span :class="['status-badge', section.is_active ? 'active' : 'inactive']">
                {{ section.is_active ? 'Активен' : 'Неактивен' }}
              </span>
              <div class="action-buttons">
                <button @click="editItem('section', section)" class="action-btn edit">
                  ✏️
                </button>
                <button 
                  @click="deleteItem('section', section.id)" 
                  :disabled="section.product_count > 0"
                  class="action-btn delete"
                  :title="section.product_count > 0 ? 'Нельзя удалить раздел с товарами' : ''"
                >
                  🗑️
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Вкладка контактов -->
      <div v-else-if="activeTab === 'contacts'" class="tab-content">
        <div class="section-header">
          <h2>Управление контактами</h2>
          <button @click="openModal('contact')" class="add-btn">
            ➕ Добавить контакт
          </button>
        </div>

        <div v-if="loading.contacts" class="loading">
          <div class="spinner"></div>
          <p>Загрузка контактов...</p>
        </div>

        <div v-else-if="contacts.length === 0" class="empty-state">
          <p>Контактов пока нет</p>
          <button @click="openModal('contact')" class="add-first-btn">
            Добавить первый контакт
          </button>
        </div>

        <div v-else class="contacts-grid">
          <div v-for="contact in contacts" :key="contact.id" class="contact-card">
            <div class="contact-icon">
              {{ getPlatformIcon(contact.platform) }}
            </div>
            <div class="contact-details">
              <h3>{{ contact.platform }}</h3>
              <a :href="contact.url" target="_blank" class="contact-url">
                {{ contact.url }}
              </a>
              <div class="contact-meta">
                <span class="order">Порядок: {{ contact.order }}</span>
                <span :class="['status-badge', contact.is_active ? 'active' : 'inactive']">
                  {{ contact.is_active ? 'Активен' : 'Неактивен' }}
                </span>
              </div>
            </div>
            <div class="action-buttons">
              <button @click="editItem('contact', contact)" class="action-btn edit">
                ✏️
              </button>
              <button @click="deleteItem('contact', contact.id)" class="action-btn delete">
                🗑️
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Модальные окна -->
    <!-- Товар -->
    <ProductForm
      v-if="showModal === 'product'"
      :show="showModal === 'product'"
      :product="editingItem"
      :sections="sections"
      @close="closeModal"
      @saved="handleSaved"
    />

    <!-- Раздел -->
    <SectionForm
      v-if="showModal === 'section'"
      :show="showModal === 'section'"
      :section="editingItem"
      @close="closeModal"
      @saved="handleSaved"
    />

    <!-- Контакт -->
    <ContactForm
      v-if="showModal === 'contact'"
      :show="showModal === 'contact'"
      :contact="editingItem"
      @close="closeModal"
      @saved="handleSaved"
    />
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { admin, auth, productAPI, sectionAPI, contactAPI } from '@/services/api'
import ProductForm from '@/components/ProductForm.vue'
import SectionForm from '@/components/SectionForm.vue'
import ContactForm from '@/components/ContactForm.vue'

export default {
  name: 'AdminView',
  components: {
    ProductForm,
    SectionForm,
    ContactForm
  },
  
  setup() {
    const router = useRouter()

    // Состояние
    const products = ref([])
    const sections = ref([])
    const contacts = ref([])
    
    const loading = ref({
      products: false,
      sections: false,
      contacts: false
    })

    // Вкладки
    const tabs = [
      { id: 'products', label: 'Товары' },
      { id: 'sections', label: 'Разделы' },
      { id: 'contacts', label: 'Контакты' }
    ]
    
    const activeTab = ref('products')

    // Модальное окно
    const showModal = ref('')
    const editingItem = ref(null)

    const handleProductSaved = async () => {
      console.log('Данные сохранены')
      closeModal()
      await loadProducts()    // дожидаемся загрузки товаров
      await loadSections()    // если нужно
    }
    
    // Загрузка данных
    const loadProducts = async () => {
      loading.value.products = true
      try {
        const response = await productAPI.getAll({ limit: 100 })
        products.value = response.data.products || []
      } catch (error) {
        console.error('Ошибка загрузки товаров:', error)
      } finally {
        loading.value.products = false
      }
    }

    const loadSections = async () => {
      loading.value.sections = true
      try {
        const response = await sectionAPI.getAll(false)
        sections.value = response.data || []
      } catch (error) {
        console.error('Ошибка загрузки разделов:', error)
      } finally {
        loading.value.sections = false
      }
    }

    const loadContacts = async () => {
      loading.value.contacts = true
      try {
        const response = await contactAPI.getAll(false)
        contacts.value = response.data || []
      } catch (error) {
        console.error('Ошибка загрузки контактов:', error)
      } finally {
        loading.value.contacts = false
      }
    }

    // Вспомогательные функции
    const formatPrice = (price) => {
      return new Intl.NumberFormat('ru-RU').format(price)
    }

    const getSectionName = (sectionId) => {
      const section = sections.value.find(s => s.id === sectionId)
      return section ? section.name : 'Неизвестно'
    }

    const getPlatformIcon = (platform) => {
      const icons = {
        'WhatsApp': '💬',
        'Telegram': '✈️',
        'VKontakte': '🌐',
        'Instagram': '📷',
        'Email': '📧',
        'Phone': '📞',
        'Website': '🌍'
      }
      return icons[platform] || '📱'
    }

    const handleImageError = (event) => {
      event.target.src = 'https://via.placeholder.com/40x40?text=Нет+фото'
    }

    // Управление вкладками
    const switchTab = (tabId) => {
      console.log('Переключение на вкладку:', tabId)
      activeTab.value = tabId
    }

    // Управление модальными окнами
    const openModal = (type, item = null) => {
      console.log('Открытие модального окна:', type, item)
      editingItem.value = item
      showModal.value = type
    }

    const closeModal = () => {
      console.log('Закрытие модального окна')
      showModal.value = ''
      editingItem.value = null
    }

    const editItem = (type, item) => {
      console.log('Редактирование:', type, item)
      openModal(type, item)
    }

    const deleteItem = async (type, id) => {
      if (!confirm(`Удалить этот ${type === 'product' ? 'товар' : type === 'section' ? 'раздел' : 'контакт'}?`)) {
        return
      }

      try {
        switch (type) {
          case 'product':
            await admin.deleteProduct(id)
            await loadProducts()
            break
          case 'section':
            await admin.deleteSection(id)
            await loadSections()
            break
          case 'contact':
            await admin.deleteContact(id)
            await loadContacts()
            break
        }
        
        // Обновляем разделы для обновления счетчиков товаров
        if (type === 'product' || type === 'section') {
          await loadSections()
        }
        
        alert('Удаление выполнено успешно')
      } catch (error) {
        console.error('Ошибка при удалении:', error)
        alert(`Ошибка при удалении: ${error.message}`)
      }
    }

    const handleSaved = () => {
      console.log('Данные сохранены')
      closeModal()
      
      // Перезагружаем данные
      loadProducts()
      loadSections()
      loadContacts()
    }

    const logout = () => {
      auth.logout()
      router.push('/')
    }

    // Инициализация
    onMounted(async () => {
      console.log('Инициализация админ-панели')
      
      if (!auth.isAuthenticated()) {
        console.log('Пользователь не аутентифицирован, перенаправление...')
        router.push('/admin/login')
        return
      }

      console.log('Загрузка данных...')
      await loadProducts()
      await loadSections()
      await loadContacts()
      console.log('Данные загружены')
    })

    return {
      // Данные
      products,
      sections,
      contacts,
      loading,
      
      // Вкладки
      tabs,
      activeTab,
      
      // Модальные окна
      showModal,
      editingItem,
      
      // Методы
      formatPrice,
      getSectionName,
      getPlatformIcon,
      handleImageError,
      
      // Управление вкладками
      switchTab,
      handleProductSaved,
      
      // Управление модальными окнами
      openModal,
      closeModal,
      editItem,
      deleteItem,
      handleSaved,
      
      // Выход
      logout
    }
  }
}
</script>

<style scoped>
.admin-view {
  min-height: calc(100vh - 70px);
  background: #f5f7fa;
  padding: 20px;
  margin-top: 70px;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #e9ecef;
}

.admin-header h1 {
  margin: 0;
  color: #333;
  font-size: 28px;
}

.logout-btn {
  background: #ff6b6b;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.logout-btn:hover {
  background: #ff5252;
  transform: translateY(-2px);
}

.admin-tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  background: white;
  border-radius: 10px;
  padding: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.tab-btn {
  flex: 1;
  padding: 12px 20px;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  color: #666;
  transition: all 0.3s;
  text-align: center;
}

.tab-btn:hover {
  background: #f8f9fa;
  color: #333;
}

.tab-btn.active {
  background: #667eea;
  color: white;
  box-shadow: 0 4px 6px rgba(102, 126, 234, 0.2);
}

.admin-content {
  background: white;
  border-radius: 10px;
  padding: 30px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.section-header h2 {
  margin: 0;
  color: #333;
  font-size: 24px;
}

.add-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  color: #666;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  text-align: center;
  color: #666;
}

.add-first-btn {
  margin-top: 20px;
  background: #667eea;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
}

.add-first-btn:hover {
  background: #764ba2;
}

.products-table {
  overflow-x: auto;
}

.products-table table {
  width: 100%;
  border-collapse: collapse;
}

.products-table th {
  background: #f8f9fa;
  padding: 12px 15px;
  text-align: left;
  font-weight: 600;
  color: #333;
  border-bottom: 2px solid #e9ecef;
}

.products-table td {
  padding: 12px 15px;
  border-bottom: 1px solid #eee;
}

.products-table tbody tr:hover {
  background: #f8f9ff;
}

.product-info-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.product-thumb {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
}

.sections-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.section-card {
  background: #f8f9fa;
  border-radius: 10px;
  padding: 20px;
  border: 1px solid #e9ecef;
  transition: all 0.3s;
}

.section-card:hover {
  border-color: #667eea;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
}

.section-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.section-card h3 {
  margin: 0;
  color: #333;
  font-size: 18px;
  flex: 1;
}

.product-count {
  background: #667eea;
  color: white;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: bold;
}

.section-description {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 15px;
}

.section-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.contacts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.contact-card {
  display: flex;
  align-items: center;
  gap: 15px;
  background: #f8f9fa;
  border-radius: 10px;
  padding: 20px;
  border: 1px solid #e9ecef;
  transition: all 0.3s;
}

.contact-card:hover {
  border-color: #667eea;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
}

.contact-icon {
  font-size: 32px;
  width: 60px;
  height: 60px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
}

.contact-details {
  flex: 1;
}

.contact-details h3 {
  margin: 0 0 5px 0;
  color: #333;
  font-size: 16px;
}

.contact-url {
  display: block;
  color: #667eea;
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 10px;
  word-break: break-all;
}

.contact-url:hover {
  text-decoration: underline;
}

.contact-meta {
  display: flex;
  gap: 15px;
  align-items: center;
}

.order {
  font-size: 12px;
  color: #666;
}

.status-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background: #d4edda;
  color: #155724;
}

.status-badge.inactive {
  background: #f8d7da;
  color: #721c24;
}

.action-buttons {
  display: flex;
  gap: 5px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  transition: all 0.3s;
}

.action-btn.edit {
  background: #ffeaa7;
  color: #e17055;
}

.action-btn.edit:hover {
  background: #fab1a0;
  color: white;
  transform: scale(1.1);
}

.action-btn.delete {
  background: #ffcccc;
  color: #ff6b6b;
}

.action-btn.delete:hover:not(:disabled) {
  background: #ff6b6b;
  color: white;
  transform: scale(1.1);
}

.action-btn.delete:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

@media (max-width: 768px) {
  .admin-view {
    padding: 15px;
    margin-top: 60px;
  }
  
  .admin-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .admin-tabs {
    flex-wrap: wrap;
  }
  
  .tab-btn {
    flex: 1 0 calc(33.333% - 10px);
    min-width: 100px;
    padding: 10px 15px;
    font-size: 14px;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .sections-grid,
  .contacts-grid {
    grid-template-columns: 1fr;
  }
  
  .contact-card {
    flex-direction: column;
    text-align: center;
  }
  
  .contact-details {
    text-align: center;
  }
  
  .contact-meta {
    justify-content: center;
  }
}
</style>