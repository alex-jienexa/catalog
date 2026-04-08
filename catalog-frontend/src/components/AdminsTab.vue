<template>
  <div class="admins-tab">
    <div class="section-header">
      <h2>Управление администраторами</h2>
      <button @click="openCreate" class="add-btn">➕ Добавить администратора</button>
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Загрузка...</p>
    </div>

    <div v-else class="admins-list">
      <div v-for="admin in admins" :key="admin.id" class="admin-card">
        <div class="admin-avatar">{{ initials(admin.name) }}</div>
        <div class="admin-info">
          <div class="admin-name">
            {{ admin.name }}
            <span v-if="admin.id === currentAdminId" class="you-badge">Вы</span>
          </div>
          <div class="admin-username">@{{ admin.username }}</div>
          <div class="admin-since">С {{ formatDate(admin.created_at) }}</div>
        </div>
        <div class="admin-actions">
          <button @click="openEdit(admin)" class="action-btn edit" title="Редактировать">✏️</button>
          <button
            @click="handleDelete(admin)"
            class="action-btn delete"
            :disabled="admin.id === currentAdminId"
            :title="admin.id === currentAdminId ? 'Нельзя удалить себя' : 'Удалить'"
          >🗑️</button>
        </div>
      </div>
    </div>

    <!-- Модальное окно создания / редактирования -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editingAdmin ? 'Редактировать администратора' : 'Новый администратор' }}</h3>
          <button @click="closeModal" class="close-btn">✕</button>
        </div>

        <form @submit.prevent="handleSave" class="modal-form">
          <div class="form-group">
            <label>Имя</label>
            <input v-model="form.name" type="text" placeholder="Имя администратора" required />
          </div>

          <div class="form-group">
            <label>Логин</label>
            <input v-model="form.username" type="text" placeholder="Логин для входа" required />
          </div>

          <div class="form-group">
            <label>{{ editingAdmin ? 'Новый пароль (оставьте пустым, чтобы не менять)' : 'Пароль' }}</label>
            <input
              v-model="form.password"
              type="password"
              :placeholder="editingAdmin ? 'Не менее 6 символов' : 'Не менее 6 символов'"
              :required="!editingAdmin"
            />
          </div>

          <div v-if="formError" class="form-error">{{ formError }}</div>

          <div class="modal-footer">
            <button type="button" @click="closeModal" class="cancel-btn">Отмена</button>
            <button type="submit" :disabled="saving" class="save-btn">
              <span v-if="saving" class="spinner-sm"></span>
              <span v-else>{{ editingAdmin ? 'Сохранить' : 'Создать' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { adminAPI, auth } from '@/services/api'

const admins = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const editingAdmin = ref(null)
const formError = ref('')

const currentAdminId = auth.getAdminData()?.id ?? null

const form = ref({ name: '', username: '', password: '' })

const initials = (name) => {
  return name
    .split(' ')
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase() ?? '')
    .join('')
}

const formatDate = (iso) => {
  return new Date(iso).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

const loadAdmins = async () => {
  loading.value = true
  try {
    const res = await adminAPI.getAdmins()
    admins.value = res.data
  } catch (e) {
    console.error('Ошибка загрузки администраторов:', e)
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editingAdmin.value = null
  form.value = { name: '', username: '', password: '' }
  formError.value = ''
  showModal.value = true
}

const openEdit = (admin) => {
  editingAdmin.value = admin
  form.value = { name: admin.name, username: admin.username, password: '' }
  formError.value = ''
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const handleSave = async () => {
  formError.value = ''
  saving.value = true
  try {
    if (editingAdmin.value) {
      const payload = { name: form.value.name, username: form.value.username }
      if (form.value.password) payload.password = form.value.password
      await adminAPI.updateAdmin(editingAdmin.value.id, payload)
    } else {
      await adminAPI.createAdmin({
        name: form.value.name,
        username: form.value.username,
        password: form.value.password,
      })
    }
    closeModal()
    await loadAdmins()
  } catch (e) {
    formError.value = e.response?.data?.error || 'Произошла ошибка'
  } finally {
    saving.value = false
  }
}

const handleDelete = async (admin) => {
  if (!confirm(`Удалить администратора «${admin.name}»?`)) return
  try {
    await adminAPI.deleteAdmin(admin.id)
    await loadAdmins()
  } catch (e) {
    alert(e.response?.data?.error || 'Ошибка при удалении')
  }
}

onMounted(loadAdmins)
</script>

<style scoped>
.admins-tab {}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.section-header h2 {
  margin: 0;
  color: #333;
  font-size: 24px;
}

.add-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px;
  color: #666;
}

.admins-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.admin-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 10px;
  padding: 16px 20px;
  transition: border-color 0.2s;
}

.admin-card:hover {
  border-color: #667eea;
}

.admin-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-weight: 700;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.admin-info {
  flex: 1;
}

.admin-name {
  font-weight: 600;
  color: #333;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.you-badge {
  background: #667eea;
  color: white;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.admin-username {
  color: #667eea;
  font-size: 14px;
  margin-top: 2px;
}

.admin-since {
  color: #999;
  font-size: 12px;
  margin-top: 2px;
}

.admin-actions {
  display: flex;
  gap: 6px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  transition: all 0.2s;
}

.action-btn.edit { background: #ffeaa7; }
.action-btn.edit:hover { background: #fab1a0; transform: scale(1.1); }
.action-btn.delete { background: #ffcccc; color: #ff6b6b; }
.action-btn.delete:hover:not(:disabled) { background: #ff6b6b; color: white; transform: scale(1.1); }
.action-btn:disabled { opacity: 0.4; cursor: not-allowed; }

/* Модальное окно */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 440px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px 0;
}

.modal-header h3 { margin: 0; color: #333; font-size: 20px; }

.close-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #666;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}
.close-btn:hover { background: #f0f0f0; }

.modal-form { padding: 20px 24px 24px; }

.form-group { margin-bottom: 18px; }

.form-group label {
  display: block;
  margin-bottom: 6px;
  color: #555;
  font-weight: 500;
  font-size: 14px;
}

.form-group input {
  width: 100%;
  padding: 10px 14px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 15px;
  box-sizing: border-box;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-error {
  background: #ffeaea;
  color: #ff6b6b;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 14px;
  margin-bottom: 16px;
  border: 1px solid #ffcccc;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.cancel-btn {
  padding: 10px 20px;
  background: #f0f0f0;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
  transition: background 0.2s;
}
.cancel-btn:hover { background: #e0e0e0; }

.save-btn {
  padding: 10px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 15px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 100px;
  justify-content: center;
  transition: opacity 0.2s;
}
.save-btn:disabled { opacity: 0.7; cursor: not-allowed; }

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

.spinner-sm {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.4);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
