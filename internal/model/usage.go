package model

import "time"

// UsageStats 使用统计数据
type UsageStats struct {
	Date                  string    `json:"date"`
	UserID                *int64    `json:"user_id,omitempty"`
	UserEmail             *string   `json:"user_email,omitempty"`
	Username              *string   `json:"username,omitempty"`
	APIKeyID              *int64    `json:"api_key_id,omitempty"`
	APIKeyName            *string   `json:"api_key_name,omitempty"`
	Model                 *string   `json:"model,omitempty"`
	RequestCount          int       `json:"request_count"`
	TotalInputTokens      int64     `json:"total_input_tokens"`
	TotalOutputTokens     int64     `json:"total_output_tokens"`
	TotalCacheCreation    int64     `json:"total_cache_creation_tokens"`
	TotalCacheRead        int64     `json:"total_cache_read_tokens"`
	TotalTokens           int64     `json:"total_tokens"`
	TotalCost             float64   `json:"total_cost"`
	ModelsUsed            []string     `json:"models_used,omitempty"`
	ModelStats            []ModelStat  `json:"model_stats,omitempty"`
}

// ModelStat 按模型的统计信息
type ModelStat struct {
	Model                 string  `json:"model"`
	RequestCount          int     `json:"request_count"`
	TotalInputTokens      int64   `json:"total_input_tokens"`
	TotalOutputTokens     int64   `json:"total_output_tokens"`
	TotalCacheCreation    int64   `json:"total_cache_creation_tokens"`
	TotalCacheRead        int64   `json:"total_cache_read_tokens"`
	TotalTokens           int64   `json:"total_tokens"`
	TotalCost             float64 `json:"total_cost"`
}

// DailyReport 日报告
type DailyReport struct {
	Date       string       `json:"date"`
	TotalStats UsageStats   `json:"total_stats"`
	ByUser     []UsageStats `json:"by_user,omitempty"`
	ByAPIKey   []UsageStats `json:"by_api_key,omitempty"`
	ByModel    []UsageStats `json:"by_model,omitempty"`
}

// APIKeyInfo API Key 信息
type APIKeyInfo struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	UserEmail    string     `json:"user_email"`
	Username     *string    `json:"username,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	LastUsedAt   *time.Time `json:"last_used_at,omitempty"`
	RequestCount int        `json:"request_count"`
	TotalCost    float64    `json:"total_cost"`
	NameType     string     `json:"name_type"`
	UsageLevel   string     `json:"usage_level"`
}
