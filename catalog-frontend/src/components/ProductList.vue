<template>
  <div class="product-list">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Загрузка товаров...</p>
    </div>
    
    <div v-else-if="displayProducts.length === 0" class="empty-state">
      <div class="empty-icon">📦</div>
      <h3>Товаров не найдено</h3>
      <p v-if="searchQuery">Попробуйте изменить поисковый запрос</p>
      <p v-else>В этом разделе пока нет товаров</p>
    </div>
    
    <div v-else class="products-grid">
      <ProductCard
        v-for="product in displayProducts"
        :key="product.id"
        :product="product"
        :sections="sections"
      />
    </div>
    
    <!-- Пагинация -->
    <div v-if="totalPages > 1" class="pagination">
      <button 
        @click="changePage(currentPage - 1)"
        :disabled="currentPage === 1"
        class="page-btn"
      >
        ←
      </button>
      
      <span class="page-info">
        Страница {{ currentPage }} из {{ totalPages }}
      </span>
      
      <button 
        @click="changePage(currentPage + 1)"
        :disabled="currentPage === totalPages"
        class="page-btn"
      >
        →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted } from 'vue'
import ProductCard from './ProductCard.vue'
import { productAPI } from '@/services/api'

const props = defineProps({
  sectionId: {
    type: Number,
    default: null
  },
  sortBy: {
    type: String,
    default: 'created_at'
  },
  searchQuery: {
    type: String,
    default: ''
  },
  sections: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['page-change'])

const products = ref([])
const loading = ref(false)
const currentPage = ref(1)
const totalItems = ref(0)

// Параметры запроса
const limit = 12

// Компьютированные свойства для фильтрации и сортировки
const filteredProducts = computed(() => {
  let filtered = [...products.value]
  
  // Фильтрация по поисковому запросу
  if (props.searchQuery) {
    const query = props.searchQuery.toLowerCase().trim()
    filtered = filtered.filter(product => {
      const nameMatch = product.name.toLowerCase().includes(query)
      const descMatch = product.description && 
                       product.description.toLowerCase().includes(query)
      return nameMatch || descMatch
    })
  }
  
  return filtered
})

const sortedProducts = computed(() => {
  const productsToSort = [...filteredProducts.value]
  
  switch (props.sortBy) {
    case 'name':
      return productsToSort.sort((a, b) => a.name.localeCompare(b.name))
    case 'price':
      return productsToSort.sort((a, b) => a.price - b.price)
    case '-price':
      return productsToSort.sort((a, b) => b.price - a.price)
    case 'created_at':
    default:
      return productsToSort.sort((a, b) => 
        new Date(b.created_at) - new Date(a.created_at)
      )
  }
})

const displayProducts = computed(() => {
  const start = (currentPage.value - 1) * limit
  const end = start + limit
  return sortedProducts.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(sortedProducts.value.length / limit)
})

const loadProducts = async () => {
  loading.value = true
  console.log('Загрузка товаров...')
  console.log('Параметры загрузки:', {
    sectionId: props.sectionId,
    sortBy: props.sortBy,
    searchQuery: props.searchQuery
  })
  
  try {
    // Параметры для запроса к API
    const params = {
      limit: 100, // Загружаем много товаров для фильтрации на клиенте
      is_active: true
    }
    
    // Если указан раздел, фильтруем на бекенде
    if (props.sectionId) {
      params.section_id = props.sectionId
    }
    
    // Получаем товары с бекенда
    const response = await productAPI.getAll(params)
    console.log('Ответ от API:', response.data)
    
    products.value = response.data.products || []
    totalItems.value = response.data.total || products.value.length
    
    console.log('Загружено товаров:', products.value.length)
    console.log('Пример товара:', products.value[0])
    
  } catch (error) {
    console.error('Ошибка загрузки товаров:', error)
    products.value = []
    totalItems.value = 0
  } finally {
    loading.value = false
  }
}

const changePage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    emit('page-change', page)
    // Прокручиваем к началу списка
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

// Следим за изменениями фильтров
watch(() => [props.sectionId, props.searchQuery], () => {
  console.log('Фильтры изменились, перезагружаем товары...')
  currentPage.value = 1 // Сбрасываем на первую страницу при изменении фильтров
  loadProducts()
}, { immediate: true })

// Следим за изменением сортировки
watch(() => props.sortBy, () => {
  console.log('Сортировка изменилась на:', props.sortBy)
  currentPage.value = 1 // Сбрасываем на первую страницу при изменении сортировки
})
</script>

<style scoped>
.product-list {
  flex: 1;
  padding: 20px;
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

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  text-align: center;
  color: #666;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 20px;
}

.empty-state h3 {
  margin-bottom: 10px;
  color: #333;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 30px;
  margin-bottom: 40px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-top: 40px;
  padding: 20px;
  border-top: 1px solid #e9ecef;
}

.page-btn {
  background: #667eea;
  color: white;
  border: none;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.page-btn:hover:not(:disabled) {
  background: #764ba2;
  transform: scale(1.1);
}

.page-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.page-info {
  font-size: 14px;
  color: #666;
}

@media (max-width: 768px) {
  .products-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
  }
  
  .product-list {
    padding: 15px;
  }
}

@media (max-width: 480px) {
  .products-grid {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .product-card {
    max-width: 100%;
  }

  .product-image {
    height: 180px;
  }

    .product-name {
    font-size: 16px;
  }

  .product-price {
    font-size: 18px;
  }
  
  .product-info {
    padding: 15px;
  }
}
</style>