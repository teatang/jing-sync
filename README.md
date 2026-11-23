# Jing-Sync

[English](README.me) | [简体中文](README.md)

Jing-Sync 是一个功能强大、现代化的全栈文件同步与管理 Web 应用程序。它支持通过 OpenList / AList 等存储引擎进行跨目录、跨网盘的高效文件同步、定时调度、镜像对齐以及详细的同步日志审计。

---

## 🌟 核心功能特性

- **多存储引擎管理**：深度集成 OpenList / AList 实例，支持统一管理多个存储引擎的访问地址与鉴权 Token，并提供一键连接测试功能。
- **直观的路径浏览与选择**：
  - 支持手动输入或粘贴源目录与目标目录路径。
  - 内置路径浏览弹窗（Path Browser），支持可视化逐级展开、快捷跳转与目录拾取。
- **多元同步策略**：
  - **仅新增 (增量同步)**：仅复制目标端缺失或有变更的文件，安全高效，不破坏目标端已有数据。
  - **完全同步 (镜像同步)**：源端与目标端状态精确镜像对齐，源端删除的文件会在目标端自动同步删除。
  - **移动文件 (剪切同步)**：文件成功传输至目标端后自动从源端清除。
- **细粒度文件过滤与控制**：
  - 支持通配符排除规则（Glob Patterns，例如 `*.tmp, *.bak, .git/**`）。
  - 支持文件大小上下限阈值过滤（最小 / 最大字节控制）。
  - 内置源路径与目标路径包含重叠校验，从逻辑层面杜绝死循环传输。
- **高可用定时调度**：
  - 支持快捷分钟级间隔设定或标准 5 段式 Cron 表达式（如 `0 2 * * *`）。
  - 全局动态 Cron 调度管理器（CronManager）：在作业新增、修改、启停或删除时自动同步调度状态，无需重启服务。
- **实时审计与日志终端**：
  - 详尽记录每一次同步任务的执行状态、发现文件数、新增数、删除数、跳过数及失败错误信息。
  - 抽屉式终端风格日志详情查看器，支持高亮显示与排错审计。
- **安全与权限体系**：
  - 采用基于 JWT 的无状态安全认证。
  - 首次运行自动生成强密码管理员账号并在终端高亮提示。
  - 用户管理模块支持用户账号维护、状态控制与密码脱敏展示。
  - 完善的国际化多语言支持（简体中文 `zh-CN`、英文 `en`）与深色/浅色主题自由切换。

---

## 🛠 技术架构与技术栈

### 架构分层

```text
jing-sync/
├── cmd/main.go               # 应用程序统一启动入口
├── boot/                     # 引导与初始化模块 (配置、日志、国际化、数据库、Gin 引擎)
├── api/                      # 接口层 (路由定义、请求控制器、认证/超时/日志中间件)
├── internal/                 # 核心业务逻辑层
│   ├── models/               # GORM 数据模型 (BaseModel, User, Engine, Job, SyncLog)
│   ├── services/             # 业务服务 (SyncService, OpenListClient, CronManager, BaseService)
│   └── utils/                # 通用工具集 (JWT Token, HTTP 请求客户端, 文件与响应工具)
├── frontend/                 # 前端工程 (Vue 3 + TypeScript + Vite + Element Plus + Tailwind CSS)
└── public/                   # 静态嵌入资源 (Go embed 挂载的前端 dist 与 i18n toml)
```

### 技术选型

- **后端 (Backend)**：
  - 语言：Go 1.23+
  - Web 框架：Gin Framework
  - ORM 与存储：GORM + SQLite (轻量级嵌入式数据库，自动执行迁移)
  - 定时调度：`robfig/cron/v3`
  - 日志系统：Logrus (集成按日日志归档)
  - 认证安全：golang-jwt/v5 + SHA256 密码加盐哈希
- **前端 (Frontend)**：
  - 框架：Vue 3 (Composition API / `<script setup>`) + TypeScript
  - 构建工具：Vite
  - UI 组件库：Element Plus + Tailwind CSS 4.x
  - 状态管理：Pinia
  - 路由与国际化：Vue Router + Vue I18n

---

## 🚀 快速开始与构建命令

本项目配备了跨平台优化的 `Makefile`，支持 Windows (PowerShell) 与 Unix (macOS / Linux)。

### 常用命令指南

```bash
# 查看所有可用命令及详细帮助
make help

# 1. 一键全栈构建（自动编译前端并嵌入至后端独立二进制中）
make build        # 或 make all

# 2. 运行构建产物
./bin/jing-sync    # macOS / Linux
# .\bin\jing-sync.exe # Windows

# 3. 本地开发模式
make dev-frontend # 启动前端 Vite 开发服务器 (快捷方式: make df, 默认端口 5173)
make dev-backend  # 以源码模式启动 Go 后端 (快捷方式: make db, 默认端口 8888)

# 4. 单独模块构建
make build-frontend # 仅编译前端并拷贝至 public/web (快捷方式: make bf)
make build-backend  # 仅编译后端二进制 (快捷方式: make bb)

# 5. 清理构建产物
make clean        # 清理 bin/, dist/, public/web/
make clean-all    # 深度清理（包含前端 node_modules 与运行日志）
```

---

## 🧪 测试与质量保证

项目前后端均配备完善的自动化测试套件：

### 后端单元测试 (Go Test)
后端测试覆盖了中间件、数据库 CRUD、路由树注册、OpenList 客户端交互及文件同步核心逻辑，并内置了针对无网络沙箱环境的内存 Mock 传输层（`HandlerTransport`）：

```bash
# 运行全部后端单元测试
make test

# 运行数据竞态安全检测测试
make test-race
```

### 前端单元测试 (Vitest)
前端测试覆盖了通用工具库、Pinia 状态管理 Store 以及核心视图组件：

```bash
# 运行前端 Vitest 单元测试
make test-frontend
```

### 一键全量测试
```bash
# 依次运行后端单元测试与前端单元测试
make test-all
```

---

## ⚙️ 配置说明

运行时数据与配置文件位于根目录下的 `data/` 目录：
- `data/config_prod.json`：生产运行配置文件，支持自定义监听端口（默认 `8888`）、数据库文件名及超时时长。
- `data/secret.key`：JWT 签名加密密钥文件（首次启动自动生成）。
- `data/jing-sync.db`：SQLite 数据库持久化存储文件。
- `data/logs/`：运行时系统业务与 HTTP 请求日志。

> **提示**：首次启动程序时，终端将会输出系统自动初始化的管理员账号与高强度随机密码，请妥善保存并在登录后于「用户设置」页面修改。

---

## 🐳 Docker 容器化部署

本项目提供了优化的多阶段构建 `Dockerfile`：

```bash
# 1. 使用 Makefile 构建镜像
make docker-build
# 或原生命令: docker build -t jing-sync .

# 2. 运行 Docker 容器（挂载宿主机目录实现数据与配置持久化）
make docker-run
# 或原生命令:
docker run -d \
  --name jing-sync \
  -p 8888:8888 \
  -v $(pwd)/data:/app/data \
  --restart unless-stopped \
  jing-sync
```

访问地址：`http://127.0.0.1:8888/`

---

## 📄 开源许可证

本项目基于 [MIT 许可证](LICENSE) 开源。
