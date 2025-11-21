package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var mirrorList []string = []string{"https://yearnstudio.cdn.houlang.cloud/d/123-test", "https://resource.yearnstudio.cn/d/123pan/pan", "https://list.yearnstudio.cn/d/123pan/pan"}

// var mirrorDomain []string = []string{"https://yearnstudio.cdn.houlang.cloud", "https://resource.yearnstudio.cn", "https://list.yearnstudio.cn"}
var mirrorDomain map[string]string = map[string]string{
	"主站":          "https://resource.yearnstudio.cn",
	"厚浪云镜像":       "https://yearnstudio.cdn.houlang.cloud",
	"EdgeOne海外线路": "https://list.yearnstudio.cn",
}

type SiteInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	fmt.Println("[Info]Resource Mirror Backend")
	debug := flag.Bool("debug", false, "启用调试模式")
	flag.Parse()
	if *debug {
		fmt.Println("[Info]调试模式已开启")
	} else {
		fmt.Println("[Info]调试模式已关闭")
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://mirror.yearnstudio.cn", "http://localhost:3000"}, // 允许的源
		AllowMethods:     []string{"GET", "POST"},                                            // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type"},                                 // 允许的头部
		ExposeHeaders:    []string{"Content-Length"},                                         // 暴露的头部
		AllowCredentials: true,                                                               // 是否允许携带凭证
		MaxAge:           12 * time.Hour,                                                     // 缓存时间
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Online",
		})
	})

	r.GET("/api/list", func(c *gin.Context) {
		var sites []SiteInfo             // 结果储存
		for k, v := range mirrorDomain { // 格式化数据
			sites = append(sites, SiteInfo{Name: k, URL: v})
		}
		c.JSON(http.StatusOK, sites)
	})

	r.POST("/api/get", func(ctx *gin.Context) {
		// 获取请求体
		var reqBody struct {
			URL string `json:"url"`
		} // 使用json标签令gin框架解析JSON中的url参数
		if err := ctx.ShouldBindJSON(&reqBody); err != nil { // URL字段不存在
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		url := reqBody.URL // 获取到url
		if strings.HasPrefix(url, "http://") {
			url = strings.Replace(url, "http://", "https://", 1)
		}
		for _, v := range mirrorList {
			if strings.HasPrefix(url, v) { // 匹配域名
				url = strings.Replace(url, v, "", 1) // 替换开头的部分
				var result []string = []string{}     // 初始化结果储存
				for _, baseUrl := range mirrorList {
					newUrl := baseUrl + url
					result = append(result, newUrl) // 添加
				}
				ctx.JSON(http.StatusOK, result) // 返回列表
				return
			}
		}
		// 没有找到匹配项目
		ctx.JSON(http.StatusBadRequest, gin.H{})
	})

	r.Run(":27645")
}
