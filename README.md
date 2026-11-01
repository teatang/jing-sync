# Jing-Sync

Jing-Sync 是一个功能强大、现代化的全-stack 文件同步与管理 Web 应用程序。它支持通过 OpenList / AList 等存储引擎进行跨目录、跨网盘的高效文件同步、定时调度、镜像对齐以及详细的同步日志审计。

---

## 🌟 核心功能特性

- **多存储引擎支持**：集成 OpenList / AList 实例，支持管理多个存储引擎的地址与访问令牌（Token）。
- **灵活的路径选择与浏览**：
  - 支持直接输入或粘贴任意源路径与目标路径。
  - 内置路径浏览弹窗（Path Browser），支持点击目录逐级浏览、跳转和快捷选择。
- **多种同步模式**：
  - **仅新增 (增量)**：仅复制目标端缺少或修改的文件，安全无覆盖。
  - **完全同步 (镜像)**：目标端与源端完全对齐，源端删除的文件会自动在目标端同步删除。
  - **移动文件 (剪切)**：同步到目标端后从源端删除。
- **强大的文件过滤与限流规则**：
  - 支持排除通配符规则（Glob Pattern，如 `*.tmp, *.bak, .git/**`）。
  - 支持文件大小上下限过滤（最小/最大字节限制）。
  - 支持多种扫描速率策略（标准、快速缓存、低速限频）。
- **定时调度与动态管理**：
  - 支持按分钟定时间隔或标准的 5 段式 Cron 表达式（如 `0 2 * * *`）。
  - 全局动态 Cron 调度管理器（CronManager），在编辑、新增、启用/禁用或删除任务时自动重新调度，避免任务重复添加。
- **详细的同步日志与终端抽屉**：
  - 记录每次同步的详细统计（发现文件、新增、删除、跳过、失败）。
  - 提供终端风格的日志详情抽屉（Sync Log Detail Drawer），方便排查与审计。
- **安全与权限管理**：
  - 基于 JWT 的用户认证。
  - 用户管理页面支持账号维护、密码脱敏显示与安全加密。
  - 内置源路径与目标路径重叠安全检查，防止死循环传输。

---

## 🛠 技术栈

- **后端 (Backend)**：
  - Go 1.23+
  - Gin Web Framework
  - GORM ORM + SQLite 数据库
  - `robfig/cron` 定时调度器
- **前端 (Frontend)**：
  - Vue 3 + TypeScript
  - Vite + Element Plus + Tailwind CSS 4.x
  - Pinia 状态管理 + Vue Router + Vue I18n

---

## 🚀 快速开始与构建指令

### 1. 构建与运行

```bash
# 全局生产构建（同时编译前端并打包至 Go 后端静态目录）
make build        # 或 make all

# 运行构建后的二进制程序
./bin/jing-sync    # Linux / macOS
# 或 ./bin/jing-sync.exe # Windows
```

### 2. 开发模式

```bash
# 启动前端 Vite 开发服务器 (端口 5173)
make dev-frontend # 或 make df

# 运行 Go 后端开发服务器
make dev-backend  # 或 make db
```

### 3. 单独构建

```bash
make build-frontend  # 编译 Vue 前端并复制到 public/web
make build-backend   # 编译 Go 后端二进制到 bin/jing-sync
make clean           # 清理构建产物
```

---

## 🐳 Docker 部署

```bash
# 构建 Docker 镜像
docker build -t jing-sync .

# 运行容器 (挂载数据目录持久化 SQLite 数据库与日志)
docker run -p 8888:8888 -v /path/to/data:/app/data jing-sync
```

---

## 📄 许可证

MIT License
