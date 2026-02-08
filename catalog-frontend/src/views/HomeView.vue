<template>
  <div class="home">
    <Sidebar
      :selected-section="selectedSection"
      @section-change="handleSectionChange"
      @sort-change="handleSortChange"
      @search="handleSearch"
    />
    <main class="main-content">
      <ProductList
        :section-id="selectedSection"
        :sort-by="sortBy"
        :search-query="searchQuery"
        :sections="sections"
        @page-change="handlePageChange"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Sidebar from '@/components/Sidebar.vue'
import ProductList from '@/components/ProductList.vue'
import { sectionAPI } from '@/services/api'

const selectedSection = ref(null)
const sortBy = ref('created_at')
const searchQuery = ref('')
const sections = ref([])

console.log('HomeView инициализирован')

// Загружаем разделы при монтировании
onMounted(async () => {
  console.log('Загрузка разделов...')
  try {
    const response = await sectionAPI.getAll(true)
    sections.value = response.data || []
    console.log('Загружено разделов:', sections.value.length)
    console.log('Разделы:', sections.value)
  } catch (error) {
    console.error('Ошибка загрузки разделов:', error)
  }
})

// Обработчики событий
const handleSectionChange = (sectionId) => {
  console.log('Изменен раздел:', sectionId)
  selectedSection.value = sectionId
}

const handleSortChange = (newSortBy) => {
  console.log('Изменена сортировка:', newSortBy)
  sortBy.value = newSortBy
}

const handleSearch = (query) => {
  console.log('Поиск:', query)
  searchQuery.value = query
}

const handlePageChange = (page) => {
  console.log('Изменена страница:', page)
  // Можно добавить скролл к началу списка
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<style scoped>
.home {
  display: flex;
  min-height: calc(100vh - 70px);
  margin-top: 70px;
}

.main-content {
  flex: 1;
  margin-left: 250px;
  background: #f5f7fa;
  min-height: calc(100vh - 70px);
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  
  .home {
    flex-direction: column;
  }
}
</style>