# PaiPai 端口映射说明文档

## 1. 概述
本文档详细说明 PaiPai 项目中各服务的端口映射配置、分配原则及特殊说明，旨在提高项目清晰度和易用性，帮助团队成员快速了解和使用各服务端口。

## 2. 端口分配原则

### 2.1 基本规则
- **唯一性**：确保所有服务的主机端口（宿主机端口）不冲突
- **高位端口**：统一使用 1024 以上的高位端口，避免与系统保留端口冲突
- **内外一致**：多数服务保持容器内外端口一致，便于维护和记忆
- **类型区分**：RPC 服务与 API 服务使用不同范围的端口，体现服务类型区分

### 2.2 端口范围划分
- API 服务：8880-8900 范围内
- RPC 服务：10000-10010 范围内
- WebSocket 服务：10090 及以上

## 3. 各服务端口映射表

| 服务名称 | 测试脚本 | 容器端口 | 主机端口 | 服务类型 | 健康检查方式 |
|---------|---------|---------|---------|---------|------------|
| User API | user-api-test.sh | 8888 | 8888 | HTTP API | HTTP 请求检查 `/metrics` |
| User API | user-api-test.sh | 7888 | 7888 | HTTP API | HTTP 请求检查 |
| User RPC | user-rpc-test.sh | 10000 | 10000 | gRPC | TCP 端口连接测试 |
| Social API | social-api-test.sh | 8888 | 8891 | HTTP API | HTTP 请求检查 `/metrics` |
| Social RPC | social-rpc-test.sh | 8888 | 8890 | gRPC | TCP 端口连接测试 |
| IM API | im-api-test.sh | 8882 | 8882 | HTTP API | HTTP 请求检查 |
| IM RPC | im-rpc-test.sh | 10002 | 10002 | gRPC | TCP 端口连接测试 |
| IM WebSocket | im-ws-test.sh | 10090 | 10090 | WebSocket | 连接可用性测试 |
| Task MQ | task-mq-test.sh | - | - | 消息队列 | 内部服务 |

## 4. 特殊端口说明

### 4.1 User API 双端口配置
User API 服务同时映射两个端口：
- **8888**: 主要 API 接口服务
- **7888**: 辅助功能或内部测试接口服务

## 5. 健康检查配置

所有服务均已配置 Docker 健康检查机制，确保服务稳定性：

### 5.1 HTTP API 服务（如 User API、Social API、IM API）
```bash
--health-cmd="wget --no-verbose --tries=1 --spider http://localhost:8888/metrics || exit 1"
--health-interval=30s
--health-timeout=5s
--health-retries=3
```

### 5.2 gRPC 服务（如 User RPC、Social RPC、IM RPC）
```bash
--health-cmd="timeout 5s bash -c '</dev/tcp/localhost/8888' || exit 1"
--health-interval=30s
--health-timeout=5s
--health-retries=3
```

### 5.3 WebSocket 服务（如 IM WebSocket）
```bash
--health-cmd="timeout 5s bash -c '</dev/tcp/localhost/10090' || exit 1"
--health-interval=30s
--health-timeout=5s
--health-retries=3
```

## 6. 端口使用建议

### 6.1 开发环境
- 本地开发时，请确保不占用上述端口，避免端口冲突
- 如果需要调整端口，请同步更新对应测试脚本和本文档

### 6.2 部署环境
- 生产环境建议使用反向代理（APISIX/NGINX）统一管理端口和流量
- 考虑使用端口安全组策略，限制外部访问权限

## 7. 更新记录
**最后更新时间：2025-09-24**