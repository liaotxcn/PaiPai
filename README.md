# PaiPai - Instant Messaging IM based on GoZero microservices and AI large model applications

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Microservices-Architecture-6BA539?style=for-the-badge" alt="Architecture">
  <img src="https://img.shields.io/badge/Cloud_Native-3371E3?style=for-the-badge&logo=Docker&logoColor=white" alt="Cloud Native">
  <img src="https://img.shields.io/badge/AI_Enhanced-FF6F00?style=for-the-badge&logo=ai&logoColor=white" alt="AI Enhanced">
</div>

---

## 🚀 整体架构
<img width="2186" height="975" alt="image" src="https://github.com/user-attachments/assets/a2feb290-baf0-4490-8ccd-ee48b6d094d4" />

---

## 📂 项目结构  

```plaintext
PaiPai/
├── apps/            # Service
│   ├── im/          # 即时通信服务
│   ├── social/      # 社交服务
│   ├── task/        # 事务服务
│   ├── user/        # 用户服务
│   └── eino_chat/   # EinoService
├── components/      # API网关
├── deploy/          # Docker部署&&运行脚本
├── pkg/             # 工具包
├── test/            # 实例测试等
├── go.mod               
├── docker-compose.yaml   
└── Makefile              
```

---

## 🌟 功能特性
- **微服务三高架构**
  - 限流、熔断、降级 
  - 高可用、高性能、高扩展
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

### 功能逐步完善中...







