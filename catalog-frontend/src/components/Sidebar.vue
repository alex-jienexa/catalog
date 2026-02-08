<template>
  <aside class="sidebar">
    <div class="search-section">
      <input
        v-model="searchQuery"
        @input="onSearch"
        type="text"
        placeholder="Поиск товаров..."
        class="search-input"
      />
    </div>
    
    <div class="sort-section">
      <label>Сортировка:</label>
      <select v-model="sortBy" @change="onSortChange" class="sort-select">
        <option value="created_at">По новизне</option>
        <option value="name">По названию</option>
        <option value="price">По цене (дешевые)</option>
        <option value="-price">По цене (дорогие)</option>
      </select>
    </div>
    
    <div class="sections">
      <h3>Разделы</h3>
      <ul class="section-list">
        <li>
          <button
            @click="selectSection(null)"
            :class="{ active: !selectedSection }"
            class="section-btn"
          >
            Все товары
            <span class="count">{{ totalProducts }}</span>
          </button>
        </li>
        <li v-for="section in sections" :key="section.id">
          <button
            @click="selectSection(section.id)"
            :class="{ active: selectedSection === section.id }"
            class="section-btn"
          >
            {{ section.name }}
            <span class="count">{{ section.product_count }}</span>
          </button>
        </li>
      </ul>
    </div>
  </aside>
</template>

<script setup>
import { ref, defineEmits, watch } from 'vue'
import { sectionAPI, productAPI } from '@/services/api'

const props = defineProps({
  selectedSection: Number
})

const emit = defineEmits(['section-change', 'sort-change', 'search'])

const sections = ref([])
const totalProducts = ref(0)
const sortBy = ref('created_at')
const searchQuery = ref('')

// Загружаем разделы
const loadSections = async () => {
  try {
    const response = await sectionAPI.getAll()
    sections.value = response.data
  } catch (error) {
    console.error('Ошибка загрузки разделов:', error)
  }
}

// Загружаем общее количество товаров
const loadTotalProducts = async () => {
  try {
    const response = await productAPI.getAll({ limit: 1 })
    totalProducts.value = response.data.total || 0
  } catch (error) {
    console.error('Ошибка загрузки товаров:', error)
  }
}

// Инициализация
loadSections()
loadTotalProducts()

// Обработчики событий
const selectSection = (sectionId) => {
  emit('section-change', sectionId)
}

const onSortChange = () => {
  emit('sort-change', sortBy.value)
}

const onSearch = () => {
  emit('search', searchQuery.value)
}

// Следим за изменениями разделов для обновления счетчиков
watch(sections, () => {
  if (props.selectedSection) {
    const section = sections.value.find(s => s.id === props.selectedSection)
    if (section) {
      totalProducts.value = section.product_count
    }
  }
})
</script>

<style scoped>
.sidebar {
  width: 250px;
  background: #f8f9fa;
  padding: 20px;
  border-right: 1px solid #e9ecef;
  height: calc(100vh - 70px);
  position: fixed;
  left: 0;
  top: 70px;
  overflow-y: auto;
}

.search-section {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
}

.sort-section {
  margin-bottom: 30px;
}

.sort-section label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.sort-select {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.sections h3 {
  margin-bottom: 15px;
  color: #333;
  font-size: 18px;
}

.section-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.section-list li {
  margin-bottom: 5px;
}

.section-btn {
  width: 100%;
  text-align: left;
  padding: 12px 15px;
  border: none;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #555;
  font-size: 14px;
  border: 1px solid #e9ecef;
}

.section-btn:hover {
  background: #f1f3ff;
  border-color: #667eea;
  transform: translateX(5px);
}

.section-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: #667eea;
}

.count {
  background: rgba(0, 0, 0, 0.1);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: bold;
}

.section-btn.active .count {
  background: rgba(255, 255, 255, 0.3);
}

@media (max-width: 768px) {
  .sidebar {
    position: static;
    width: 100%;
    height: auto;
    border-right: none;
    border-bottom: 1px solid #e9ecef;
  }
}
</style>