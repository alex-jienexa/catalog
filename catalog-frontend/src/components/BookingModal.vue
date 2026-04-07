<template>
  <div v-if="visible" class="modal-overlay" @click.self="close">
    <div class="modal" @click.stop>
      <div class="modal-header">
        <h2>Оформление интереса к товару</h2>
        <button @click="close" class="close-btn">×</button>
      </div>

      <div v-if="!confirmed">
        <!-- Ввод данных -->
        <div class="modal-body">
          <div v-if="existingCustomer && !overrideCustomer" class="existing-info">
            <p>Вы уже записаны у нас как <strong>{{ existingCustomer.first_name }} {{ existingCustomer.last_name || '' }}</strong>.</p>
            <p>Использовать эти данные?</p>
            <div class="actions">
              <button @click="useExistingCustomer" class="btn-primary">Да, это я</button>
              <button @click="overrideCustomer = true" class="btn-secondary">Нет, ввести другие</button>
            </div>
          </div>

          <form v-else @submit.prevent="submit" class="booking-form">
            <div class="form-group">
              <label>Имя *</label>
              <input v-model="form.first_name" type="text" required />
            </div>
            <div class="form-group">
              <label>Фамилия</label>
              <input v-model="form.last_name" type="text" />
            </div>
            <div class="form-group">
              <label>Телефон *</label>
              <input v-model="form.phone" type="tel" required placeholder="+7XXXXXXXXXX" />
            </div>
            <div class="form-group checkbox-group">
              <label>
                <input type="checkbox" v-model="rememberMe" />
                Запомнить меня
              </label>
            </div>
            <div class="form-actions">
              <button type="button" @click="close" class="btn-secondary">Отмена</button>
              <button type="submit" :disabled="loading" class="btn-primary">
                {{ loading ? 'Отправка...' : 'Подтвердить интерес' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <div v-else class="modal-body success">
        <div class="success-icon">✅</div>
        <h3>Спасибо за интерес к товару!</h3>
        <p>Мы постараемся связаться с вами в ближайшее время.</p>
        <p>Вы также можете связаться с продавцом напрямую:</p>
        <div v-if="contacts.length === 0" class="no-contacts">
            <p>Контакты не добавлены администратором.</p>
            <p>Вы можете связаться с нами позже.</p>
        </div>
        <div v-else class="contacts-list">
            <a v-for="contact in contacts" :key="contact.id" :href="contact.url" target="_blank" class="contact-link">
                {{ contact.platform }}
            </a>
        </div>
        <button @click="close" class="btn-primary">Закрыть</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { customerAPI, reservationAPI, contactAPI } from '@/services/api'

const props = defineProps({
  visible: Boolean,
  productId: Number
})
const emit = defineEmits(['close'])

const form = ref({ first_name: '', last_name: '', phone: '' })
const loading = ref(false)
const rememberMe = ref(true)
const confirmed = ref(false)
const existingCustomer = ref(null)
const overrideCustomer = ref(false)
const contacts = ref([])

// Загрузка контактов для отображения после успеха
const loadContacts = async () => {
  try {
    const res = await contactAPI.getAll(true)
    contacts.value = res.data
  } catch (e) {
    console.error(e)
  }
}

// Проверка запомненного клиента при открытии
const checkStoredCustomer = async () => {
  const stored = localStorage.getItem('customer_data')
  if (stored) {
    try {
      const data = JSON.parse(stored)
      // Проверяем актуальность через API
      const res = await customerAPI.getOrCreate({ 
        first_name: data.first_name, 
        last_name: data.last_name || null, 
        phone: data.phone 
      })
      existingCustomer.value = res.data.customer
    } catch (e) {
      console.error(e)
      localStorage.removeItem('customer_data')
    }
  }
}

const validatePhone = (phone) => {
  // Убираем все нецифровые символы
  const cleaned = phone.replace(/\D/g, '')
  // Российские номера: 11 цифр, начинается с 7 или 8
  if (cleaned.length === 11 && (cleaned[0] === '7' || cleaned[0] === '8')) {
    return true
  }
  // Возможен формат +7...
  const cleanedPlus = phone.replace(/[^\d+]/g, '')
  if (/^\+7\d{10}$/.test(cleanedPlus)) {
    return true
  }
  return false
}

const useExistingCustomer = () => {
  // Подтверждаем бронирование с существующим клиентом
  submitWithCustomer(existingCustomer.value.id)
}

const submitWithCustomer = async (customerId) => {
  loading.value = true
  try {
    await reservationAPI.create({ customer_id: customerId, product_id: props.productId })
    confirmed.value = true
    await loadContacts()
  } catch (error) {
    alert('Ошибка бронирования: ' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

const submit = async () => {
  if (!form.value.first_name || !form.value.phone) {
    alert('Пожалуйста, заполните обязательные поля')
    return
  }
  if (!validatePhone(form.value.phone)) {
    alert('Введите корректный российский номер телефона (например, +79123456789 или 89123456789)')
    return
  }
  loading.value = true
  try {
    // Создаём или получаем клиента
    const customerRes = await customerAPI.getOrCreate({
      first_name: form.value.first_name,
      last_name: form.value.last_name || null,
      phone: form.value.phone
    })
    const customer = customerRes.data.customer
    if (rememberMe.value) {
      localStorage.setItem('customer_data', JSON.stringify({
        first_name: customer.first_name,
        last_name: customer.last_name,
        phone: customer.phone
      }))
    }
    // Создаём бронирование
    await reservationAPI.create({ customer_id: customer.id, product_id: props.productId })
    confirmed.value = true
    await loadContacts()
  } catch (error) {
    alert('Ошибка бронирования: ' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

const close = () => {
  emit('close')
  // Сброс состояния
  confirmed.value = false
  existingCustomer.value = null
  overrideCustomer.value = false
  form.value = { first_name: '', last_name: '', phone: '' }
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    checkStoredCustomer()
  } else {
    close()
  }
})
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
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal {
  background: white;
  border-radius: 20px;
  width: 90%;
  max-width: 500px;
  max-height: 85vh;
  overflow-y: auto;
  box-shadow: 0 20px 35px rgba(0, 0, 0, 0.2);
  animation: fadeSlideUp 0.3s ease;
}

@keyframes fadeSlideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #edf2f7;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #1a202c;
}

.close-btn {
  background: none;
  border: none;
  font-size: 28px;
  cursor: pointer;
  color: #a0aec0;
  transition: color 0.2s;
  line-height: 1;
  padding: 0;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: #f7fafc;
  color: #4a5568;
}

.modal-body {
  padding: 24px;
}

.existing-info {
  background: #ebf8ff;
  border-left: 4px solid #3182ce;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 20px;
}

.existing-info p {
  margin: 0 0 8px 0;
}

.actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.booking-form .form-group {
  margin-bottom: 20px;
}

.booking-form label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #2d3748;
}

.booking-form input {
  width: 100%;
  padding: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.booking-form input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.checkbox-group {
  display: flex;
  align-items: center;
  margin: 20px 0;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  margin-bottom: 0;
}

.checkbox-group input {
  width: auto;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.success {
  text-align: center;
}

.success-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.contacts-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  margin: 20px 0;
}

.contact-link {
  background: #edf2f7;
  padding: 8px 16px;
  border-radius: 40px;
  text-decoration: none;
  color: #4a5568;
  font-weight: 500;
  transition: background 0.2s;
}

.contact-link:hover {
  background: #e2e8f0;
}

.btn-primary, .btn-secondary {
  padding: 10px 20px;
  border-radius: 40px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
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
  background: #e2e8f0;
  color: #4a5568;
}

.btn-secondary:hover {
  background: #cbd5e0;
}
</style>