# PaiPai - Instant Messaging IM based on GoZero microservices and AI large model applications

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Microservices-Architecture-6BA539?style=for-the-badge" alt="Architecture">
  <img src="https://img.shields.io/badge/Cloud_Native-3371E3?style=for-the-badge&logo=Docker&logoColor=white" alt="Cloud Native">
  <img src="https://img.shields.io/badge/AI_Enhanced-FF6F00?style=for-the-badge&logo=ai&logoColor=white" alt="AI Enhanced">
</div>


## 📋 项目简介
   基于微服务架构设计的现代化社交服务平台，专注于提供高质量的即时通信、丰富的社交互动以及智能AI增强体验。项目采用云原生设计理念，结合主流技术栈，同时进一步融合AI大模型应用，旨在构建安全、可靠、高性能且易于扩展的社交服务解决方案。

---

## 🚀 整体架构
<img width="2186" height="975" alt="image" src="https://github.com/user-attachments/assets/a2feb290-baf0-4490-8ccd-ee48b6d094d4" />

---

## 🛠️ 核心技术栈
| 类别 | 技术组件 | 推荐版本 | 使用说明 |
|------|----------|----------|----------|
| **编程语言** | Golang | 1.21+ | 后端主要开发语言 |
| **微服务框架** | GoZero | 1.8.5 | 微服务开发框架 |
| **数据库** | MySQL | 8.0+ | 关系型数据存储 |
|  | Redis | 7.0+ | 缓存、会话管理 |
|  | MongoDB | 6.0+ | 文档型数据存储 |
| **消息队列** | Kafka | 3.5+ | 高吞吐量消息处理 |
|  | RabbitMQ | 3.11+ | 复杂路由消息队列 |
| **服务发现** | ETCD | 3.5+ | 服务发现与配置管理 |
| **API网关** | Apisix | 3.7+ | API网关与流量管理 |
| **监控观测** | Prometheus | latest | 指标收集与监控 |
|  | Grafana | latest | 数据可视化仪表板 |
|  | Jaeger | latest | 分布式链路追踪 |
| **日志系统** | Elasticsearch | latest | 日志存储与检索 |
|  | Logstash | latest | 日志收集与处理 |
|  | Kibana | latest | 日志可视化分析 |
| **容器编排** | Docker | latest | 容器运行时 |
|  | Kubernetes | latest | 容器编排管理 |
| **AI应用/集成** | DeepSeek | - | 智能代码审查 |
|  | Eino | latest | 大模型应用框架 |

## 🏗️架构层次
| 架构层 | 技术组件 | 功能描述 |
|--------|----------|----------|
| **接入层** | Apisix, Docker, Kubernetes | 流量接入、负载均衡、容器化管理 |
| **服务层** | Golang, GoZero, ETCD | 业务逻辑处理、服务发现与注册 |
| **数据层** | MySQL, Redis, MongoDB | 数据存储、缓存、持久化 |
| **消息层** | Kafka, RabbitMQ | 异步消息处理、系统解耦 |
| **观测层** | Prometheus, Grafana, Jaeger, ELK | 系统监控、日志分析、链路追踪 |
| **智能层** | DeepSeek, Eino | AI大模型集成、智能业务处理 |

## 🔄系统数据流
| 数据流向 | 技术组件 | 协议/接口 |
|----------|----------|----------|
| **客户端请求** | Apisix → GoZero微服务 | HTTP/HTTPS/WebSocket |
| **服务间通信** | 微服务间调用 | gRPC/HTTP + ETCD服务发现 |
| **消息处理** | Kafka/RabbitMQ → 业务处理 | AMQP/MQTT/自定义 |
| **数据持久化** | → MySQL/Redis/MongoDB | SQL/NoSQL接口 |
| **监控数据** | → Prometheus/ELK | 指标采集/日志收集 |
| **AI大模型集成** | → DeepSeek/Eino | REST API/gRPC |

---

## 📂 项目结构  

```plaintext
PaiPai/
├── apps/            # 核心业务
│   ├── im/          # 即时通信服务
│   ├── social/      # 社交服务
│   ├── task/        # 任务消息队列
│   ├── user/        # 用户服务
│   └── eino_chat/   # AI服务
├── components/      # 基础设施组件
│   ├── apisix/           # API网关配置
│   ├── apisix-dashboard/ # API网关管理后台
│   ├── filebeat/         # 日志收集
│   ├── grafana/          # 监控可视化
│   ├── kibana/           # 日志分析
│   ├── logstash/         # 日志处理
│   ├── prometheus/       # 监控系统
│   └── sail/             # 服务编排
├── deploy/          # 部署配置文件
│   ├── cicd/             # CI/CD、代码审查
│   ├── dockerfile/       # Docker构建
│   ├── makefile/         # 构建脚本
│   ├── script/           # 启动测试脚本
│   └── sql/              # 数据库初始化
├── pkg/             # 公共工具包
├── test/            # 测试代码和示例
├── go.mod           
├── go.sum           
├── docker-compose.yaml   # 容器编排配置
├── PORT_MAP.md      # 端口映射文档
└── README.md        # 项目说明文档
```

---

## 🌟 功能特性
- **微服务三高架构**
  - 限流、熔断、降级机制
  - 高可用、高性能、高扩展性设计
- **高效IM通信引擎**
  - WebSocket + gRPC 高效通信
  - 智能路由节点消息中转优化
  - 保障高并发、低延迟
- **完善消息收发体系**
  - 全场景覆盖，支持文本/图片/语音/视频/文件/位置 
  - 消息漫游，云端历史消息存储
  - 消息安全加密，隐私保护
- **深度融合 AI + 云原生**
  - AI 智能代码审计
  - AI 大模型应用融合
  - LLM 驱动智能回复诉求
  - 消息语义分析与风险识别
  - 企业级知识库高效构建
  - AIOps 异常流量自动识别与防控
- **自动化容器化便捷部署**
  - Docker 全容器化
  - 智能编排部署
- **全链路监测保障**
  - Metrics/Logging/Tracing 三位一体
  - Prometheus + Grafana + Jaeger 全面监测
 
---

## 🤝 贡献指南

欢迎对项目进行贡献！感谢！

1. **Fork 仓库**并克隆到本地
2. **创建分支**进行开发（`git checkout -b feature/your-feature`）
3. **提交代码**并确保通过测试
4. **创建 Pull Request** 描述您的更改
5. 等待**代码审查**并根据反馈进行修改

---

### <div align="center"> <strong>✨ 持续更新完善中... ✨</strong> </div>
