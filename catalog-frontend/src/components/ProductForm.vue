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
          <label for="image_url">URL изображения</label>
          <input
            id="image_url"
            v-model="form.image_url"
            type="text"
            placeholder="https://example.com/image.jpg"
          />
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
      image_url: '',
      is_active: true
    })

    const editing = computed(() => !!props.product)

    // Сбрасываем форму
    const resetForm = () => {
      form.value = {
        name: '',
        price: 0,
        section_id: '',
        description: '',
        image_url: '',
        is_active: true
      }
    }

    // Заполняем форму данными товара при редактировании
    watch(() => props.product, (product) => {
      if (product) {
        form.value = {
          name: product.name || '',
          price: product.price || 0,
          section_id: product.section_id || '',
          description: product.description || '',
          image_url: product.image_url || '',
          is_active: product.is_active !== undefined ? product.is_active : true
        }
      } else {
        resetForm()
      }
    }, { immediate: true })

    const closeModal = () => {
      resetForm()
      emit('close')
    }

    const handleSubmit = async () => {
      loading.value = true
      console.log('Отправка данных:', form.value)
      
      try {
        const submitData = {
          name: form.value.name.trim(),
          price: parseFloat(form.value.price),
          section_id: parseInt(form.value.section_id),
          is_active: form.value.is_active
        }

        // Добавляем опциональные поля, если они заполнены
        if (form.value.description?.trim()) {
          submitData.description = form.value.description.trim()
        }
        
        if (form.value.image_url?.trim()) {
          submitData.image_url = form.value.image_url.trim()
        }

        console.log('Данные для отправки:', submitData)

        if (props.product) {
          // Редактирование существующего товара
          console.log('Редактирование товара с ID:', props.product.id)
          await admin.updateProduct(props.product.id, submitData)
        } else {
          // Создание нового товара
          console.log('Создание нового товара')
          await admin.createProduct(submitData)
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
      closeModal,
      handleSubmit
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
</style>