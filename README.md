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
  - **支持随机延迟执行** (Random Delay)，避免任务特征检测。
  - 任务执行日志持久化与结果回溯。
- **消息通知**：
  - 支持 SMTP 邮件通知。
  - 支持 Telegram Bot 消息推送。
  - 支持 **Bark** (iOS) 实时推送。
  - 可配置任务执行失败/成功时的即时告警。
- **安全认证**：
  - 独立的 Web 登录系统（JWT 认证）。
  - 支持 SSH 密码与密钥（Private Key）认证。
- **现代化设计**：全新的毛玻璃拟态 UI，响应式布局，支持夜间模式。
- **灵活的数据存储**：支持 SQLite（默认）及外部 MySQL/MariaDB 数据库。

## 快速部署

### 1. Docker 镜像启动 (推荐)

**默认模式 (SQLite)**：
直接启动，数据存储在挂载的 `data` 目录中。

```bash
docker run -d \
  -p 8888:8888 \
  -e USER=admin \
  -e PASS=admin123 \
  -e TZ=Asia/Shanghai \
  -v $(pwd)/data:/app/data \
  --name webssh \
  --restart always \
  ghcr.io/workerspages/webssh:external-database
```

**连接外部数据库 (MySQL/MariaDB)**：
通过环境变量配置连接外部数据库（需提前创建好数据库，例如 `webssh`）。

```bash
docker run -d \
  -p 8888:8888 \
  -e USER=admin \
  -e PASS=admin123 \
  -e TZ=Asia/Shanghai \
  -e DB_TYPE=mariadb \
  -e DB_HOST=192.168.1.100 \
  -e DB_PORT=3306 \
  -e DB_USER=root \
  -e DB_PASS=your_db_password \
  -e DB_NAME=webssh \
  -v $(pwd)/data:/app/data \
  --name webssh \
  --restart always \
  ghcr.io/workerspages/webssh:external-database
```

### 2. Docker Compose

```yaml
version: '3.8'
services:
  webssh:
    image: ghcr.io/workerspages/webssh:external-database
    container_name: webssh
    ports:
      - "8888:8888"
    environment:
      - USER=admin        # 初始管理员用户名
      - PASS=admin123     # 初始管理员密码
      - TZ=Asia/Shanghai  # 时区设置
      
      # 数据库配置 (可选，不填默认使用 SQLite)
      # - DB_TYPE=mariadb
      # - DB_HOST=192.168.1.100
      # - DB_PORT=3306
      # - DB_USER=root
      # - DB_PASS=password
      # - DB_NAME=webssh
    volumes:
      - ./data:/app/data  # 即使使用 MySQL，建议挂载 data 目录以保存日志文件等
    restart: unless-stopped
```

### 环境变量说明

| 变量名 | 说明 | 默认值 |
| :--- | :--- | :--- |
| `PORT` | Web服务端口 (仅二进制运行有效，Docker请用端口映射) | `8888` |
| `USER` | 初始管理员用户名 (仅首次初始化数据库时有效) | `admin` |
| `PASS` | 初始管理员密码 (仅首次初始化数据库时有效) | `admin123` |
| `TZ` | 系统时区 (影响定时任务触发时间) | `Asia/Shanghai` |
| **数据库配置** | | |
| `DB_TYPE` | 数据库类型，支持 `sqlite`, `mysql`, `mariadb` | `sqlite` |
| `DB_HOST` | 数据库地址 (IP 或域名) | `127.0.0.1` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASS` | 数据库密码 | - |
| `DB_NAME` | 数据库名称 (请先手动创建库) | `webssh` |

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
- **后端**: Go 1.24, Gin, GORM (SQLite/MySQL), Gorilla WebSocket, Robfig Cron
- **部署**: Docker, GitHub Actions

## 开源协议

MIT License
