<template>
  <v-app>
    <v-app-bar app color="primary" prominent>
      <!-- <v-app-bar-nav-icon>
        <div class="logo-container">
          <v-img 
            src="@/assets/logo.png" 
            alt="Yearnstudio Logo" 
            width="40"
            height="40"
            class="logo-img"
          />
        </div>
      </v-app-bar-nav-icon> -->
      
      <v-app-bar-title class="font-weight-bold">
        Yearnstudio资源站镜像
      </v-app-bar-title>
      
      <v-spacer></v-spacer>
      
      <v-btn 
        @click="toggleTheme"
        :title="currentTheme === 'light' ? '切换到暗色主题' : '切换到亮色主题'"
      >
        <v-icon v-if="currentTheme === 'light'">mdi-weather-sunny</v-icon>
        <v-icon v-else>mdi-weather-night</v-icon>
      </v-btn>
    </v-app-bar>
    
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useTheme } from 'vuetify'

const theme = useTheme()
const currentTheme = ref('light')

// 切换主题
const toggleTheme = () => {
  currentTheme.value = theme.global.current.value.dark ? 'light' : 'dark'
  theme.global.name.value = currentTheme.value
  
  // 保存主题偏好到本地存储
  localStorage.setItem('yearnstudio-theme', currentTheme.value)
}

// 初始化主题
onMounted(() => {
  const savedTheme = localStorage.getItem('yearnstudio-theme')
  if (savedTheme) {
    currentTheme.value = savedTheme
    theme.global.name.value = savedTheme
  }
})
</script>

<style scoped>
.logo-container {
  border-radius: 0 !important; /* 移除圆角 */
  overflow: hidden; /* 确保内容不会溢出 */
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.logo-img {
  border-radius: 0 !important; /* 确保图片本身也是直角 */
}
</style>