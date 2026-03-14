<template>
  <div v-if="show" class="modal-overlay" @click.self="closeModal">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h2>{{ editing ? 'Редактировать раздел' : 'Добавить раздел' }}</h2>
        <button @click="closeModal" class="close-btn">×</button>
      </div>
      
      <form @submit.prevent="handleSubmit" class="modal-form">
        <div class="form-group">
          <label for="name">Название раздела *</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            required
            placeholder="Введите название раздела"
          />
        </div>
        
        <div class="form-group">
          <label for="description">Описание</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            placeholder="Введите описание раздела"
          ></textarea>
        </div>
        
        <div class="form-group checkbox-group">
          <label>
            <input type="checkbox" v-model="form.is_active" />
            Активный раздел
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
  name: 'SectionForm',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    section: {
      type: Object,
      default: null
    }
  },
  emits: ['close', 'saved'],
  
  setup(props, { emit }) {
    const loading = ref(false)
    const form = ref({
      name: '',
      description: '',
      is_active: true
    })

    const editing = computed(() => !!props.section)

    const resetForm = () => {
      form.value = {
        name: '',
        description: '',
        is_active: true
      }
    }

    watch(() => props.section, (section) => {
      if (section) {
        form.value = {
          name: section.name || '',
          description: section.description || '',
          is_active: section.is_active !== undefined ? section.is_active : true
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
      
      try {
        const submitData = {
          name: form.value.name.trim(),
          is_active: form.value.is_active
        }

        if (form.value.description.trim()) {
          submitData.description = form.value.description.trim()
        }

        if (props.section) {
          await admin.updateSection(props.section.id, submitData)
        } else {
          await admin.createSection(submitData)
        }

        emit('saved')
        closeModal()
      } catch (error) {
        console.error('Ошибка при сохранении раздела:', error)
        alert(`Ошибка: ${error.message || 'Не удалось сохранить раздел'}`)
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
/* Стили такие же как в ProductForm */
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

@media (max-width: 480px) {
  .modal {
    width: 95%;
    max-height: 95vh;
  }
  .modal-form {
    padding: 15px;
  }
}
</style>