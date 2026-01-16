# SK-IM-Bot

A modern, multi-platform chatbot system supporting QQ and Discord, integrated with Large Language Models (LLM) and a stunning Web dashboard.

## 🌟 Key Features

-   **Multi-Platform Support**: Seamless integration with QQ (via OneBot 11) and Discord servers.
-   **Intelligent Conversation**: Powered by LLMs (OpenAI, etc.) with context-aware chat management.
-   **Multimedia Support**: Capable of receiving and sending image messages.
-   **Web Management Dashboard**:
    -   **Real-time Monitoring**: Live chat stream via WebSockets.
    -   **Dynamic Configuration**: Update configurations and tokens on the fly.
    -   **System Statistics**: Overview of system status and activity.
-   **Anti-Ban Mechanisms**: Human-like behavior simulation with random delays (500ms - 3s) and rate limiting.
-   **Modern Aesthetics**: Premium dark theme UI with Glassmorphism using React 18 and Ant Design.

## 🛠️ Technology Stack

### Backend (Go)
-   **Framework**: Gin
-   **ORM**: GORM (PostgreSQL)
-   **Adapters**: OneBot 11 (QQ), Discordgo
-   **Utility**: Viper (Config), Zap (Logging), JWT (Auth)

### Frontend (React)
-   **Framework**: React 18 + TypeScript
-   **Build Tool**: Vite
-   **UI System**: Ant Design 5 + Lucide Icons
-   **State Management**: Zustand

## 🚀 Quick Start

### 1. Prerequisites
-   Docker & Docker Compose
-   Go 1.25+ (for local development)
-   Node.js 18+ (for local development)

### 2. Deployment with Docker (Recommended)
1.  Clone the `.env.example` to `.env` and fill in your keys:
    ```bash
    cp .env.example .env
    ```
2.  Launch the stack:
    ```bash
    docker-compose up -d
    ```
3.  Access the Dashboard: `http://localhost:5173` (Default Admin: `admin` / `admin`)

### 3. Local Development
**Backend:**
```bash
cd backend
go mod tidy
go run cmd/main.go
```
**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## 📂 Project Structure
```
.
├── backend/            # Go service
│   ├── cmd/           # Entrance
│   ├── internal/      # Core logic (API, Bot Adapters, etc.)
│   └── pkg/           # Shared utilities
├── frontend/           # React application
│   └── src/           # Source code (Pages, Components, Stores)
└── docker-compose.yml  # Infrastructure as Code
```

## ⚖️ License
MIT License
