<template>
  <Header />
  <main>
    <router-view />
  </main>
</template>

<script>
import Header from './components/Header.vue'

export default {
  name: 'App',
  components: {
    Header
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  max-width: 100%;
}

body {
  overflow-x: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #f5f7fa;
}

#app {
  min-height: 100vh;
}

main {
  min-height: calc(100vh - 70px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  /* display: flex;
  justify-content: space-between; */
}

@media (max-width: 768px) {
  .container {
    padding: 0 15px;
  }
}
</style>