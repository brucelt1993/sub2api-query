package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wei-Shaw/sub2api-query/internal/config"
	"github.com/Wei-Shaw/sub2api-query/internal/database"
	"github.com/Wei-Shaw/sub2api-query/internal/handler"
	"github.com/Wei-Shaw/sub2api-query/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✓ Database connected successfully")

	// 初始化服务
	queryService := service.NewQueryService(db)
	queryHandler := handler.NewQueryHandler(queryService, cfg.Auth.AccessToken)

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 注册路由
	api := router.Group("/api")
	{
		query := api.Group("/query")
		query.Use(queryHandler.AuthMiddleware()) // 添加认证中间件
		{
			query.GET("/usage-by-date", queryHandler.GetUsageByDate)
			query.GET("/api-keys", queryHandler.ListAPIKeys)
		}
	}

	router.GET("/health", queryHandler.HealthCheck)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("🚀 Query service starting on %s", cfg.Server.Address())
		log.Println("📊 Endpoints:")
		log.Println("   GET /api/query/usage-by-date?date=2026-06-04")
		log.Println("   GET /api/query/api-keys")
		log.Println("   GET /health")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server stopped gracefully")
}
