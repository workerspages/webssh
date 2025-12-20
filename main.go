package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"
	"webssh/controller"
	"webssh/model"
	"webssh/service"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//go:embed public/*
var f embed.FS

var (
	port = flag.Int("p", 8888, "服务运行端口")
	// -a 参数现在作为初始化管理员密码的快捷方式，或者依旧用于简单的基本认证（为了兼容旧习惯，我们这里保留基本认证用于ws，API使用JWT）
	authInfo = flag.String("a", "", "user:pass")
	timeout  int
)

func init() {
	flag.IntVar(&timeout, "t", 120, "ssh连接超时时间(min)")
	flag.Parse()
}

// JWTAuthMiddleware 中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"msg": "Unauthorized"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &service.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return service.SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"msg": "Invalid Token"})
			return
		}
		c.Next()
	}
}

func main() {
	// 1. 初始化数据库
	model.InitDB()
	// 2. 启动定时任务调度
	service.StartCron()

	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(gzip.Gzip(gzip.DefaultCompression))

	// API 路由
	api := server.Group("/api")
	{
		// 登录
		api.POST("/login", controller.LoginHandler)

		// 需认证的接口
		auth := api.Group("/", JWTAuthMiddleware())
		{
			auth.GET("/check", func(c *gin.Context) {
				responseBody := controller.CheckSSH(c)
				c.JSON(200, responseBody)
			})
			
			// 文件管理
			file := auth.Group("/file")
			{
				file.GET("/list", func(c *gin.Context) { c.JSON(200, controller.FileList(c)) })
				file.GET("/download", func(c *gin.Context) { controller.DownloadFile(c) })
				file.POST("/upload", func(c *gin.Context) { c.JSON(200, controller.UploadFile(c)) })
			}

			// 任务管理
			cron := auth.Group("/cron")
			{
				cron.GET("/list", controller.GetCronJobs)
				cron.POST("/save", controller.AddCronJob)
				cron.POST("/delete/:id", controller.DeleteCronJob)
				cron.POST("/run/:id", controller.RunCronJobManually)
			}

			// 通知配置
			notify := auth.Group("/notify")
			{
				notify.GET("/config", controller.GetNotifyConfig)
				notify.POST("/save", controller.SaveNotifyConfig)
				notify.POST("/test", controller.TestNotify)
			}
		}
	}

	// WebSocket (无需JWT头，通过 Query 参数传递或者简单的 BasicAuth，这里沿用原逻辑或做简单修改)
	// 为了兼容前端 xterm 的连接方式，这里保持原样，或者在 Query 中校验 token
	server.GET("/term", func(c *gin.Context) {
		controller.TermWs(c, time.Duration(timeout)*time.Minute)
	})
	server.GET("/file/progress", func(c *gin.Context) {
		controller.UploadProgressWs(c)
	})

	// 静态文件托管
	staticFS, _ := fs.Sub(f, "public/static")
	server.StaticFS("/static", http.FS(staticFS))

	server.NoRoute(func(c *gin.Context) {
		indexHTML, err := f.ReadFile("public/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "index.html not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	fmt.Printf("WebSSH Started on :%d\n", *port)
	server.Run(fmt.Sprintf(":%d", *port))
}
