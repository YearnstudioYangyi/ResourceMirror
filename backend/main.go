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

type SiteInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// 新版: 站点信息
//
//	站点类型常量
type SiteType int

const (
	DownloadProxy = iota
	ShareLink
)

// 站点信息结构体
type MirrorSite struct {
	Name     string             `json:"name"`     // 站点名称
	Url      string             `json:"url"`      // 显示URL(DownloadProxy)  |  分享链接(ShareLink)
	BaseUrl  string             `json:"base_url"` // 下载基URL(DownloadProxy)  |  文件开始路径(ShareLink)
	Showed   bool               `json:"showed"`   // 站点是否显示
	SiteType `json:"site_type"` // 站点类型
}

var sites []MirrorSite = []MirrorSite{
	{"主站", "https://resource.yearnstudio.cn", "https://resource.yearnstudio.cn/d/123pan/pan", true, DownloadProxy},
	{"厚浪云镜像", "https://yearnstudio.cdn.houlang.cloud", "https://yearnstudio.cdn.houlang.cloud/d/123-test", true, DownloadProxy},
	{"EdgeOne海外线路", "https://list.yearnstudio.cn", "https://list.yearnstudio.cn/d/123pan/pan", true, DownloadProxy},
	{"123直链", "https://vip.123pan.cn", "https://dl.yearnstudio.cn/1814376442/alist/%E5%88%86%E4%BA%AB%E7%94%A8", true, DownloadProxy},
	{"音乐 - 123云盘分享", "https://www.123865.com/s/km6bVv-cJ1W3", "/%F0%9F%8E%B5%E9%9F%B3%E4%B9%90%20%C2%B7%20Music%E2%80%8E%E2%80%8E%E2%80%8E/", false, ShareLink},
	{"游戏 - 123云盘分享", "https://www.123865.com/s/km6bVv-bJ1W3", "/%F0%9F%8E%AE%E6%B8%B8%E6%88%8F%20%C2%B7%20Game%E2%80%8E%E2%80%8E%E2%80%8E%E2%80%8E%E2%80%8E%E2%80%8E", false, ShareLink},
	{"软件 - 123云盘分享", "https://www.123865.com/s/km6bVv-7J1W3", "/%F0%9F%92%BE%E8%BD%AF%E4%BB%B6%20%C2%B7%20Software%E2%80%8E%E2%80%8E%E2%80%8E/%E5%85%B6%E4%BB%96", false, ShareLink},
	{"系统镜像 - 123云盘分享", "https://www.123865.com/s/km6bVv-SB1W3", "/%E7%B3%BB%E7%BB%9F%E9%95%9C%E5%83%8F%20%C2%B7%20System%20Images%E2%80%8E%E2%80%8E%E2%80%8E%E2%80%8E%E2%80%8E%E2%80%8E", false, ShareLink},
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
		var result []SiteInfo // 结果储存
		for _, v := range sites {
			if v.Showed { // 判断是否显示
				result = append(result, SiteInfo{Name: v.Name, URL: v.Url})
			}
		}
		c.JSON(http.StatusOK, result)
	})

	r.POST("/api/get", func(ctx *gin.Context) {
		// 获取请求体
		var reqBody struct {
			URL string `json:"url" binding:"required"`
		} // 使用json标签令gin框架解析JSON中的url参数
		if err := ctx.ShouldBindJSON(&reqBody); err != nil { // URL字段不存在
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		url := reqBody.URL // 获取到url
		if strings.HasPrefix(url, "http://") {
			url = strings.Replace(url, "http://", "https://", 1)
		}
		url = strings.Replace(url, "/p/", "/d/", 1) // 将/p/的下载路径格式转换为/d/以便于处理
		for _, v := range sites {
			if v.SiteType == DownloadProxy && strings.HasPrefix(url, v.BaseUrl) {
				// 匹配成功
				url = strings.Replace(url, v.BaseUrl, "", 1) // 替换开头的
				var result []string = []string{}             // 初始化结果储存
				for _, v := range sites {
					switch v.SiteType {
					case DownloadProxy:
						// 是下载代理
						result = append(result, v.BaseUrl+url) // 拼接
					case ShareLink:
						// 是分享链接
						if strings.HasPrefix(url, v.BaseUrl) {
							result = append(result, v.Url) // 直接返回即可
						}
					}
				}
				ctx.JSON(http.StatusOK, result) // 添加结果到返回体
				return                          // 跳出循环
			}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{})
	})

	r.Run(":27645")
}
