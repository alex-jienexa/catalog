<template>
  <div class="product-detail">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Загрузка товара...</p>
    </div>
    
    <div v-else-if="product" class="product-container">
      <button @click="$router.back()" class="back-btn">← Назад</button>
      
      <div class="product-grid">
        <div class="product-images">
          <img :src="mainImageSrc" :alt="product.name" class="main-image" @error="handleImageError" />
        </div>
        
        <div class="product-info">
          <div class="section-badge">
            {{ getSectionName(product.section_id) }}
          </div>
          
          <h1 class="product-title">{{ product.name }}</h1>
          
          <div class="product-price">
            {{ formatPrice(product.price) }} ₽
          </div>
          
          <div class="product-meta">
            <span class="meta-item">
              <span class="meta-label">Добавлен:</span>
              {{ formatDate(product.created_at) }}
            </span>
            <span class="meta-item">
              <span class="meta-label">Статус:</span>
              <span :class="['status', product.is_active ? 'active' : 'inactive']">
                {{ product.is_active ? 'В наличии' : 'Нет в наличии' }}
              </span>
            </span>
          </div>
          
          <div class="product-description">
            <h3>Описание</h3>
            <p v-if="product.description">{{ product.description }}</p>
            <p v-else class="no-description">Описание отсутствует</p>
          </div>
          
          <button @click.stop="showBookingModal = true" class="contact-seller-btn">
            📦 Хочу забрать товар
          </button>
          <BookingModal 
            :visible="showBookingModal" 
            :product-id="product.id" 
            @close="showBookingModal = false"
          />
          
          <!-- <div class="seller-info">
            <h3>Информация о связи</h3>
            <p>Для связи с продавцом используйте <a href="/about">контакты</a>.</p>
          </div> -->
        </div>
      </div>
      
      <!-- Похожие товары -->
      <div v-if="similarProducts.length > 0" class="similar-products">
        <h2>Похожие товары</h2>
        <div class="similar-grid">
          <div
            v-for="similar in similarProducts"
            :key="similar.id"
            class="similar-card"
            @click="goToProduct(similar.id)"
          >
            <img :src="similar.image_url && similar.image_url.startsWith('http') ? similar.image_url : similar.image_url || 'https://via.placeholder.com/150x100?text=Нет+фото'" />
            <div class="similar-info">
              <h4>{{ similar.name }}</h4>
              <div class="similar-price">{{ formatPrice(similar.price) }} ₽</div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="not-found">
      <h2>Товар не найден</h2>
      <p>К сожалению, запрашиваемый товар не существует или был удален.</p>
      <router-link to="/" class="home-link">Вернуться в каталог</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { productAPI, sectionAPI } from '@/services/api'
import BookingModal from '@/components/BookingModal.vue'

const route = useRoute()
const product = ref(null)
const sections = ref([])
const similarProducts = ref([])
const loading = ref(true)
const showContactModal = ref(false)
const showBookingModal = ref(false)

const mainImageSrc = computed(() => {
  if (!product.value?.image_url) {
    return 'https://placehold.co/500x400?text=Нет+фото'
  }
  if (product.value.image_url.startsWith('http')) {
    return product.value.image_url
  }
  return product.value.image_url
})

const loadProduct = async () => {
  loading.value = true
  try {
    const id = parseInt(route.params.id)
    
    // Загружаем товар
    const response = await productAPI.getById(id)
    product.value = response.data
    
    // Загружаем разделы
    const sectionsResponse = await sectionAPI.getAll()
    sections.value = sectionsResponse.data
    
    // Загружаем похожие товары
    await loadSimilarProducts(id, response.data.section_id)
    
  } catch (error) {
    console.error('Ошибка загрузки товара:', error)
    product.value = null
  } finally {
    loading.value = false
  }
}

const loadSimilarProducts = async (currentId, sectionId) => {
  try {
    const response = await productAPI.getBySection(sectionId, { 
      limit: 4,
      is_active: true 
    })
    
    // Фильтруем текущий товар и берем только 3 похожих
    similarProducts.value = response.data.products
      .filter(p => p.id !== currentId)
      .slice(0, 3)
  } catch (error) {
    console.error('Ошибка загрузки похожих товаров:', error)
    similarProducts.value = []
  }
}

const getSectionName = (sectionId) => {
  const section = sections.value.find(s => s.id === sectionId)
  return section ? section.name : 'Неизвестный раздел'
}

const formatPrice = (price) => {
  return new Intl.NumberFormat('ru-RU').format(price)
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const handleImageError = (event) => {
  event.target.src = 'https://placehold.co/500x400?text=Нет+фото'
}

const goToProduct = (id) => {
  // Используем window.location для полной перезагрузки страницы
  window.location.href = `/product/${id}`
}

onMounted(() => {
  loadProduct()
})
</script>

<style scoped>
.product-detail {
  padding: 40px 20px;
  max-width: 1200px;
  margin: 70px auto 0;
  min-height: calc(100vh - 70px);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  color: #666;
}

.spinner {
  width: 50px;
  height: 50px;
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

.back-btn {
  background: none;
  border: none;
  color: #667eea;
  font-size: 16px;
  cursor: pointer;
  padding: 10px 0;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 5px;
  transition: color 0.3s;
}

.back-btn:hover {
  color: #764ba2;
}

.product-container {
  background: white;
  border-radius: 15px;
  padding: 30px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.product-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
  margin-bottom: 40px;
}

.product-images {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.main-image {
  width: 100%;
  height: 400px;
  object-fit: cover;
  border-radius: 10px;
}

.product-info {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-badge {
  display: inline-block;
  background: #f1f3ff;
  color: #667eea;
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: 500;
  font-size: 14px;
  align-self: flex-start;
}

.product-title {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  margin: 0;
  line-height: 1.3;
}

.product-price {
  font-size: 36px;
  font-weight: 700;
  color: #667eea;
}

.product-meta {
  display: flex;
  gap: 30px;
  padding: 15px 0;
  border-top: 1px solid #eee;
  border-bottom: 1px solid #eee;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.meta-label {
  font-size: 12px;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status {
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 14px;
}

.status.active {
  background: #d4edda;
  color: #155724;
}

.status.inactive {
  background: #f8d7da;
  color: #721c24;
}

.product-description h3 {
  margin-bottom: 15px;
  color: #333;
  font-size: 20px;
}

.product-description p {
  line-height: 1.6;
  color: #555;
  font-size: 16px;
}

.no-description {
  color: #888;
  font-style: italic;
}

.contact-seller-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 16px 32px;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
  margin-top: 10px;
}

.contact-seller-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
}

.seller-info {
  margin-top: 20px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.seller-info h3 {
  margin-bottom: 10px;
  color: #333;
}

.seller-info p {
  color: #666;
  line-height: 1.5;
}

.similar-products {
  margin-top: 60px;
}

.similar-products h2 {
  margin-bottom: 30px;
  color: #333;
  font-size: 24px;
  text-align: center;
}

.similar-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.similar-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s;
}

.similar-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.similar-image {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.similar-info {
  padding: 15px;
}

.similar-info h4 {
  margin: 0 0 10px 0;
  font-size: 16px;
  color: #333;
  line-height: 1.3;
}

.similar-price {
  font-size: 18px;
  font-weight: 600;
  color: #667eea;
}

.not-found {
  text-align: center;
  padding: 100px 20px;
}

.not-found h2 {
  color: #333;
  margin-bottom: 20px;
}

.not-found p {
  color: #666;
  margin-bottom: 30px;
  font-size: 18px;
}

.home-link {
  display: inline-block;
  background: #667eea;
  color: white;
  padding: 12px 24px;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s;
}

.home-link:hover {
  background: #764ba2;
  transform: translateY(-2px);
}

@media (max-width: 992px) {
  .product-grid {
    grid-template-columns: 1fr;
    gap: 30px;
  }
  
  .main-image {
    height: 300px;
  }
}

@media (max-width: 576px) {
  .product-detail {
    padding: 20px 15px;
    margin-top: 60px;
  }
  
  .product-container {
    padding: 20px;
  }
  
  .product-title {
    font-size: 24px;
  }
  
  .product-price {
    font-size: 28px;
  }
}
</style>