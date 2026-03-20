<template>
  <header class="header">
    <div class="container">
        <router-link to="/" class="logo">
          🛒 Магазин
        </router-link>
        <button class="mobile-nav-toggle" @click="toggleNav" v-if="isMobile">
          <span class="burger-icon"></span>
        </button>
        
        <div v-if="mobileNavOpen && isMobile" class="nav-overlay" @click="closeNav"></div>

        <nav class="nav" :class="{ 'open': mobileNavOpen }">
          <router-link to="/" class="nav-link" @click="closeNav">Главная</router-link>
          <router-link to="/about" class="nav-link" @click="closeNav">О магазине</router-link>
          <router-link v-if="isAuthenticated" to="/admin" class="nav-link admin-link" @click="closeNav">👑 Админка</router-link>
          <router-link v-else to="/admin/login" class="nav-link" @click="closeNav">🔑 Вход</router-link>
        </nav>
    </div>
  </header>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { auth } from '@/services/api'
import { useRouter } from 'vue-router';

const router = useRouter()
const isMobile = ref(window.innerWidth <= 768)
const isAuthenticated = computed(() => auth.isAuthenticated())
const mobileNavOpen = ref(false)

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
  if (!isMobile.value) mobileNavOpen.value = false
}

const handleScroll = () => {
  if (mobileNavOpen.value) {
    closeNav()
  }
}

onMounted(() => {
  window.addEventListener('resize', checkMobile)
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  window.removeEventListener('scroll', handleScroll)
})

const toggleNav = () => {
  mobileNavOpen.value = !mobileNavOpen.value
}

const closeNav = () => {
  mobileNavOpen.value = false
}

router.afterEach(() => {
  closeNav()
})
</script>

<style scoped>
.header {
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1rem 0;
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 101;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
  text-decoration: none;
  color: white;
  transition: transform 0.3s ease;
}

.logo:hover {
  transform: scale(1.05);
}

.nav-links {
  display: flex;
  gap: 2rem;
  align-items: center;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: all 0.3s ease;
  margin: 0 20px;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.1);
}

.nav-link.router-link-active {
  background: rgba(255, 255, 255, 0.2);
}

.admin-link {
  background: rgba(255, 215, 0, 0.2);
  border: 1px solid gold;
}

.admin-link:hover {
  background: rgba(255, 215, 0, 0.3);
}

.mobile-nav-toggle {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  width: 40px;
  height: 40px;
  position: relative;
  z-index: 102;
}

.burger-icon,
.burger-icon::before,
.burger-icon::after {
  content: '';
  display: block;
  width: 24px;
  height: 2px;
  background: white;
  position: absolute;
  transition: all 0.3s;
}

.burger-icon {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.burger-icon::before {
  top: -8px;
}

.burger-icon::after {
  bottom: -8px;
}

.nav-overlay {
  position: fixed;
  top: 70px;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 100;
}

@media (max-width: 768px) {
  .mobile-nav-toggle {
    display: block;
  }

  .nav {
    position: fixed;
    top: 70px;
    left: 0;
    right: 0;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    flex-direction: column;
    gap: 0;
    transform: translateY(-200%);
    transition: transform 0.3s ease;
    z-index: 101;
    padding: 10px 0;
  }

  .nav.open {
    transform: translateY(0);
  }

  .nav-link {
    padding: 10px;
    margin: 10px;
    width: 100%;
  }
}
</style>