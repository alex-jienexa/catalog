<template>
  <div class="admin-login">
    <div class="login-container">
      <div class="login-header">
        <h1>{{ isFirstAdmin ? '🚀 Первый запуск' : '🔐 Вход в админ-панель' }}</h1>
        <p v-if="isFirstAdmin">
          Администраторов ещё нет. Создайте первого администратора.
        </p>
        <p v-else>
          Для доступа к управлению каталогом требуется авторизация
        </p>
      </div>

      <!-- Форма регистрации первого администратора -->
      <form v-if="isFirstAdmin" @submit.prevent="handleRegister" class="login-form">
        <div class="form-group">
          <label for="name">Имя</label>
          <input
            v-model="form.name"
            type="text"
            id="name"
            placeholder="Введите имя"
            required
            :class="{ error: error }"
          />
        </div>

        <div class="form-group">
          <label for="username">Логин</label>
          <input
            v-model="form.username"
            type="text"
            id="username"
            placeholder="Придумайте логин"
            required
            :class="{ error: error }"
          />
        </div>

        <div class="form-group">
          <label for="password">Пароль</label>
          <input
            v-model="form.password"
            type="password"
            id="password"
            placeholder="Не менее 6 символов"
            required
            :class="{ error: error }"
          />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <button type="submit" :disabled="loading" class="login-btn">
          <span v-if="loading" class="spinner"></span>
          <span v-else>Создать администратора</span>
        </button>
      </form>

      <!-- Форма входа -->
      <form v-else @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="username">Логин</label>
          <input
            v-model="form.username"
            type="text"
            id="username"
            placeholder="Введите логин"
            required
            :class="{ error: error }"
          />
        </div>

        <div class="form-group">
          <label for="password">Пароль</label>
          <input
            v-model="form.password"
            type="password"
            id="password"
            placeholder="Введите пароль"
            required
            :class="{ error: error }"
          />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <button type="submit" :disabled="loading" class="login-btn">
          <span v-if="loading" class="spinner"></span>
          <span v-else>Войти</span>
        </button>
      </form>

      <router-link to="/" class="back-link">← Вернуться в каталог</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { auth, authAPI } from '@/services/api'

const router = useRouter()

const isFirstAdmin = ref(false)
const loading = ref(false)
const error = ref('')

const form = ref({
  name: '',
  username: '',
  password: '',
})

onMounted(async () => {
  try {
    const res = await authAPI.isFirst()
    isFirstAdmin.value = res.data.is_first
  } catch (e) {
    console.error('Не удалось проверить наличие администраторов', e)
  }
})

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  try {
    await auth.login(form.value.username, form.value.password)
    router.push('/admin')
  } catch (err) {
    error.value = err.response?.data?.error || 'Ошибка авторизации'
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  error.value = ''
  loading.value = true
  try {
    const res = await authAPI.register({
      name: form.value.name,
      username: form.value.username,
      password: form.value.password,
    })
    const { token, admin: adminData } = res.data
    localStorage.setItem('token', token)
    localStorage.setItem('adminData', JSON.stringify(adminData))
    router.push('/admin')
  } catch (err) {
    error.value = err.response?.data?.error || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-login {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-container {
  background: white;
  border-radius: 20px;
  padding: 40px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  color: #333;
  margin: 0 0 10px 0;
  font-size: 26px;
}

.login-header p {
  color: #666;
  margin: 0;
  font-size: 14px;
}

.login-form {
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-weight: 500;
  font-size: 14px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e1e5e9;
  border-radius: 8px;
  font-size: 16px;
  transition: all 0.3s;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group input.error {
  border-color: #ff6b6b;
}

.error-message {
  background: #ffeaea;
  color: #ff6b6b;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 14px;
  border: 1px solid #ffcccc;
}

.login-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 48px;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.back-link {
  display: block;
  text-align: center;
  color: #667eea;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
}

.back-link:hover {
  color: #764ba2;
  text-decoration: underline;
}

@media (max-width: 480px) {
  .login-container {
    padding: 30px 20px;
  }

  .login-header h1 {
    font-size: 22px;
  }
}
</style>
