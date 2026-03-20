<template>
  <div class="store-editor">
    <div class="section-header">
      <h2>Редактирование страницы "О магазине"</h2>
      <button @click="saveStore" class="save-btn" :disabled="loading">
        {{ loading ? 'Сохранение...' : 'Сохранить изменения' }}
      </button>
    </div>

    <div v-if="loadingData" class="loading">
      <div class="spinner"></div>
      <p>Загрузка данных...</p>
    </div>

    <form v-else @submit.prevent="saveStore" class="store-form">
      <!-- Заголовок -->
      <div class="form-group">
        <label for="title">Заголовок *</label>
        <input
          id="title"
          v-model="store.title"
          type="text"
          required
          placeholder="🛒 Добро пожаловать в наш каталог!"
        />
      </div>

      <!-- Описание -->
      <div class="form-group">
        <label for="description">Описание *</label>
        <textarea
          id="description"
          v-model="store.description"
          rows="4"
          required
          placeholder="Мы - современный онлайн-магазин, который предлагает широкий ассортимент товаров..."
        ></textarea>
      </div>

      <!-- URL изображения -->
      <div class="form-group">
        <label for="image_url">URL изображения</label>
        <input
          id="image_url"
          v-model="store.image_url"
          type="text"
          placeholder="https://example.com/store-image.jpg"
        />
        <div v-if="store.image_url" class="image-preview">
          <img :src="store.image_url" alt="Preview" />
        </div>
      </div>

      <!-- Функции магазина -->
      <div class="features-section">
        <h3>Функции магазина</h3>
        <div class="features-list">
          <div v-for="(feature, idx) in store.features" :key="idx" class="feature-card">
            <div class="feature-icon">{{ feature.icon || '📦' }}</div>
            <div class="feature-content">
              <h4>{{ feature.title || 'Без названия' }}</h4>
              <p>{{ feature.description || 'Нет описания' }}</p>
            </div>
            <div class="feature-actions">
              <button @click="editFeature(idx)" class="edit-btn">✏️</button>
              <button @click="removeFeature(idx)" class="delete-btn">🗑️</button>
            </div>
          </div>
          <button @click="addFeature" class="add-feature-btn">
            ➕ Добавить функцию
          </button>
        </div>
      </div>
    </form>

    <!-- Модальное окно для редактирования функции -->
    <div v-if="showFeatureModal" class="modal-overlay" @click.self="closeFeatureModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ editingFeatureIndex === -1 ? 'Добавить функцию' : 'Редактировать функцию' }}</h2>
          <button @click="closeFeatureModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>Иконка (эмодзи) *</label>
            <input v-model="featureForm.icon" type="text" placeholder="⭐" maxlength="2" required />
          </div>
          <div class="form-group">
            <label>Название *</label>
            <input v-model="featureForm.title" type="text" placeholder="Качество товаров" required />
          </div>
          <div class="form-group">
            <label>Описание *</label>
            <textarea v-model="featureForm.description" rows="2" placeholder="Краткое описание функции" required></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeFeatureModal" class="btn-secondary">Отмена</button>
          <button @click="saveFeature" class="btn-primary">Сохранить</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { storeAPI } from '@/services/api'

const store = ref({
  title: '',
  description: '',
  image_url: '',
  features: []
})

const loadingData = ref(false)
const loading = ref(false)

// Для модального окна функций
const showFeatureModal = ref(false)
const editingFeatureIndex = ref(-1)
const featureForm = ref({
  icon: '',
  title: '',
  description: ''
})

// Загрузка данных
const loadStore = async () => {
  loadingData.value = true
  try {
    const response = await storeAPI.get()
    store.value = response.data
    // Если image_url null, делаем пустую строку
    if (!store.value.image_url) store.value.image_url = ''
    // Убедимся, что features всегда массив
    if (!store.value.features) store.value.features = []
  } catch (error) {
    console.error('Ошибка загрузки данных магазина:', error)
    alert('Не удалось загрузить данные магазина')
  } finally {
    loadingData.value = false
  }
}

// Сохранение
const saveStore = async () => {
  loading.value = true
  try {
    // Подготовка данных: image_url может быть пустой строкой – отправляем null
    const dataToSend = {
      title: store.value.title.trim(),
      description: store.value.description.trim(),
      image_url: store.value.image_url?.trim() || null,
      features: store.value.features.filter(f => f.icon && f.title && f.description)
    }
    await storeAPI.update(dataToSend)
    alert('Информация о магазине успешно обновлена')
  } catch (error) {
    console.error('Ошибка сохранения:', error)
    alert('Ошибка при сохранении: ' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

// Функции
const addFeature = () => {
  editingFeatureIndex.value = -1
  featureForm.value = { icon: '', title: '', description: '' }
  showFeatureModal.value = true
}

const editFeature = (idx) => {
  editingFeatureIndex.value = idx
  const f = store.value.features[idx]
  featureForm.value = { ...f }
  showFeatureModal.value = true
}

const removeFeature = (idx) => {
  if (confirm('Удалить эту функцию?')) {
    store.value.features.splice(idx, 1)
  }
}

const saveFeature = () => {
  const newFeature = {
    icon: featureForm.value.icon.trim(),
    title: featureForm.value.title.trim(),
    description: featureForm.value.description.trim()
  }
  if (!newFeature.icon || !newFeature.title || !newFeature.description) {
    alert('Все поля функции должны быть заполнены')
    return
  }
  if (editingFeatureIndex.value === -1) {
    store.value.features.push(newFeature)
  } else {
    store.value.features[editingFeatureIndex.value] = newFeature
  }
  closeFeatureModal()
}

const closeFeatureModal = () => {
  showFeatureModal.value = false
  editingFeatureIndex.value = -1
  featureForm.value = { icon: '', title: '', description: '' }
}

onMounted(() => {
  loadStore()
})
</script>

<style scoped>
.store-editor {
  padding: 20px;
}
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}
.save-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
}
.save-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.store-form {
  max-width: 800px;
}
.form-group {
  margin-bottom: 20px;
}
.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
}
.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.image-preview {
  margin-top: 10px;
}
.image-preview img {
  max-width: 300px;
  max-height: 200px;
  border-radius: 4px;
}
.features-section {
  margin-top: 30px;
}
.features-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.feature-card {
  display: flex;
  align-items: center;
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  gap: 15px;
}
.feature-icon {
  font-size: 32px;
  width: 60px;
  text-align: center;
}
.feature-content {
  flex: 1;
}
.feature-content h4 {
  margin: 0 0 5px 0;
}
.feature-content p {
  margin: 0;
  color: #666;
}
.feature-actions {
  display: flex;
  gap: 5px;
}
.add-feature-btn {
  background: #f0f0f0;
  border: 1px dashed #999;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  margin-top: 10px;
}
.edit-btn, .delete-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
}
.edit-btn {
  background: #ffeaa7;
}
.delete-btn {
  background: #ffcccc;
}
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.modal {
  background: white;
  border-radius: 10px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}
.modal-body {
  padding: 20px;
}
.modal-footer {
  padding: 20px;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.btn-primary {
  background: #667eea;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}
.btn-secondary {
  background: #6c757d;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}
</style>