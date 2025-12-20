# WebSSH

## 项目简介

WebSSH 是一个基于 **Go (Backend)** 和 **Vue2 (Frontend)** 构建的现代化 Web 端 SSH 连接与任务调度平台。它不仅提供了流畅的 Web 终端体验和 SFTP 文件管理功能，还内置了强大的定时任务调度系统和通知服务，是运维管理的得力助手。

> **核心特性**： 现代化 UI | SSH 终端 | SFTP 文件管理 | Cron 定时任务 (支持随机延迟) | 邮件/Telegram/Bark 通知

## 功能特性

- **Web 终端**：基于 Xterm.js 的高性能终端，支持自定义主题、字体、快捷连接及会话管理。
- **文件管理**：集成的 SFTP 面板，支持拖拽上传、下载、实时浏览服务器文件。
- **定时任务**：
  - 支持标准 Cron 表达式（精确到秒）。
  - 支持多步骤命令链（Command Chaining）。
  - 支持标准 Cron 表达式（精确到秒）。
  - 支持多步骤命令链（Command Chaining）。
  - **支持随机延迟执行** (Random Delay)，避免任务特征检测。
  - 任务执行日志持久化与结果回溯。
- **消息通知**：
  - 支持 SMTP 邮件通知。
  - 支持 SMTP 邮件通知。
  - 支持 Telegram Bot 消息推送。
  - 支持 **Bark** (iOS) 实时推送。
  - 可配置任务执行失败/成功时的即时告警。
- **安全认证**：
  - 独立的 Web 登录系统（JWT 认证）。
  - 支持 SSH 密码与密钥（Private Key）认证。
- **现代化设计**：全新的毛玻璃拟态 UI，响应式布局，支持夜间模式。

## 快速部署

### 1. Docker 镜像启动 (推荐)

直接使用 GitHub Container Registry 托管的最新镜像：

```bash
docker run -d \
  -p 8888:8888 \
  -e USER=admin        # 初始管理员用户名
  -e PASS=admin123     # 初始管理员密码
  -e TZ=Asia/Shanghai  # 设置时区 (重要: 影响定时任务触发时间)
  --name webssh \
  --restart always \
  ghcr.io/workerspages/webssh:latest
```

### 2. Docker Compose

```yaml
version: '3.8'
services:
  webssh:
    image: ghcr.io/workerspages/webssh:latest
    container_name: webssh
    ports:
      - "8888:8888"
    environment:
      - USER=admin        # 自定义用户名
      - PASS=admin123     # 自定义密码
      - TZ=Asia/Shanghai  # 设置时区，这对 Cron 定时任务非常重要
    volumes:
      - ./data:/app/data  # 挂载数据目录以持久化数据库(webssh.db)
    restart: unless-stopped
```

## 源码构建

如果您需要二次开发或自行构建：

### 前置要求

- **Go**: 1.24+
- **Node.js**: 18+

### 构建步骤

1. **克隆仓库**

   ```bash
   git clone https://github.com/workerspages/webssh.git
   cd webssh
   ```

2. **构建前端**

   ```bash
   cd frontend
   npm install
   npm run build
   # 构建产物将自动输出到 ../public 目录
   ```

3. **构建后端**

   ```bash
   cd ..
   # 启用 CGO_ENABLED=0 以确保静态链接兼容性
   set CGO_ENABLED=0
   go build -ldflags "-s -w" -o webssh main.go
   ```

4. **运行**

   ```bash
   ./webssh
   ```

   访问 `http://localhost:8888` 即可。

## 技术栈

- **前端**: Vue 2.7, Element UI, Xterm.js, Axios
- **后端**: Go 1.24, Gin, GORM (SQLite), Gorilla WebSocket, Robfig Cron
- **部署**: Docker, GitHub Actions

## 开源协议

MIT License
