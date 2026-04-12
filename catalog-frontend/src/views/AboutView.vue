<template>
  <div class="about">
    <div class="container">
      <div v-if="loading_storeInfo" class="loading">
        <div class="spinner"></div>
        <p>Загрузка...</p>
      </div>
      <div v-else>
        <!-- Информация о магазине -->
        <section class="store-info">
          <h1 class="title">О нашем магазине</h1>
          
          <div class="info-card">
            <div class="info-content">
              <h2>{{ storeInfo.title }}</h2>
              <p>
                {{ storeInfo.description }}
              </p>
              
              <div v-if="storeInfo.features && storeInfo.features.length" class="features">
                <div v-for="feature in storeInfo.features" :key="feature.id" class="feature">
                  <div class="feature-icon">{{ feature.icon }}</div>
                  <div class="feature-text">
                    <h3>{{ feature.title }}</h3>
                    <p>{{ feature.description }}</p>
                  </div>
                </div>
              </div>
            </div>
            
            <div v-if="storeInfo.image_url" class="info-image">
              <img :src="storeInfo.image_url" 
                  alt="Магазин" />
            </div>
          </div>
        </section>
      </div>
      
      
      <!-- Контакты -->
      <section class="contacts-section">
        <h2 class="section-title">📞 Свяжитесь с нами</h2>
        <p class="section-subtitle">Выберите удобный способ связи</p>
        
        <div v-if="loading_contacts" class="loading">
          <div class="spinner"></div>
          <p>Загрузка контактов...</p>
        </div>
        
        <div v-else-if="contacts.length === 0" class="no-contacts">
          <p>Контакты пока не добавлены</p>
        </div>
        
        <div v-else class="contacts-grid">
          <a
            v-for="contact in contacts"
            :key="contact.id"
            :href="contact.url"
            target="_blank"
            :class="['contact-card', { inactive: !contact.is_active }]"
          >
            <div class="contact-icon">
              {{ getPlatformIcon(contact.platform) }}
            </div>
            <div class="contact-info">
              <h3>{{ contact.platform }}</h3>
              <p class="contact-url">{{ getContactDisplay(contact.url) }}</p>
              <span v-if="!contact.is_active" class="contact-status">Неактивно</span>
            </div>
            <div class="contact-arrow">→</div>
          </a>
        </div>
        
        <!-- <div class="working-hours">
          <h3>⏰ Часы работы поддержки:</h3>
          <ul>
            <li><strong>Пн-Пт:</strong> 9:00 - 20:00</li>
            <li><strong>Сб:</strong> 10:00 - 18:00</li>
            <li><strong>Вс:</strong> 12:00 - 16:00</li>
          </ul>
        </div> -->
      </section>
      
      <!-- Карта -->
      <!-- <section class="map-section">
        <h2 class="section-title">📍 Мы находимся</h2>
        <div class="map-placeholder">
          <div class="map-content">
            <p>Москва, ул. Примерная, д. 1</p>
            <p>БЦ "Деловой", 5 этаж</p>
          </div>
        </div>
      </section> -->
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { contactAPI, storeAPI } from '@/services/api'

const storeInfo = ref({
  title: '',
  description: '',
  image_url: '',
  features: []
})
const contacts = ref([])
const loading_contacts = ref(false)
const loading_storeInfo = ref(false)

const loadStoreInfo = async () => {
  loading_storeInfo.value = true
  try {
    const res = await storeAPI.get()
    storeInfo.value = res.data
  } catch (err) {
    console.error("Ошибка загрузки данных о магазине:", err)
  } finally {
    loading_storeInfo.value = false
  }
}

const loadContacts = async () => {
  loading_contacts.value = true
  try {
    const response = await contactAPI.getAll(true)
    if (response.data != null) {
        contacts.value = response.data
    }
  } catch (error) {
    console.error('Ошибка загрузки контактов:', error)
    contacts.value = []
  } finally {
    loading_contacts.value = false
  }
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

const getContactDisplay = (url) => {
  // Упрощаем URL для отображения
  return url
    .replace(/^https?:\/\//, '')
    .replace(/^mailto:/, '')
    .replace(/^tel:/, '')
    .replace(/\/$/, '')
}

onMounted(() => {
  loadContacts()
  loadStoreInfo()
})
</script>

<style scoped>
.about {
  padding: 70px 20px 40px;
  background: #f5f7fa;
  min-height: calc(100vh - 70px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.title {
  text-align: center;
  font-size: 36px;
  color: #333;
  margin-bottom: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.info-card {
  background: white;
  border-radius: 20px;
  padding: 40px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  margin-bottom: 60px;
}

.info-content h2 {
  font-size: 28px;
  color: #333;
  margin-bottom: 20px;
}

.info-content p {
  font-size: 18px;
  line-height: 1.6;
  color: #555;
  margin-bottom: 30px;
}

.features {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-top: 30px;
}

.feature {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 10px;
  transition: all 0.3s;
}

.feature:hover {
  background: #f1f3ff;
  transform: translateX(5px);
}

.feature-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.feature-text h3 {
  font-size: 16px;
  color: #333;
  margin: 0 0 5px 0;
}

.feature-text p {
  font-size: 14px;
  color: #666;
  margin: 0;
  line-height: 1.4;
}

.info-image {
  display: flex;
  align-items: center;
  justify-content: center;
}

.info-image img {
  width: 100%;
  height: 400px;
  object-fit: cover;
  border-radius: 10px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

/* Контакты */
.contacts-section {
  margin-bottom: 60px;
}

.section-title {
  font-size: 32px;
  color: #333;
  text-align: center;
  margin-bottom: 10px;
}

.section-subtitle {
  text-align: center;
  color: #666;
  font-size: 18px;
  margin-bottom: 40px;
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

.no-contacts {
  text-align: center;
  padding: 40px;
  background: white;
  border-radius: 10px;
  color: #666;
}

.contacts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.contact-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  text-decoration: none;
  color: inherit;
  border: 2px solid transparent;
  transition: all 0.3s;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
}

.contact-card:hover:not(.inactive) {
  border-color: #667eea;
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.2);
}

.contact-card.inactive {
  opacity: 0.5;
  cursor: not-allowed;
}

.contact-icon {
  font-size: 32px;
  width: 60px;
  height: 60px;
  background: #f1f3ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.contact-info {
  flex: 1;
}

.contact-info h3 {
  margin: 0 0 5px 0;
  color: #333;
  font-size: 18px;
}

.contact-url {
  color: #666;
  font-size: 14px;
  margin: 0;
  word-break: break-all;
}

.contact-status {
  display: inline-block;
  font-size: 12px;
  color: #ff6b6b;
  background: #ffeaea;
  padding: 2px 8px;
  border-radius: 10px;
  margin-top: 5px;
}

.contact-arrow {
  color: #667eea;
  font-size: 20px;
  opacity: 0;
  transform: translateX(-10px);
  transition: all 0.3s;
}

.contact-card:hover:not(.inactive) .contact-arrow {
  opacity: 1;
  transform: translateX(0);
}

/* Часы работы */
.working-hours {
  background: white;
  border-radius: 12px;
  padding: 25px;
  margin-bottom: 40px;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
}

.working-hours h3 {
  color: #333;
  margin: 0 0 15px 0;
  font-size: 20px;
}

.working-hours ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.working-hours li {
  padding: 8px 0;
  border-bottom: 1px solid #eee;
  color: #555;
  font-size: 16px;
}

.working-hours li:last-child {
  border-bottom: none;
}

.working-hours li strong {
  display: inline-block;
  width: 100px;
  color: #333;
}

/* Карта */
.map-section {
  margin-bottom: 40px;
}

.map-placeholder {
  background: white;
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
  background-image: 
    linear-gradient(45deg, #f5f5f5 25%, transparent 25%),
    linear-gradient(-45deg, #f5f5f5 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #f5f5f5 75%),
    linear-gradient(-45deg, transparent 75%, #f5f5f5 75%);
  background-size: 20px 20px;
  background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
}

.map-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  display: inline-block;
}

.map-content p {
  margin: 10px 0;
  color: #333;
  font-size: 18px;
}

@media (max-width: 768px) {
  .info-card {
    grid-template-columns: 1fr;
    padding: 20px;
    gap: 30px;
  }
  
  .features {
    grid-template-columns: 1fr;
  }
  
  .info-image img {
    height: 250px;
  }
  
  .title {
    font-size: 28px;
  }
  
  .info-content h2 {
    font-size: 24px;
  }
  
  .contacts-grid {
    grid-template-columns: 1fr;
  }
}
</style>