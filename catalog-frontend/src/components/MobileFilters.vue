<template>
  <div class="mobile-filters">
    <!-- Верхняя панель -->
    <div class="filter-bar">
      <button class="section-btn" @click="openSectionModal">
        <span class="label">Раздел:</span>
        <span class="value">{{ selectedSectionName }}</span>
        <span class="arrow">▼</span>
      </button>
      <div class="search-group">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="Поиск..." 
          @input="onSearch"
          class="search-input"
        />
        <button class="filter-toggle" @click="openSortModal">
          ⚙️
        </button>
      </div>
    </div>

    <!-- Модальное окно выбора раздела -->
    <div v-if="showSectionModal" class="modal-overlay" @click.self="closeSectionModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Выберите раздел</h3>
          <button class="close-btn" @click="closeSectionModal">×</button>
        </div>
        <ul class="modal-list">
          <li>
            <button 
              @click="selectSection(null)"
              :class="{ active: selectedSection === null }"
            >
              Все товары
            </button>
          </li>
          <li v-for="section in sections" :key="section.id">
            <button 
              @click="selectSection(section.id)"
              :class="{ active: selectedSection === section.id }"
            >
              {{ section.name }}
              <span class="count">{{ section.product_count }}</span>
            </button>
          </li>
        </ul>
      </div>
    </div>

    <!-- Модальное окно сортировки -->
    <div v-if="showSortModal" class="modal-overlay" @click.self="closeSortModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Сортировка</h3>
          <button class="close-btn" @click="closeSortModal">×</button>
        </div>
        <ul class="modal-list">
          <li>
            <button 
              @click="setSort('created_at')"
              :class="{ active: currentSort === 'created_at' }"
            >
              По новизне
            </button>
          </li>
          <li>
            <button 
              @click="setSort('name')"
              :class="{ active: currentSort === 'name' }"
            >
              По названию
            </button>
          </li>
          <li>
            <button 
              @click="setSort('price')"
              :class="{ active: currentSort === 'price' }"
            >
              По цене (дешевые)
            </button>
          </li>
          <li>
            <button 
              @click="setSort('-price')"
              :class="{ active: currentSort === '-price' }"
            >
              По цене (дорогие)
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  sections: Array,
  selectedSection: Number,
  sortBy: String,
  searchQuery: String
})

const emit = defineEmits(['section-change', 'sort-change', 'search'])

// Локальное состояние модалок
const showSectionModal = ref(false)
const showSortModal = ref(false)

// Текущие значения (пропсы)
const selectedSection = ref(props.selectedSection)
const currentSort = ref(props.sortBy)
const searchQuery = ref(props.searchQuery || '')

// Название выбранного раздела для отображения
const selectedSectionName = computed(() => {
  if (!selectedSection.value) return 'Все товары'
  const section = props.sections.find(s => s.id === selectedSection.value)
  return section ? section.name : 'Все товары'
})

// Методы
const openSectionModal = () => {
  showSectionModal.value = true
}

const closeSectionModal = () => {
  showSectionModal.value = false
}

const selectSection = (id) => {
  selectedSection.value = id
  emit('section-change', id)
  closeSectionModal()
}

const openSortModal = () => {
  showSortModal.value = true
}

const closeSortModal = () => {
  showSortModal.value = false
}

const setSort = (value) => {
  currentSort.value = value
  emit('sort-change', value)
  closeSortModal()
}

const onSearch = () => {
  emit('search', searchQuery.value)
}
</script>

<style scoped>
.mobile-filters {
  background: white;
  padding: 10px 15px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  margin-bottom: 15px;
  border-radius: 8px;
}

.filter-bar {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.section-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 10px 12px;
  background: #f5f7fa;
  border: 1px solid #e0e4e8;
  border-radius: 8px;
  width: 100%;
  cursor: pointer;
  font-size: 14px;
}

.section-btn .label {
  color: #666;
  font-weight: 500;
}

.section-btn .value {
  flex: 1;
  text-align: left;
  font-weight: 600;
  color: #333;
}

.section-btn .arrow {
  color: #667eea;
  font-size: 12px;
}

.search-group {
  display: flex;
  gap: 8px;
}

.search-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #e0e4e8;
  border-radius: 8px;
  font-size: 14px;
}

.filter-toggle {
  width: 44px;
  height: 44px;
  background: #f5f7fa;
  border: 1px solid #e0e4e8;
  border-radius: 8px;
  cursor: pointer;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Модальные окна */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: flex-end;
  z-index: 1000;
}

.modal-content {
  background: white;
  width: 100%;
  border-radius: 20px 20px 0 0;
  max-height: 70vh;
  overflow-y: auto;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.modal-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.modal-list li {
  border-bottom: 1px solid #f0f0f0;
}

.modal-list button {
  width: 100%;
  padding: 15px 20px;
  background: none;
  border: none;
  text-align: left;
  font-size: 16px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-list button.active {
  background: #f0f3ff;
  color: #667eea;
  font-weight: 500;
}

.count {
  background: #e0e4e8;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  color: #666;
}

@media (min-width: 769px) {
  .mobile-filters {
    display: none;
  }
}
</style>