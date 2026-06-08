package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api-query/internal/service"
	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	service     *service.QueryService
	accessToken string
}

func NewQueryHandler(service *service.QueryService, accessToken string) *QueryHandler {
	return &QueryHandler{
		service:     service,
		accessToken: accessToken,
	}
}

// AuthMiddleware Token 认证中间件
func (h *QueryHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// 支持 Bearer token 格式
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if token != h.accessToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid access token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUsageByDate 按日期查询使用情况
// GET /api/query/usage-by-date?date=2026-06-04&page=1&page_size=20
func (h *QueryHandler) GetUsageByDate(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date format, expected YYYY-MM-DD",
		})
		return
	}

	// 分页参数
	page := 0
	pageSize := 0

	// 只有当传了 page 或 page_size 参数时才启用分页
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 1000 {
			pageSize = parsed
		}
	}

	report, err := h.service.GetUsageByDate(c.Request.Context(), date, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := gin.H{
		"success": true,
		"data":    report,
	}

	// 只有启用分页时才返回分页信息
	if page > 0 && pageSize > 0 {
		response["pagination"] = gin.H{
			"page":      page,
			"page_size": pageSize,
		}
	}

	c.JSON(http.StatusOK, response)
}

// ListAPIKeys 列出所有 API Key
// GET /api/query/api-keys
func (h *QueryHandler) ListAPIKeys(c *gin.Context) {
	keys, err := h.service.ListAPIKeys(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 统计各类型数量
	stats := make(map[string]int)
	for _, key := range keys {
		stats[key.NameType]++
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"keys":       keys,
			"total":      len(keys),
			"statistics": stats,
		},
	})
}

// HealthCheck 健康检查
// GET /health
func (h *QueryHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	})
}
