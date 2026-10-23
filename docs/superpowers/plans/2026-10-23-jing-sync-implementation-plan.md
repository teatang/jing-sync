# jing-sync 功能移植与架构重构开发计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 taosync（Python）的 OpenList 数据同步核心功能完整移植到 jing-sync（Go + Vue 3），优化同步引擎逻辑，并实现现代化、极简、支持日夜模式切换的前端界面，同时为关键逻辑补充 Testify 与 Vitest 自动化单元测试。

**Architecture:** 
- 后端采用 Go 1.23+ / Gin 框架分层架构（Route -> Controller -> Service -> Model/GORM -> SQLite），封装 `OpenListClient` 与 `SyncService` 实现对 OpenList/AList v3+ 接口的文件扫描、差异比对、任务提交与状态轮询，并提供 Cron 定时调度。
- 前端基于 Vue 3 + TypeScript + Vite + Pinia + Element Plus + Tailwind CSS 4.x 构建极简响应式界面，集成 `@vueuse/core` 实现浅色/深色/跟随系统三态切换。

**Tech Stack:** 
- 后端: Go 1.23, Gin, GORM, SQLite, Testify, robfig/cron/v3, Logrus
- 前端: Vue 3.5, TypeScript 5.8, Vite 6, Pinia 3, Element Plus, Tailwind CSS 4.x, Vitest, @vueuse/core, pnpm

**Spec:** `/Users/tea/www/cc-test-pro1/CLAUDE.md` 及用户需求说明

## Global Constraints

- 前端包管理必须使用 `pnpm`。
- 前端必须基于 TypeScript，代码添加清晰规范的中文注释。
- UI 必须遵循极简风格，支持浅色（Light）、深色（Dark）、跟随系统（Auto）模式。
- 后端代码必须严格基于 Gin 与 SQLite，关键逻辑编写 Testify 单元测试。
- Go 端在编译或运行测试前，必须确保 `jing-sync/public/web` 目录存在，避免 `embed.FS` 编译错误。

## Review Focus

1. **虚拟路径重叠与安全检查**: 源路径与目标路径完全相同或相互嵌套时，必须拒绝同步，避免无限递归与死循环。
2. **文件大小与通配符过滤准确性**: 文件大小区间（Min/Max）和通配符过滤（排除文件/文件夹）必须在文件差异对比前生效，避免不必要的文件网络请求。
3. **OpenList API 状态轮询防悬挂**: 提交复制任务后，轮询任务状态需设定超时与重试上限，任务被取消或异常退出时需正确更新日志状态。
4. **主题切换闪烁与多端同步**: 切换日夜模式时需同步更新 `html.dark` 类名、Element Plus 暗色变量及本地 `localStorage`，跟随系统时响应 `prefers-color-scheme` 媒体查询。
5. **任务并发安全与状态隔离**: 同一个 Job 正在执行时，禁止重复并发触发，调度器与手动触发共享相同的锁或状态保护。

---

### Task 1: 后端引入 Testify 测试框架并为基础工具库编写单元测试

**Files:**
- Modify: `jing-sync/go.mod`
- Modify: `jing-sync/internal/utils/utils.go`
- Create: `jing-sync/internal/utils/utils_test.go`
- Test: `jing-sync/internal/utils/utils_test.go`

**Interfaces:**
- Consumes: `internal/utils` 中的 `SecureRandString`, `Password2hash`, `CheckPasswordHash`, `GenerateToken`, `ParseToken`, `FormatBytes`
- Produces: 完备的 Testify 测试用例，确保密码加盐哈希、JWT 鉴权凭证生成校验、容量格式化逻辑 100% 可靠。

- [ ] **Step 1: 添加 `github.com/stretchr/testify` 依赖到 go.mod**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && go get github.com/stretchr/testify
```

- [ ] **Step 2: 在 `internal/utils/utils.go` 中补充辅助工具函数（如字节格式化 `FormatBytes`）并写出失败测试**

在 `jing-sync/internal/utils/utils_test.go` 中编写 Testify 测试用例：
```go
func TestFormatBytes(t *testing.T) {
    assert.Equal(t, "1.00 KB", utils.FormatBytes(1024))
    assert.Equal(t, "1.50 MB", utils.FormatBytes(1572864))
}
```

- [ ] **Step 3: 运行测试验证失败**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/utils -run TestFormatBytes`
Expected: FAIL (函数未定义)

- [ ] **Step 4: 实现 `FormatBytes(bytes int64) string` 并补充完善 JWT/密码哈希测试**

在 `internal/utils/utils.go` 中实现 `FormatBytes`，并在 `utils_test.go` 中使用 `assert` 与 `require` 覆盖密码生成、哈希比对、JWT 签发解析。

- [ ] **Step 5: 运行全部 utils 测试验证通过**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/utils`
Expected: PASS

- [ ] **Step 6: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && git add go.mod go.sum internal/utils/
git commit -m "feat(backend): add testify and comprehensive tests for utils"
```

---

### Task 2: 移植并实现文件过滤引擎（大小过滤与通配符规则匹配）

**Files:**
- Create: `jing-sync/internal/services/filter/filter.go`
- Create: `jing-sync/internal/services/filter/filter_test.go`
- Test: `jing-sync/internal/services/filter/filter_test.go`

**Interfaces:**
- Consumes: 无
- Produces: 
  - `IsFileSizeAllowed(size int64, minSize *int64, maxSize *int64) bool`
  - `PathFilter`: 基于通配符（如 `*.tmp`, `node_modules/*`, `.git/**`）的路径排除过滤器
  - `VirtualPathsOverlap(p1, p2 string) bool`: 校验路径是否重叠避免循环同步

- [ ] **Step 1: 编写过滤引擎的 Testify 失败测试**

在 `jing-sync/internal/services/filter/filter_test.go` 中：
```go
func TestIsFileSizeAllowed(t *testing.T) {
    min := int64(100)
    max := int64(500)
    assert.True(t, IsFileSizeAllowed(200, &min, &max))
    assert.False(t, IsFileSizeAllowed(50, &min, &max))
    assert.False(t, IsFileSizeAllowed(600, &min, &max))
}

func TestVirtualPathsOverlap(t *testing.T) {
    assert.True(t, VirtualPathsOverlap("/data/a", "/data/a"))
    assert.True(t, VirtualPathsOverlap("/data/a", "/data/a/sub"))
    assert.False(t, VirtualPathsOverlap("/data/a", "/data/b"))
}
```

- [ ] **Step 2: 运行测试验证失败**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/services/filter`
Expected: FAIL (package or function not defined)

- [ ] **Step 3: 实现 `filter.go`**

实现：
1. `IsFileSizeAllowed(size int64, minSize, maxSize *int64) bool`
2. `NormalizeVirtualPath(path string) string`
3. `VirtualPathsOverlap(p1, p2 string) bool`
4. `PathMatcher`: 支持以行分隔的排除规则列表，支持目录尾部 `/`、`*`、`**` 前缀后缀匹配。

- [ ] **Step 4: 运行测试验证全部通过**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/services/filter`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && git add internal/services/filter/
git commit -m "feat(backend): implement file size and path pattern filter engine"
```

---

### Task 3: 扩充 Job 与 SyncLog 数据模型（支持同步模式与高级过滤参数）

**Files:**
- Modify: `jing-sync/internal/models/jobs.go`
- Modify: `jing-sync/internal/models/sync_logs.go`
- Create: `jing-sync/internal/models/models_test.go`
- Test: `jing-sync/internal/models/models_test.go`

**Interfaces:**
- Consumes: `gorm.io/gorm`
- Produces: 
  - `models.Job` 包含 `SyncMode` (0: 全同步, 1: 仅新增, 2: 移动), `MinSize`, `MaxSize`, `Excludes`, `Status` (启用/禁用)
  - `models.SyncLog` 包含结构化统计指标（扫描数、新增数、耗时、详细步骤）

- [ ] **Step 1: 编写模型迁移与默认字段赋值的失败测试**

在 `jing-sync/internal/models/models_test.go` 中测试在内存 SQLite 中 AutoMigrate，并验证 Job 字段与状态常量。

- [ ] **Step 2: 运行测试验证失败**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/models`
Expected: FAIL (字段未定义)

- [ ] **Step 3: 完善 `jobs.go` 与 `sync_logs.go` 结构体**

```go
type Job struct {
    BaseModel
    Remark       string `json:"remark" gorm:"column:remark;size:255"`
    SrcPath      string `json:"src_path" gorm:"column:src_path;size:500"`
    DstPath      string `json:"dst_path" gorm:"column:dst_path;size:500"`
    EngineId     uint   `json:"engine_id" gorm:"column:engine_id"`
    UserId       uint   `json:"user_id" gorm:"column:user_id"`
    Speed        uint8  `json:"speed" gorm:"column:speed;default:0"` // 0:正常 1:快速(不刷新) 2:慢速(等待1s)
    SyncMode     uint8  `json:"sync_mode" gorm:"column:sync_mode;default:1"` // 0:全同步 1:仅新增 2:移动
    Interval     uint   `json:"interval" gorm:"column:interval;default:60"`
    IsCron       uint8  `json:"is_cron" gorm:"column:is_cron;default:0"`
    CronExpr     string `json:"cron_expr" gorm:"column:cron_expr;size:50"`
    MinSize      *int64 `json:"min_size" gorm:"column:min_size"`
    MaxSize      *int64 `json:"max_size" gorm:"column:max_size"`
    Excludes     string `json:"excludes" gorm:"column:excludes;type:text"` // 换行分隔的排除通配符
    Enabled      uint8  `json:"enabled" gorm:"column:enabled;default:1"`   // 0:禁用 1:启用
}
```

- [ ] **Step 4: 运行测试验证迁移成功**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/models`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && git add internal/models/
git commit -m "feat(backend): expand Job and SyncLog models with sync modes and filter rules"
```

---

### Task 4: 完整实现 OpenListClient（目录递归扫描、文件复制、任务查询）

**Files:**
- Modify: `jing-sync/internal/services/open_list_service.go`
- Create: `jing-sync/internal/services/open_list_client_test.go`
- Test: `jing-sync/internal/services/open_list_client_test.go`

**Interfaces:**
- Consumes: OpenList / AList v3 REST API
- Produces: 
  - `ListFiles(path string, refresh bool) ([]FileInfo, error)`（同时获取文件与目录，包括大小与修改时间）
  - `CopyFile(srcDir, dstDir, name string) (taskId string, err error)`
  - `MoveFile(srcDir, dstDir, name string) (taskId string, err error)`
  - `GetTaskStatus(taskId string) (state int, progress float64, err error)`
  - `TestConnection() error`

- [ ] **Step 1: 编写使用 httptest.Server Mock 的失败测试**

在 `jing-sync/internal/services/open_list_client_test.go` 中 mock OpenList 响应（`/api/fs/list`, `/api/fs/copy`, `/api/me`），断言方法返回正确的结构化文件列表和任务 ID。

- [ ] **Step 2: 运行测试验证失败**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/services -run TestOpenListClient`
Expected: FAIL

- [ ] **Step 3: 在 `open_list_service.go` 中实现核心 API 交互方法**

1. 扩展 `ChildPathRawInfo` 包含 `Size int64`。
2. 实现 `ListFiles(path string, speed int) ([]FileInfo, error)`。
3. 实现 `CopyFile` 调用 `/api/fs/copy`。
4. 实现 `MoveFile` 调用 `/api/fs/move`。
5. 实现 `TestConnection` 调用 `/api/me` 或 `/api/fs/list` 根路径。

- [ ] **Step 4: 运行测试验证通过**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/services -run TestOpenListClient`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && git add internal/services/
git commit -m "feat(backend): implement complete OpenList API client operations"
```

---

### Task 5: 完整实现 SyncService 同步比对与执行引擎

**Files:**
- Modify: `jing-sync/internal/services/sync_service.go`
- Create: `jing-sync/internal/services/sync_service_test.go`
- Test: `jing-sync/internal/services/sync_service_test.go`

**Interfaces:**
- Consumes: `OpenListClient`, `filter.PathFilter`, `models.Job`, `models.SyncLog`
- Produces: `ExecuteSync(jobId uint) (*models.SyncLog, error)`
  - 路径合法性与重叠性检测
  - 递归扫描源端与目的端目录树
  - 应用过滤规则（排除模式、文件大小范围）
  - 差异计算（Mode 0 全量比对、Mode 1 增量新增、Mode 2 移动）
  - 顺序/并发提交任务并等待结果，写入详细的同步进度记录

- [ ] **Step 1: 编写 SyncService 差异算法与执行流程的单元测试**

在 `sync_service_test.go` 中，构造虚拟源文件树和目标文件树，测试增量（Add-only）模式下正确识别出待同步文件，排除过滤文件，拦截重叠路径。

- [ ] **Step 2: 运行测试验证失败**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/services -run TestSyncService`
Expected: FAIL

- [ ] **Step 3: 重构并实现 `sync_service.go`**

实现真正的对比与同步循环，替换原有的 mock 代码。支持错误重试、记录每批次同步细节到 `SyncLog.Details`（JSON 格式）。

- [ ] **Step 4: 运行测试验证通过**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./internal/services -run TestSyncService`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && git add internal/services/sync_service*
git commit -m "feat(backend): implement real differential synchronization engine"
```

---

### Task 6: 完善 API 接口层（引擎连通性测试、任务手动触发、日志查询）

**Files:**
- Modify: `jing-sync/api/controllers/engine_controller.go`
- Modify: `jing-sync/api/controllers/job_controller.go`
- Modify: `jing-sync/api/controllers/sync_log_controller.go`
- Modify: `jing-sync/api/route/protected_route.go`
- Create: `jing-sync/api/controllers/controllers_test.go`
- Test: `jing-sync/api/controllers/controllers_test.go`

**Interfaces:**
- Consumes: HTTP Gin 请求
- Produces: 
  - `POST /api/engine/test`: 测试引擎连通性
  - `POST /api/job/:job_id/sync`: 手动立即触发同步任务
  - `GET /api/job/:job_id/sync-logs`: 分页查询任务同步日志
  - `PUT /api/job/:job_id/toggle`: 快速启用/禁用定时任务

- [ ] **Step 1: 编写控制器接口的 HTTP 测试**

使用 Gin 的测试模式 `httptest.NewRecorder()`，编写接口调用与 JSON 返回结构断言。

- [ ] **Step 2: 运行测试验证失败**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./api/controllers`
Expected: FAIL (新路由或方法未就绪)

- [ ] **Step 3: 在控制器和路由中实现这些端点**

在 `engine_controller.go` 中增加 `TestEngine` 方法，在 `job_controller.go` 中增加状态切换方法，并在 `protected_route.go` 中注册路由。

- [ ] **Step 4: 运行测试验证全部通过**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync && go test -v ./api/controllers`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && git add api/
git commit -m "feat(backend): add engine connectivity test and job management endpoints"
```

---

### Task 7: 前端工程配置 Vitest 测试框架与 @vueuse/core

**Files:**
- Modify: `jing-sync/frontend/package.json`
- Create: `jing-sync/frontend/vitest.config.ts`
- Create: `jing-sync/frontend/src/utils/__tests__/rest-client.spec.ts`
- Create: `jing-sync/frontend/src/stores/__tests__/localeStore.spec.ts`
- Test: `pnpm run test:unit`

**Interfaces:**
- Consumes: `pnpm`
- Produces: 可运行的 Vitest 单元测试命令与配置环境

- [ ] **Step 1: 安装 `@vueuse/core` 与 `vitest` 相关依赖**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm add @vueuse/core && pnpm add -D vitest @vue/test-utils jsdom
```

- [ ] **Step 2: 创建 `vitest.config.ts` 并在 `package.json` 添加测试脚本**

配置 `test: { environment: 'jsdom' }`，在 `package.json` scripts 中添加 `"test:unit": "vitest run"`。

- [ ] **Step 3: 编写 token 与 rest-client 的 Vitest 测试**

在 `src/utils/__tests__/rest-client.spec.ts` 中断言 token 设置、提取、清除逻辑。

- [ ] **Step 4: 运行 Vitest 测试验证通过**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run test:unit`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && git add package.json pnpm-lock.yaml vitest.config.ts src/
git commit -m "chore(frontend): configure Vitest and VueUse with initial unit tests"
```

---

### Task 8: 实现日夜模式（浅色/深色/跟随系统）主题管理

**Files:**
- Create: `jing-sync/frontend/src/stores/themeStore.ts`
- Create: `jing-sync/frontend/src/components/ThemeSwitch.vue`
- Create: `jing-sync/frontend/src/stores/__tests__/themeStore.spec.ts`
- Modify: `jing-sync/frontend/src/assets/tailwind.css`
- Test: `pnpm run test:unit`

**Interfaces:**
- Consumes: `@vueuse/core` (`useColorMode`, `useDark`)
- Produces: 
  - `useThemeStore`: 管理 `'light' | 'dark' | 'auto'` 状态，响应式同步 `html.dark` 类名
  - `ThemeSwitch.vue`: 极简下拉菜单或按钮组，清晰呈现当前模式图标与切换项

- [ ] **Step 1: 编写 `themeStore` 的 Vitest 测试**

测试手动设置 `'dark'`、`'light'`、`'auto'` 时，状态正确切换且持久化到 localStorage。

- [ ] **Step 2: 运行测试验证失败**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run test:unit`
Expected: FAIL (store not found)

- [ ] **Step 3: 实现 `themeStore.ts` 与 `ThemeSwitch.vue`**

结合 Element Plus 的暗色主题模式（引入 `element-plus/theme-chalk/dark/css-vars.css`）与 Tailwind 的 `.dark` 选择器，确保全局背景、边框、卡片文字自适应。

- [ ] **Step 4: 运行测试验证通过**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run test:unit`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && git add src/stores/ src/components/ src/assets/
git commit -m "feat(frontend): implement comprehensive light/dark/system theme switcher"
```

---

### Task 9: 重构全局布局 Layout（极简审美、响应式导航与头部工具栏）

**Files:**
- Modify: `jing-sync/frontend/src/views/Layout.vue`
- Modify: `jing-sync/frontend/src/router/index.ts`
- Modify: `jing-sync/frontend/src/locales/zh-CN.ts`
- Modify: `jing-sync/frontend/src/locales/en.ts`
- Test: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`

**Interfaces:**
- Consumes: `ThemeSwitch.vue`, `localeStore`, `userStore`, 路由系统
- Produces: 具有现代极简视觉的顶栏（折叠控制、Logo、主题切换、语言切换、用户头像/退出登录）与侧边菜单栏。

- [ ] **Step 1: 扩展路由配置增加“同步日志”页面路由**

在 `router/index.ts` 中注册 `/sync-logs` 路由，更新多语言字典中的导航文案。

- [ ] **Step 2: 重新设计并精简 `Layout.vue` 样式与结构**

移除硬编码背景色，采用 Tailwind 类名（`bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800`），嵌入 `ThemeSwitch` 组件。

- [ ] **Step 3: 执行前端编译检查**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`
Expected: PASS

- [ ] **Step 4: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && git add src/views/Layout.vue src/router/ src/locales/
git commit -m "refactor(frontend): modernize layout with minimalist aesthetics and responsive navigation"
```

---

### Task 10: 优化引擎管理视图（新增连通性测试与清晰反馈）

**Files:**
- Modify: `jing-sync/frontend/src/views/engine/Engine.vue`
- Modify: `jing-sync/frontend/src/views/engine/EngineDialog.vue`
- Test: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`

**Interfaces:**
- Consumes: `/api/engine`, `/api/engine/test`
- Produces: 具备“测试连接”功能、卡片/表格化展示、状态徽章的引擎管理页面。

- [ ] **Step 1: 在 `EngineDialog.vue` 中添加“测试连接”按钮与请求逻辑**

点击测试时调用 `/api/engine/test`，测试成功或失败以 Element Plus `ElMessage` 和视觉标签给出明确提示。

- [ ] **Step 2: 优化 `Engine.vue` 列表展示**

添加操作区（编辑、删除、连通性复查），在暗黑模式下呈现高对比度极简排版。

- [ ] **Step 3: 执行前端类型与构建校验**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`
Expected: PASS

- [ ] **Step 4: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && git add src/views/engine/
git commit -m "feat(frontend): enhance engine management with connection testing and modern cards"
```

---

### Task 11: 升级定时任务视图（同步模式选择、文件过滤、手动立即执行）

**Files:**
- Modify: `jing-sync/frontend/src/views/cron-job/CronJob.vue`
- Modify: `jing-sync/frontend/src/views/cron-job/CronJobDialog.vue`
- Create: `jing-sync/frontend/src/views/cron-job/PathSelectorDialog.vue`
- Test: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`

**Interfaces:**
- Consumes: `/api/job`, `/api/job/:job_id/sync`, `/api/open-list`
- Produces: 
  - 支持配置同步模式（全同步 / 仅新增 / 移动）
  - 支持配置文件大小区间（带 B/KB/MB/GB 单位选择）与排除通配符列表
  - 支持在表单中点击弹窗浏览 OpenList 远程路径并选定
  - 任务列表支持一键“立即执行”并给出反馈

- [ ] **Step 1: 创建 `PathSelectorDialog.vue` 路径选择器组件**

调用 `/api/open-list?engine_id=...&path=...`，支持层级深入浏览并一键选择目录。

- [ ] **Step 2: 改造 `CronJobDialog.vue` 表单项**

增加同步模式单选组、文件大小输入行（单位下拉）、多行文本输入排除规则（并附带通配符提示示例）。

- [ ] **Step 3: 改造 `CronJob.vue` 列表与操作栏**

增加启用/禁用快捷开关、立即同步按钮、查看日志入口按钮。

- [ ] **Step 4: 执行构建测试**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`
Expected: PASS

- [ ] **Step 5: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && git add src/views/cron-job/
git commit -m "feat(frontend): upgrade cron-job views with sync modes, filters and manual trigger"
```

---

### Task 12: 实现同步日志与详情抽屉视图

**Files:**
- Create: `jing-sync/frontend/src/views/sync-log/SyncLogView.vue`
- Create: `jing-sync/frontend/src/views/sync-log/SyncLogDetailDrawer.vue`
- Modify: `jing-sync/frontend/src/router/index.ts`
- Test: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`

**Interfaces:**
- Consumes: `/api/job/:job_id/sync-logs`, `/api/sync-log/:id`
- Produces: 
  - 同步记录分页列表（任务名称、状态、成功/失败标签、耗时、新增文件数、开始时间）
  - 详情抽屉展开展示 JSON 解析后的详细同步文件列表与错误原因排查

- [ ] **Step 1: 创建 `SyncLogView.vue` 与 `SyncLogDetailDrawer.vue`**

构建表格、状态 Tag（成功显示绿色/失败显示红色）、查看详情抽屉组件。

- [ ] **Step 2: 在路由与菜单中联调**

在 `Layout.vue` 菜单中增加“同步日志”，路由配置连接到 `SyncLogView.vue`。

- [ ] **Step 3: 执行构建测试**

Run: `cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run build`
Expected: PASS

- [ ] **Step 4: 提交代码**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && git add src/views/sync-log/ src/router/ src/views/Layout.vue
git commit -m "feat(frontend): add sync logs view and execution details drawer"
```

---

### Task 13: 全流程集成构建、全量自动化测试与工程验证

**Files:**
- Modify: `jing-sync/Makefile`
- Modify: `jing-sync/CLAUDE.md`
- Test: 全量单元测试与打包构建

**Interfaces:**
- Consumes: 所有前后端产物
- Produces: 自动化通过的测试套件与独立的生产级打包可执行文件

- [ ] **Step 1: 运行后端全部 Testify 单元测试**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && make test
```
Expected: PASS (所有 package 均通过测试)

- [ ] **Step 2: 运行前端全部 Vitest 单元测试**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync/frontend && pnpm run test:unit
```
Expected: PASS

- [ ] **Step 3: 执行全流程打包**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && make build
```
Expected: 
1. `frontend/dist` 成功生成
2. 复制到 `public/web`
3. 编译出 `bin/jing-sync` 单一二进制

- [ ] **Step 4: 验证启动**

运行 `./bin/jing-sync` 并验证 HTTP 请求响应正常。

- [ ] **Step 5: 提交所有文档与配置更新**

```bash
cd /Users/tea/www/cc-test-pro1/jing-sync && git add Makefile CLAUDE.md
git commit -m "chore: finalize build workflow, documentation and full test suite verification"
```
