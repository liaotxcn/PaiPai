# PaiPai - Instant Messaging IM based on GoZero microservices and AI large model applications

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Microservices-Architecture-6BA539?style=for-the-badge" alt="Architecture">
  <img src="https://img.shields.io/badge/Cloud_Native-3371E3?style=for-the-badge&logo=Docker&logoColor=white" alt="Cloud Native">
  <img src="https://img.shields.io/badge/AI_Enhanced-FF6F00?style=for-the-badge&logo=ai&logoColor=white" alt="AI Enhanced">

  **Language Options:** [中文](README.zh-CN.md) | [English](README.md) 
</div>


## 📋 Project Overview
   A modern social service platform based on microservices architecture, focusing on providing high-quality instant messaging, rich social interactions, and intelligent AI-enhanced experiences. The project adopts cloud-native design principles, combines mainstream technology stacks, and further integrates AI large model applications, aiming to build a secure, reliable, high-performance, and easily scalable social service solution.

---

## 🚀 Overall Architecture
<img width="2186" height="975" alt="image" src="https://github.com/user-attachments/assets/a2feb290-baf0-4490-8ccd-ee48b6d094d4" />

---

## 🛠️ Core Technology Stack
| Category | Technology Components | Recommended Version | Usage Description |
|----------|----------------------|---------------------|-------------------|
| **Programming Language** | Golang | 1.21+ | Primary backend development language |
| **Microservices Framework** | GoZero | 1.8.5 | Microservices development framework |
| **Database** | MySQL | 8.0+ | Relational data storage |
|  | Redis | 7.0+ | Caching, session management |
|  | MongoDB | 6.0+ | Document-based data storage |
| **Message Queue** | Kafka | 3.5+ | High-throughput message processing |
|  | RabbitMQ | 3.11+ | Complex routing message queue |
| **Service Discovery** | ETCD | 3.5+ | Service discovery and configuration management |
| **API Gateway** | Apisix | 3.7+ | API gateway and traffic management |
| **Monitoring & Observability** | Prometheus | latest | Metrics collection and monitoring |
|  | Grafana | latest | Data visualization dashboard |
|  | Jaeger | latest | Distributed tracing |
| **Logging System** | Elasticsearch | latest | Log storage and retrieval |
|  | Logstash | latest | Log collection and processing |
|  | Kibana | latest | Log visualization and analysis |
| **Container Orchestration** | Docker | latest | Container runtime |
|  | Kubernetes | latest | Container orchestration management |
| **AI Applications/Integration** | DeepSeek | - | Intelligent code review |
|  | Eino | latest | Large model application framework |

## 🏗️Architecture Layers
| Architecture Layer | Technology Components | Functional Description |
|--------------------|----------------------|------------------------|
| **Access Layer** | Apisix, Docker, Kubernetes | Traffic access, load balancing, container management |
| **Service Layer** | Golang, GoZero, ETCD | Business logic processing, service discovery and registration |
| **Data Layer** | MySQL, Redis, MongoDB | Data storage, caching, persistence |
| **Message Layer** | Kafka, RabbitMQ | Asynchronous message processing, system decoupling |
| **Observability Layer** | Prometheus, Grafana, Jaeger, ELK | System monitoring, log analysis, distributed tracing |
| **Intelligence Layer** | DeepSeek, Eino | AI large model integration, intelligent business processing |

## 🔄System Data Flow
| Data Flow | Technology Components | Protocol/Interface |
|-----------|----------------------|-------------------|
| **Client Requests** | Apisix → GoZero Microservices | HTTP/HTTPS/WebSocket |
| **Inter-service Communication** | Microservices inter-calls | gRPC/HTTP + ETCD service discovery |
| **Message Processing** | Kafka/RabbitMQ → Business processing | AMQP/MQTT/Custom |
| **Data Persistence** | → MySQL/Redis/MongoDB | SQL/NoSQL interfaces |
| **Monitoring Data** | → Prometheus/ELK | Metrics collection/Log collection |
| **AI Large Model Integration** | → DeepSeek/Eino | REST API/gRPC |

---

## 📂 Project Structure

```plaintext
PaiPai/
├── apps/            # Core business services
│   ├── im/          # Instant Messaging service
│   ├── social/      # Social service
│   ├── task/        # Task message queue
│   ├── user/        # User service
│   └── eino_chat/   # AI service
├── components/      # Infrastructure components
│   ├── apisix/           # API gateway configuration
│   ├── apisix-dashboard/ # API gateway admin dashboard
│   ├── filebeat/         # Log collection
│   ├── grafana/          # Monitoring visualization
│   ├── kibana/           # Log analysis
│   ├── logstash/         # Log processing
│   ├── prometheus/       # Monitoring system
│   └── sail/             # Service orchestration
├── deploy/          # Deployment configuration files
│   ├── cicd/             # CI/CD, code review
│   ├── dockerfile/       # Docker builds
│   ├── makefile/         # Build scripts
│   ├── script/           # Startup and test scripts
│   └── sql/              # Database initialization
├── pkg/             # Common utility packages
├── test/            # Test code and examples
├── go.mod           
├── go.sum           
├── docker-compose.yaml   # Container orchestration configuration
├── PORT_MAP.md      # Port mapping documentation
└── README.md        # Project documentation
```

---

## 🌟 Key Features
- **High Availability Microservices Architecture**
  - Rate limiting, circuit breaking, and fallback mechanisms
  - High availability, high performance, and high scalability design
- **Efficient IM Communication Engine**
  - WebSocket + gRPC high-efficiency communication
  - Intelligent routing node message relay optimization
  - Ensures high concurrency and low latency
- **Comprehensive Message System**
  - Full scenario coverage: text/image/voice/video/file/location support
  - Message roaming with cloud-based historical message storage
  - Message security encryption and privacy protection
- **Deep Integration of AI + Cloud Native**
  - AI-powered intelligent code audit
  - AI large model application integration
  - LLM-driven intelligent response capabilities
  - Message semantic analysis and risk identification
  - Efficient enterprise knowledge base construction
  - AIOps for automatic anomaly traffic identification and prevention
- **Automated Containerized Deployment**
  - Full Docker containerization
  - Intelligent orchestration deployment
- **Full-Link Monitoring and Assurance**
  - Metrics/Logging/Tracing integrated monitoring
  - Comprehensive monitoring with Prometheus + Grafana + Jaeger

---

## 🤝 Contributing Guidelines

We welcome contributions to the project! Thank you for your interest!

1. **Fork the repository** and clone it locally
2. **Create a branch** for your development (`git checkout -b feature/your-feature`)
3. **Commit your code** and ensure all tests pass
4. **Create a Pull Request** describing your changes
5. Wait for **code review** and make modifications based on feedback

---

### <div align="center"> <strong>✨ Continuously Improving and Updating... ✨</strong> </div>
