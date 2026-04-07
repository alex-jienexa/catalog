<template>
  <div class="product-card" @click="goToProduct">
    <div class="product-image">
      <img 
        :src="image_src" 
        :alt="product.name"
        @error="handleImageError"
      />
    </div>
    <div class="product-info">
      <h3 class="product-name">{{ product.name }}</h3>
      <div class="product-price">
        {{ formatPrice(product.price) }} ₽
      </div>
      <div class="product-description">
        {{ truncateDescription(product.description) }}
      </div>
      <div class="product-meta">
        <span class="section-badge">
          {{ getSectionName(product.section_id) }}
        </span>
        <span class="date">
          {{ formatDate(product.created_at) }}
        </span>
      </div>
      <button @click="$emit('book')" class="contact-btn">📦 Хочу забрать товар</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  product: {
    type: Object,
    required: true,
    default: () => ({})
  },
  sections: {
    type: Array,
    default: () => []
  }
})

const router = useRouter()

const image_src = computed(() => {
  if (!props.product.image_url) {
    return 'https://placehold.co/600x400?text=Нет+фото'
  }
  // Если путь уже абсолютный (начинается с http), используем как есть
  if (props.product.image_url.startsWith('http')) {
    return props.product.image_url
  }
  // Если путь относительный (начинается с /uploads), добавляем текущий origin
  // В production это будет тот же домен
  return props.product.image_url
})

const getSectionName = (sectionId) => {
  const section = props.sections.find(s => s.id === sectionId)
  return section ? section.name : 'Неизвестный раздел'
}

const formatPrice = (price) => {
  if (!price && price !== 0) return '0'
  return new Intl.NumberFormat('ru-RU').format(price)
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  try {
    return new Date(dateString).toLocaleDateString('ru-RU')
  } catch (e) {
    return ''
  }
}

const truncateDescription = (description) => {
  if (!description) return 'Нет описания'
  return description.length > 100 
    ? description.substring(0, 100) + '...' 
    : description
}

const handleImageError = (event) => {
  console.error('Ошибка загрузки изображения:', props.product.image_url)
  event.target.src = 'https://placehold.co/600x400?text=Нет+фото'
}

const goToProduct = () => {
  if (props.product && props.product.id) {
    router.push(`/product/${props.product.id}`)
  }
}

console.log('ProductCard получил товар:', props.product)
</script>

<style scoped>
.product-card {
  background: white;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15);
}

.product-image {
  height: 200px;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.product-card:hover .product-image img {
  transform: scale(1.05);
}

.product-info {
  padding: 20px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.product-name {
  margin: 0 0 10px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
  line-height: 1.3;
}

.product-price {
  font-size: 22px;
  font-weight: bold;
  color: #667eea;
  margin-bottom: 10px;
}

.product-description {
  color: #666;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 15px;
  flex-grow: 1;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 12px;
  color: #888;
}

.section-badge {
  background: #f1f3ff;
  color: #667eea;
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 500;
}

.date {
  font-style: italic;
}

.contact-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
  margin-top: auto;
}

.contact-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}
</style>