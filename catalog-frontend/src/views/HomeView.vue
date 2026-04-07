<template>
  <div class="home">
    <Sidebar
      v-if="!isMobile"
      :selected-section="selectedSection"
      @section-change="handleSectionChange"
      @sort-change="handleSortChange"
      @search="handleSearch"
    />
    <MobileFilters
      v-else
      :sections="sections"
      :selected-section="selectedSection"
      :sort-by="sortBy"
      :search-query="searchQuery"
      @section-change="handleSectionChange"
      @sort-change="handleSortChange"
      @search="handleSearch"
    />

    <main class="main-content" :class="{ 'mobile': isMobile }">
      <ProductList
        :section-id="selectedSection"
        :sort-by="sortBy"
        :search-query="searchQuery"
        :sections="sections"
        @page-change="handlePageChange"
      />
      <BookingModal
        :visible="bookingModalVisible"
        :product-id="selectedProductId"
        @close="bookingModalVisible = false"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import Sidebar from '@/components/Sidebar.vue'
import ProductList from '@/components/ProductList.vue'
import MobileFilters from '@/components/MobileFilters.vue'
import { sectionAPI } from '@/services/api'
import BookingModal from '@/components/BookingModal.vue'

const selectedSection = ref(null)
const sortBy = ref('created_at')
const searchQuery = ref('')
const sections = ref([])
const isMobile = ref(window.innerWidth <= 768)

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

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
  if (!isMobile.value) sidebarOpen.value = false
}

onMounted(() => {
  window.addEventListener('resize', checkMobile)
  loadSections()
})
onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

const loadSections = async () => {
  try {
    const response = await sectionAPI.getAll(true)
    sections.value = response.data || []
  } catch (error) {
    console.error('Ошибка загрузки разделов:', error)
  }
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

.main-content.mobile {
  margin-left: 0;
  padding: 10px;
}

@media (max-width: 768px) {
  .home {
    flex-direction: column;
  }
}
</style>