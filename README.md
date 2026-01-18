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
      # 管理员账号配置
      - USER=admin
      - PASS=admin123
      - TZ=Asia/Shanghai
      
      # ----------- 外部数据库配置 -----------
      - DB_TYPE=mariadb           # 或 mysql
      - DB_HOST=192.168.1.100     # 【重点】修改为您外部数据库的真实 IP 或域名
      - DB_PORT=3306              # 数据库端口
      - DB_USER=root              # 数据库用户名
      - DB_PASS=your_db_password  # 数据库密码
      - DB_NAME=webssh            # 数据库库名 (请确保该库已存在)
      # ------------------------------------
      
      # ----------- Tailscale 配置 (可选) -----------
      # - TS_AUTHKEY=tskey-auth-xxxxx  # Tailscale Auth Key (设置后自动启用)
      # - TS_HOSTNAME=webssh           # Tailscale 网络中的主机名

      # 果你不加这一行，WebSSH 启动时会默认尝试以“个人设备”身份注册，但你的 Key 只有“服务器 Tag”的权限，可能会导致权限不匹配而报错。加上这行就万无一失了。
      # - TS_EXTRA_ARGS=--advertise-tags=tag:webssh
      # ---------------------------------------------
      
    volumes:
      - ./data:/app/data          # 即使使用 MySQL，建议挂载 data 目录以保存日志文件等
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
| **Tailscale 配置** | | |
| `TS_AUTHKEY` | Tailscale Auth Key (设置后自动启用 Tailscale) | - |
| `TS_HOSTNAME` | Tailscale 网络中的主机名 | `webssh` |
| `TS_EXTRA_ARGS` | “服务器 Tag”的权限 | `--advertise-tags=tag:webssh` |

### IPv6 与 Tailscale 支持

WebSSH 原生支持连接 IPv6 服务器。对于部署在只有 IPv6 出站能力的 PaaS 平台（如 Koyeb），可通过集成 Tailscale 获得 IPv4 出站能力。

**启用 Tailscale：**

1. 在 [Tailscale Admin Console](https://login.tailscale.com/admin/settings/keys) 生成 Auth Key
2. 设置 `TS_AUTHKEY` 环境变量启动容器

```bash
docker run -d \
  -p 8888:8888 \
  -e TS_AUTHKEY=tskey-auth-xxxxx \
  -e USER=admin \
  -e PASS=admin123 \
  --name webssh \
  ghcr.io/workerspages/webssh:latest
```

启用后，WebSSH 可通过 Tailscale 网络连接到其他 Tailscale 节点（使用 `100.x.y.z` 地址）。

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
