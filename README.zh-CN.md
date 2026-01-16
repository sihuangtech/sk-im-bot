# SK-IM-Bot

一个现代化的多平台聊天机器人系统，支持 QQ 个人号和 Discord，并集成了大语言模型（LLM）对话功能。配套提供极其美观的 Web 管理后台。

## 🌟 核心功能

-   **多平台支持**: 接入 QQ（通过 OneBot 11 协议）和 Discord 服务器。
-   **智能对话**: 集成 OpenAI 等 LLM 服务，支持上下文理解。
-   **多媒体支持**: 支持接收和发送图片消息。
-   **Web 管理后台**:
    -   **实时监控**: 通过 WebSocket 实时展示所有平台的聊天记录。
    -   **配置管理**: 在线修改 LLM 参数及平台 Token。
    -   **系统概览**: 查看系统运行状态及统计数据。
-   **防封号机制**: 模拟真人行为，包含随机回复延迟（500ms - 3s）和消息频率限制。
-   **现代 UI**: 基于 React 18 和 Ant Design，采用毛玻璃（Glassmorphism）设计风格。

## 🛠️ 技术栈

### 后端 (Go)
-   **Framework**: Gin
-   **ORM**: GORM (PostgreSQL)
-   **Protocol**: OneBot 11 (QQ), Discordgo
-   **Config**: Viper
-   **Logs**: Zap
-   **Auth**: JWT

### 前端 (React)
-   **Runtime**: React 18 + TypeScript
-   **Build Tool**: Vite
-   **UI Library**: Ant Design 5
-   **State**: Zustand
-   **Icons**: Lucide React

## 🚀 快速开始

### 1. 环境準備
-   Docker & Docker Compose
-   Go 1.25+ (本地运行)
-   Node.js 18+ (本地运行)

### 2. 通过 Docker 一键部署 (推荐)
1.  复制 `.env.example` 为 `.env` 并填写你的 API 密钥：
    ```bash
    cp .env.example .env
    ```
2.  启动服务：
    ```bash
    docker-compose up -d
    ```
3.  访问管理界面：`http://localhost:5173` (默认账号: `admin` / `admin`)

### 3. 本地开发
**后端:**
```bash
cd backend
go mod tidy
go run cmd/main.go
```
**前端:**
```bash
cd frontend
npm install
npm run dev
```

## 📂 项目结构
```
.
├── backend/            # Go 后端服务
│   ├── cmd/           # 入口程序
│   ├── internal/      # 核心逻辑 (API, Bot 适配器, 模型)
│   └── pkg/           # 内部公共工具
├── frontend/           # React 前端服务
│   └── src/           # 源代码 (页面, 组件, 状态管理)
└── docker-compose.yml  # 编排配置
```

## ⚖️ 许可
MIT License
