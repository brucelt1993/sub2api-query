# Sub2API Query Service

Sub2API 使用统计查询服务，提供按日期、API Key、模型等多维度的查询功能。

## ✨ 功能特性

- ✅ 按日期查询使用统计
- ✅ 完整的 Token 统计（input + output + cache_creation + cache_read = total）
- ✅ 按 API Key 分组统计（含每个模型的详细数据）
- ✅ 智能过滤（自动排除测试 Key）
- ✅ 支持分页查询
- ✅ Access Token 认证保护
- ✅ Docker 一键部署

## 🚀 快速部署

### 方式 1: Docker Compose（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/brucelt1993/sub2api-query.git
cd sub2api-query

# 2. 复制配置文件模板
cp config.yaml.example config.yaml

# 3. 修改配置
nano config.yaml
# 修改 database 连接信息和 auth.access_token

# 4. 启动服务
docker-compose up -d

# 5. 查看日志
docker-compose logs -f
```

### 方式 2: 手动编译

```bash
# 1. 编译
go build -o query-service cmd/server/main.go

# 2. 启动
./query-service
```

## 📡 API 使用

### 认证

所有查询接口都需要在请求头中携带 Access Token：

```bash
Authorization: Bearer your-secret-token
```

### 示例

```bash
# 健康检查（无需认证）
curl http://localhost:8086/health

# 查询今天的数据（需要认证）
curl -H "Authorization: Bearer your-token" \
     "http://localhost:8086/api/query/usage-by-date?date=2026-06-08"

# 分页查询
curl -H "Authorization: Bearer your-token" \
     "http://localhost:8086/api/query/usage-by-date?date=2026-06-08&page=1&page_size=50"
```

## 🌐 Nginx 配置

参考 `nginx.conf` 文件配置反向代理，将服务暴露到公网。

## 📖 完整文档

详细部署指南请查看 [DEPLOY.md](./DEPLOY.md)

## 🔐 安全提示

- ⚠️ 请修改 `config.yaml` 中的 `auth.access_token`
- ⚠️ 生产环境务必使用 HTTPS
- ⚠️ 定期更换 Access Token

## 📊 返回数据示例

```json
{
  "success": true,
  "data": {
    "total_stats": {
      "request_count": 5000,
      "total_tokens": 1400000,
      "total_cost": 125.50
    },
    "by_api_key": [
      {
        "api_key_name": "张三",
        "request_count": 450,
        "total_tokens": 265000,
        "total_cost": 5.20,
        "model_stats": [
          {
            "model": "claude-3-sonnet",
            "request_count": 300,
            "total_tokens": 197000,
            "total_cost": 3.80
          }
        ]
      }
    ]
  }
}
```

## 🛠️ 技术栈

- Go 1.23+
- Gin Web Framework
- PostgreSQL
- Docker
