<template>
  <v-container class="fill-height d-flex flex-column justify-center align-center">
    <v-card width="600" class="pa-6" elevation="8">
      <v-card-title class="text-center mb-4 d-flex align-center justify-center">
        文件镜像获取
      </v-card-title>
      
      <v-card-text>
        <v-form @submit.prevent="fetchMirrorList">
          <v-text-field
            v-model="mirrorUrl"
            label="镜像站点URL"
            placeholder="请输入镜像站点URL"
            prepend-icon="mdi-link"
            :rules="[v => !!v || '请输入镜像站点URL']"
            class="mb-4"
          />
          
          <div class="text-center">
            <v-btn 
              type="submit"
              color="primary" 
              size="x-large" 
              :loading="loading"
              prepend-icon="mdi-download"
              class="mb-4"
            >
              获取列表
            </v-btn>
          </div>
        </v-form>
        
        <v-divider class="my-4" />
        
        <div v-if="mirrorList.length > 0">
          <v-list-subheader class="text-h6">镜像列表</v-list-subheader>
          <v-list lines="two" class="mb-4">
            <v-list-item
              v-for="(item, index) in mirrorList"
              :key="index"
              :subtitle="item.url || item"
            >
              <template v-slot:prepend>
                <v-icon color="success">mdi-file-download</v-icon>
              </template>
              <template v-slot:append>
                <v-btn 
                  color="primary" 
                  size="small"
                  @click="downloadItem(item)"
                  prepend-icon="mdi-download"
                >
                  下载
                </v-btn>
              </template>
            </v-list-item>
          </v-list>
        </div>
        
        <v-alert
          v-if="error"
          type="error"
          class="mt-4"
        >
          {{ error }}
        </v-alert>
        
        <div class="text-center mt-4">
          <v-btn 
            color="secondary" 
            @click="$router.push('/')"
            prepend-icon="mdi-arrow-left"
          >
            返回首页
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref } from 'vue'

const mirrorUrl = ref('')
const mirrorList = ref([])
const loading = ref(false)
const error = ref('')

// 获取镜像列表的函数
const fetchMirrorList = async () => {
  if (!mirrorUrl.value) {
    error.value = '请输入镜像站点URL'
    return
  }
  
  loading.value = true
  error.value = ''
  mirrorList.value = []
  
  try {
    // 确定API端点
    const isLocalhost = window.location.hostname === 'localhost'
    const apiEndpoint = isLocalhost 
      ? 'http://localhost:27645/api/get' 
      : 'https://api.mirror.yearnstudio.cn/api/get'
    
    // 发起API请求
    const response = await fetch(apiEndpoint, { 
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ url: mirrorUrl.value }),
      cache: 'no-cache' // 不使用缓存
    })
    
    // 检查响应状态
    if (!response.ok) {
      // 处理特定的HTTP状态码
      if (response.status === 400) {
        throw new Error('请求参数错误，请检查输入的URL是否正确')
      } else if (response.status === 404) {
        throw new Error('API端点未找到，请联系管理员')
      } else if (response.status === 500) {
        throw new Error('服务器内部错误，请稍后重试')
      } else {
        throw new Error(`API请求失败: ${response.status}`)
      }
    }
    
    // 获取响应数据
    const data = await response.json()
    
    // 处理响应数据
    if (data && Array.isArray(data)) {
      mirrorList.value = data
    } else if (data && data.list && Array.isArray(data.list)) {
      mirrorList.value = data.list
    } else {
      throw new Error('响应数据格式不正确')
    }
    
    // 确保每个资源项都有必要的字段
    mirrorList.value = mirrorList.value.map(item => {
      // 如果item是字符串，则转换为对象格式
      if (typeof item === 'string') {
        return { url: item }
      }
      // 如果item是对象，确保有url字段
      return {
        url: item.url || item.path || item,
        // 不再需要name字段
      }
    })
    
  } catch (err) {
    error.value = '获取镜像列表失败：' + err.message
  } finally {
    loading.value = false
  }
}

// 下载项目
const downloadItem = (item) => {
  // 处理不同格式的数据
  const url = typeof item === 'string' ? item : item.url || item.path || item
  if (url) {
    window.open(url, '_blank')
  }
}
</script>

<style scoped>
.logo-container {
  display: flex;
  align-items: center;
}

.logo-img {
  border-radius: 0 !important;
}
</style>