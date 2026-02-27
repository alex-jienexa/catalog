<template>
  <div v-if="show" class="modal-overlay" @click.self="closeModal">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h2>{{ editing ? 'Редактировать товар' : 'Добавить товар' }}</h2>
        <button @click="closeModal" class="close-btn">×</button>
      </div>
      
      <form @submit.prevent="handleSubmit" class="modal-form">
        <div class="form-group">
          <label for="name">Название товара *</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            placeholder="Введите название товара"
          />
        </div>
        
        <div class="form-group">
          <label for="price">Цена (₽) *</label>
          <input
            id="price"
            v-model.number="form.price"
            type="number"
            step="0.01"
            min="0"
            required
            placeholder="Введите цену"
          />
        </div>
        
        <div class="form-group">
          <label for="section_id">Раздел *</label>
          <select id="section_id" v-model.number="form.section_id" required>
            <option value="">Выберите раздел</option>
            <option v-for="section in sections" :key="section.id" :value="section.id">
              {{ section.name }}
            </option>
          </select>
        </div>
        
        <div class="form-group">
          <label for="description">Описание</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            placeholder="Введите описание товара"
          ></textarea>
        </div>
        
        <div class="form-group">
          <label>Изображение</label>
          <input 
            type="file" 
            @change="handleFileSelect" 
            accept="image/jpeg,image/png,image/gif,image/webp"
            ref="fileInput"
          />
          <div v-if="uploadProgress > 0" class="progress">
            <div class="progress-bar" :style="{ width: uploadProgress + '%' }"></div>
            <span class="progress-text">{{ uploadProgress }}%</span>
          </div>
          <div v-if="imagePreview" class="image-preview">
            <img :src="imagePreview" alt="Preview" />
            <button type="button" @click="removeImage" class="remove-image">×</button>
          </div>
          <p v-if="uploadError" class="error">{{ uploadError }}</p>
        </div>
        
        <div class="form-group checkbox-group">
          <label>
            <input type="checkbox" v-model="form.is_active" />
            Активный товар
          </label>
        </div>
        
        <div class="form-actions">
          <button type="button" @click="closeModal" class="btn btn-secondary">
            Отмена
          </button>
          <button type="submit" :disabled="loading" class="btn btn-primary">
            {{ loading ? 'Сохранение...' : (editing ? 'Сохранить' : 'Создать') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, watch, computed } from 'vue'
import { admin } from '@/services/api'

export default {
  name: 'ProductForm',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    product: {
      type: Object,
      default: null
    },
    sections: {
      type: Array,
      default: () => []
    }
  },
  emits: ['close', 'saved'],
  
  setup(props, { emit }) {
    const loading = ref(false)
    const form = ref({
      name: '',
      price: 0,
      section_id: '',
      description: '',
      is_active: true
    })

    // Для изображения
    const selectedFile = ref(null)
    const imagePreview = ref('')
    const uploadProgress = ref(0)
    const uploadError = ref('')
    const fileInput = ref(null)

    const editing = computed(() => !!props.product)

    // Сбрасываем форму
    const resetForm = () => {
      form.value = {
        name: '',
        price: 0,
        section_id: '',
        description: '',
        is_active: true
      }
      selectedFile.value = null
      imagePreview.value = ''
      uploadProgress.value = 0
      uploadError.value = ''
      if (fileInput.value) fileInput.value.value = ''
    }

    // Заполняем форму данными товара при редактировании
    watch(() => props.product, (product) => {
      if (product) {
        form.value = {
          name: product.name || '',
          price: product.price || 0,
          section_id: product.section_id || '',
          description: product.description || '',
          is_active: product.is_active !== undefined ? product.is_active : true
        }
        imagePreview.value = product.image_url || ''
        selectedFile.value = null
        uploadProgress.value = 0
        uploadError.value = ''
      } else {
        resetForm()
      }
    }, { immediate: true })

    const closeModal = () => {
      resetForm()
      emit('close')
    }

    // Обработка выбора файла
    const handleFileSelect = (event) => {
      const file = event.target.files[0]
      if (!file) return
      
      // Проверка типа
      const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
      if (!allowedTypes.includes(file.type)) {
        uploadError.value = 'Допустимые форматы: JPEG, PNG, GIF, WEBP'
        fileInput.value.value = ''
        return
      }
      
      // Проверка размера (например, 5MB)
      const maxSize = 5 * 1024 * 1024
      if (file.size > maxSize) {
        uploadError.value = 'Файл слишком большой. Максимум 5MB'
        fileInput.value.value = ''
        return
      }
      
      uploadError.value = ''
      selectedFile.value = file
      
      // Создание превью
      const reader = new FileReader()
      reader.onload = (e) => {
        imagePreview.value = e.target.result
      }
      reader.readAsDataURL(file)
    }

    // Удаление выбранного изображения
    const removeImage = () => {
      selectedFile.value = null
      imagePreview.value = ''
      uploadProgress.value = 0
      if (fileInput.value) fileInput.value.value = ''
    }

    // Загрузка изображения на сервер (отдельная функция)
    const uploadImage = async (productId) => {
      if (!selectedFile.value) return
      
      const formData = new FormData()
      formData.append('image', selectedFile.value)
      
      try {
        uploadProgress.value = 0
        await admin.uploadProductImage(productId, formData, (progress) => {
          uploadProgress.value = progress
        })
        // После успешной загрузки можно обновить превью (ответ приходит с новым URL)
        // Но проще просто перезагрузить список товаров через emit
      } catch (error) {
        console.error('Ошибка загрузки изображения:', error)
        uploadError.value = 'Не удалось загрузить изображение: ' + (error.response?.data?.error || error.message)
        throw error // пробрасываем, чтобы handleSubmit знал об ошибке
      }
    }

    const handleSubmit = async () => {
      loading.value = true
      uploadError.value = ''
      console.log('Отправка данных:', form.value)
      
      try {
        const productData = {
          name: form.value.name.trim(),
          price: parseFloat(form.value.price),
          section_id: parseInt(form.value.section_id),
          description: form.value.description.trim() || null,
          is_active: form.value.is_active
        }

        let savedProduct
    
        if (editing.value) {
          // Обновление товара (без изображения)
          const response = await admin.updateProduct(props.product.id, productData)
          savedProduct = response.data
        } else {
          // Создание товара (без изображения)
          const response = await admin.createProduct(productData)
          savedProduct = response.data
        }
    
        // Если есть выбранный файл, загружаем его для этого товара
        if (selectedFile.value) {
          await uploadImage(savedProduct.id)
        }

        console.log('Товар успешно сохранен')
        emit('saved')
        closeModal()
      } catch (error) {
        console.error('Ошибка при сохранении товара:', error)
        alert(`Ошибка: ${error.message || 'Не удалось сохранить товар'}`)
      } finally {
        loading.value = false
      }
    }

    return {
      loading,
      form,
      editing,
      imagePreview,
      uploadError,
      uploadProgress,
      closeModal,
      handleSubmit,
      handleFileSelect,
      removeImage,
    }
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
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
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h2 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.3s;
}

.close-btn:hover {
  background: #f5f5f5;
}

.modal-form {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #555;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 5px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.checkbox-group input[type="checkbox"] {
  width: auto;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #5a6268;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.progress {
  margin-top: 10px;
  background: #f0f0f0;
  border-radius: 4px;
  height: 20px;
  position: relative;
  overflow: hidden;
}
.progress-bar {
  background: #667eea;
  height: 100%;
  transition: width 0.3s;
}
.progress-text {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  color: #333;
  font-size: 12px;
  line-height: 20px;
}

.image-preview {
  margin-top: 10px;
  position: relative;
  display: inline-block;
}
.image-preview img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 4px;
  border: 1px solid #ddd;
}

.remove-image {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #ff6b6b;
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.error {
  color: #ff6b6b;
  font-size: 12px;
  margin-top: 5px;
}
</style>