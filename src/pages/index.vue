<template>
  <v-container class="fill-height d-flex flex-column justify-center align-center">
    <v-card width="700" class="pa-6" elevation="8">
      <v-card-title class="text-center mb-4 d-flex flex-column align-center">
        <div class="logo-container mb-3">
          <v-img 
            src="@/assets/logo.png" 
            alt="Yearnstudio Logo" 
            width="80"
            height="80"
            class="logo-img"
          />
        </div>
        欢迎使用Yearnstudio资源站镜像站
      </v-card-title>
      
      <v-card-text>
        <!-- 加载状态 -->
        <div v-if="loading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
          <div class="mt-2">正在获取站点列表...</div>
        </div>
        
        <!-- 错误提示 -->
        <v-alert
          v-if="error"
          type="error"
          class="mb-4"
        >
          {{ error }}
        </v-alert>
        
        <!-- 站点列表 -->
        <v-list v-if="!loading && !error && mirrors.length > 0" lines="two" class="mb-4">
          <v-list-subheader class="text-h6">
            <span>镜像站点</span>
          </v-list-subheader>
          <v-list-item
            v-for="mirror in sortedMirrors"
            :key="mirror.id"
            :title="mirror.name"
            :subtitle="mirror.url"
          >
            <template v-slot:prepend>
              <v-icon color="primary">mdi-server</v-icon>
            </template>
            <template v-slot:append>
                  <div class="text-right">
                    <div v-if="mirror.latency !== null" class="d-flex align-center" :class="speedTestComplete ? 'animate-pulse' : ''">
                      <v-icon 
                        :color="getLatencyColor(mirror.latency)" 
                        size="small" 
                        class="mr-1"
                      >
                        {{ getLatencyIcon(mirror.latency) }}
                      </v-icon>
                      <span :class="`text-${getLatencyColor(mirror.latency)}`">
                        {{ (mirror.latency / 1000).toFixed(1) }}s
                      </span>
                    </div>
                    <div v-else-if="speedTestComplete" class="text-error">无法连接</div>
                    <div v-else class="text-grey">检测中</div>
                  </div>
                </template>
          </v-list-item>
        </v-list>
        
        <!-- 空状态提示 -->
        <div v-if="!loading && !error && mirrors.length === 0" class="text-center py-8">
          <v-icon size="48" color="grey">mdi-server-off</v-icon>
          <div class="mt-2 text-grey">暂无可用站点</div>
        </div>
        
        <v-divider class="my-4" />
        
        <div class="text-center">
          <v-btn 
            color="primary" 
            size="x-large" 
            @click="$router.push('/use')"
            prepend-icon="mdi-play"
            class="mb-2 mr-2"
            :disabled="mirrors.length === 0"
          >
            开始使用
          </v-btn>
          
          <v-btn 
            color="secondary" 
            size="x-large" 
            @click="refreshData"
            prepend-icon="mdi-refresh"
            class="mb-2"
            :loading="loading"
          >
            刷新
          </v-btn>
          
          <div class="text-caption text-grey">点击开始使用按钮进入镜像资源管理页面</div>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// 镜像站点数据
const mirrors = ref([])
const loading = ref(false)
const error = ref('')

// 获取镜像站点列表
const fetchMirrorList = async () => {
  loading.value = true
  error.value = ''
  
  try {
    // 确定API端点
    const isLocalhost = window.location.hostname === 'localhost'
    const apiEndpoint = isLocalhost 
      ? 'http://localhost:27645/api/list' 
      : 'https://api.mirror.yearnstudio.cn/api/list'
    
    // 发起API请求
    const response = await fetch(apiEndpoint, { 
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      },
      cache: 'no-cache' // 不使用缓存
    })
    
    // 检查响应状态
    if (!response.ok) {
      // 处理特定的HTTP状态码
      if (response.status === 400) {
        throw new Error('请求参数错误')
      } else if (response.status === 404) {
        throw new Error('API端点未找到')
      } else if (response.status === 500) {
        throw new Error('服务器内部错误')
      } else {
        throw new Error(`API请求失败: ${response.status}`)
      }
    }
    
    // 获取响应数据
    const data = await response.json()
    
    // 处理响应数据
    if (data && Array.isArray(data)) {
      // 为每个站点添加id和latency字段
      mirrors.value = data.map((item, index) => ({
        id: index + 1,
        name: item.name || '未命名站点',
        url: item.url,
        latency: null
      }))
    } else {
      throw new Error('响应数据格式不正确')
    }
    
  } catch (err) {
    error.value = '获取站点列表失败：' + err.message
    console.error('获取站点列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 检测延迟状态
const testingSpeed = ref(false)
const speedTestComplete = ref(false)

// 按延迟排序的镜像站点
const sortedMirrors = computed(() => {
  // 如果有延迟数据，按延迟从低到高排序
  if (mirrors.value.some(m => m.latency !== null)) {
    return [...mirrors.value].sort((a, b) => {
      // 将null延迟放在最后
      if (a.latency === null && b.latency === null) return 0
      if (a.latency === null) return 1
      if (b.latency === null) return -1
      return a.latency - b.latency
    })
  }
  // 如果没有延迟数据，返回原始顺序
  return mirrors.value
})

// 测试单个镜像站点的延迟
const testMirrorLatency = async (url) => {
  try {
    // 使用fetch API测试延迟
    const startTime = Date.now()
    
    // 创建一个超时Promise，如果请求超过15秒就认为失败
    const timeoutPromise = new Promise((_, reject) => {
      setTimeout(() => reject(new Error('请求超时')), 15000)
    })
    
    // 发起请求
    const fetchPromise = fetch(url, { 
      method: 'HEAD', // 只获取头部信息，减少数据传输
      mode: 'no-cors', // 使用no-cors模式避免CORS问题
      cache: 'no-cache' // 不使用缓存
    })
    
    // 等待请求完成或超时
    await Promise.race([fetchPromise, timeoutPromise])
    
    // 计算延迟
    const latency = Date.now() - startTime
    return latency
  } catch (error) {
    // 如果请求失败，返回null表示无法连接
    console.error(`测试 ${url} 延迟失败:`, error)
    return null
  }
}

// 测试所有镜像站点的延迟
const testAllMirrors = async () => {
  testingSpeed.value = true
  speedTestComplete.value = false
  
  // 重置所有延迟数据
  mirrors.value.forEach(mirror => {
    mirror.latency = null
  })
  
  // 并行测试所有镜像站点
  const testPromises = mirrors.value.map(async (mirror) => {
    const latency = await testMirrorLatency(mirror.url)
    mirror.latency = latency
    return { id: mirror.id, latency }
  })
  
  // 等待所有测试完成
  await Promise.all(testPromises)
  
  testingSpeed.value = false
  speedTestComplete.value = true
}

// 根据延迟获取颜色
const getLatencyColor = (latency) => {
  if (latency === null) return 'grey'
  if (latency < 5000) return 'success' // 绿色：延迟低于5秒
  if (latency < 8000) return 'warning' // 黄色：延迟5-8秒
  return 'error' // 红色：延迟高于8秒
}

// 根据延迟获取图标
const getLatencyIcon = (latency) => {
  if (latency === null) return 'mdi-help-circle'
  if (latency < 5000) return 'mdi-check-circle'
  if (latency < 8000) return 'mdi-alert-circle'
  return 'mdi-alert'
}

// 刷新数据
const refreshData = async () => {
  await fetchMirrorList()
  if (mirrors.value.length > 0) {
    testAllMirrors()
  }
}

// 页面加载时自动获取站点列表并检测延迟
onMounted(async () => {
  await fetchMirrorList()
  if (mirrors.value.length > 0) {
    testAllMirrors()
  }
})
</script>

<style scoped>
.logo-container {
  display: flex;
  justify-content: center;
}

.logo-img {
  border-radius: 0 !important;
}

.animate-pulse {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
  100% {
    opacity: 1;
  }
}
</style>