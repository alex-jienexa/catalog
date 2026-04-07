<template>
  <div class="reservations-list">
    <h2>Управление бронированиями</h2>
    
    <div v-if="loading" class="loading">Загрузка...</div>
    
    <div v-else-if="groupedReservations.length === 0" class="empty-state">
      <div class="empty-icon">📋</div>
      <p>Пока нет ни одного бронирования</p>
    </div>
    
    <div v-else class="reservations-container">
      <!-- Группы по клиентам -->
      <div v-for="group in groupedReservations" :key="group.customerId" class="customer-group">
        <div class="customer-header">
          <div class="customer-info">
            <span class="customer-icon">👤</span>
            <strong>{{ group.customerName }}</strong>
            <span class="customer-phone">{{ group.customerPhone }}</span>
          </div>
          <div class="customer-stats">
            Всего заявок: {{ group.reservations.length }}
          </div>
        </div>
        
        <div class="reservations-table">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Товар</th>
                <th>Цена</th>
                <th>Дата</th>
                <th>Статус</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="res in group.reservations" :key="res.id">
                <td>{{ res.id }}</td>
                <td>{{ res.product_name }}</td>
                <td>{{ formatPrice(res.product_price) }} ₽</td>
                <td>{{ formatDate(res.created_at) }}</td>
                <td>
                  <span :class="['status-badge', res.status]">
                    {{ getStatusText(res.status) }}
                  </span>
                </td>
                <td>
                  <select v-model="res.status" @change="updateStatus(res.id, res.status)" class="status-select">
                    <option value="pending">Ожидает</option>
                    <option value="contacted">Связались</option>
                    <option value="completed">Завершено</option>
                  </select>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <div class="pagination" v-if="totalPages > 1">
        <button @click="changePage(currentPage-1)" :disabled="currentPage===1">←</button>
        <span>Страница {{ currentPage }} из {{ totalPages }}</span>
        <button @click="changePage(currentPage+1)" :disabled="currentPage===totalPages">→</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { reservationAPI } from '@/services/api'

const reservations = ref([])
const loading = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)
const limit = 20

// Группировка по клиентам
const groupedReservations = computed(() => {
  const groups = new Map()
  for (const res of reservations.value) {
    const customerId = res.customer_id
    if (!groups.has(customerId)) {
      groups.set(customerId, {
        customerId,
        customerName: `${res.customer_first_name} ${res.customer_last_name || ''}`.trim(),
        customerPhone: res.customer_phone,
        reservations: []
      })
    }
    groups.get(customerId).reservations.push(res)
  }
  // Сортировка групп: сначала клиенты с самыми старыми ожидающими заявками (сохраняем порядок из БД)
  return Array.from(groups.values())
})

const load = async () => {
  loading.value = true
  try {
    const res = await reservationAPI.getAll({ page: currentPage.value, limit })
    reservations.value = res.data.reservations || []
    const total = res.data.total || 0
    totalPages.value = Math.ceil(total / limit)
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const updateStatus = async (id, newStatus) => {
  try {
    await reservationAPI.updateStatus(id, { status: newStatus })
    // Опционально: показать уведомление
  } catch (e) {
    console.error('Ошибка обновления статуса:', e)
    alert('Не удалось обновить статус')
    // Откатываем изменение в UI
    await load()
  }
}

const changePage = (page) => {
  currentPage.value = page
  load()
}

const formatDate = (date) => new Date(date).toLocaleString('ru-RU')
const formatPrice = (price) => new Intl.NumberFormat('ru-RU').format(price)
const getStatusText = (status) => {
  const map = { pending: 'Ожидает', contacted: 'Связались', completed: 'Завершено' }
  return map[status] || status
}

onMounted(load)
</script>

<style scoped>
.reservations-list {
  padding: 20px;
}
.customer-group {
  background: #fff;
  border-radius: 12px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  overflow: hidden;
}
.customer-header {
  background: #f7fafc;
  padding: 16px 20px;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
}
.customer-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}
.customer-icon {
  font-size: 20px;
}
.customer-phone {
  color: #4a5568;
  font-size: 14px;
}
.customer-stats {
  font-size: 14px;
  color: #718096;
}
.reservations-table {
  overflow-x: auto;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th, td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #edf2f7;
}
th {
  background: #f9fafb;
  font-weight: 600;
}
.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}
.status-badge.pending {
  background: #fefcbf;
  color: #975a16;
}
.status-badge.contacted {
  background: #bee3f8;
  color: #2c5282;
}
.status-badge.completed {
  background: #c6f6d5;
  color: #22543d;
}
.status-select {
  padding: 6px 12px;
  border-radius: 6px;
  border: 1px solid #cbd5e0;
  background: white;
  cursor: pointer;
}
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  margin-top: 30px;
}
.pagination button {
  background: #edf2f7;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
}
.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #718096;
}
.empty-icon {
  font-size: 48px;
  margin-bottom: 20px;
}
.loading {
  text-align: center;
  padding: 40px;
}
</style>