# ==============================================================================
# jing-sync 项目自动化管理 Makefile
# 包含前后端编译、本地开发热重载、自动化测试、Docker 容器化管理及清理维护命令
# ==============================================================================

# ------------------------------------------------------------------------------
# 伪目标声明 (.PHONY)
# 声明所有非实际文件生成的命令目标，防止与同名文件产生冲突，并提升执行效率
# ------------------------------------------------------------------------------
.PHONY: all help build clean clean-all \
        build-frontend dev-frontend frontend-deps check-frontend-deps bf df \
        build-backend dev-backend run-backend run bb db ensure-web-dir \
        test test-race test-frontend test-all \
        docker-build docker-run

# ==============================================================================
# 1. 项目基础配置与路径定义
# ==============================================================================

# 项目名称（用于输出二进制文件名、Docker 镜像名等）
PROJECT_NAME       := jing-sync

# Go 编译器可执行程序与构建参数
GO                 ?= go
# -ldflags="-s -w" 用于去除符号表与调试信息（DWARF），显著缩减编译后的二进制体积
GO_BUILD_FLAGS     ?= -ldflags="-s -w"
# Go 入口文件路径定义
GO_MAIN_DIR        := ./cmd
GO_MAIN_FILE       := main.go

# 前端包管理器（支持 pnpm / npm / yarn，可通过命令行覆盖，如 make build NPM=npm）
NPM                ?= pnpm

# 目录定义
BACKEND_BIN_DIR    := ./bin
FRONTEND_DIR       := ./frontend
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/dist
# 前端静态资源在后端的输出挂载目录（供 Go embed 静态打包嵌入）
FRONTEND_OUTPUT_TO := ./public/web

# Docker 镜像与容器配置
DOCKER_IMAGE_NAME  ?= $(PROJECT_NAME)
DOCKER_PORT        ?= 8888

# ==============================================================================
# 2. 跨平台操作系统适配 (Windows 与 类 Unix / macOS / Linux)
# ==============================================================================

ifeq ($(OS),Windows_NT)
    PLATFORM          := windows
    BACKEND_BIN_NAME  := $(PROJECT_NAME).exe
    # Windows 环境下使用 PowerShell 执行脚本并禁用加载配置文件以保证执行效率
    POWERSHELL_PREFIX := powershell.exe -NoProfile -Command
else
    PLATFORM          := unix
    BACKEND_BIN_NAME  := $(PROJECT_NAME)
endif

# 后端最终输出二进制文件完整路径
BACKEND_OUTPUT     := $(BACKEND_BIN_DIR)/$(BACKEND_BIN_NAME)

# ------------------------------------------------------------------------------
# 跨平台通用工具宏函数定义
# ------------------------------------------------------------------------------

# RMDIR: 递归删除指定目录（若目录不存在则忽略，静默退出）
# 参数 $(1): 待删除的目录路径
define RMDIR
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { if (Test-Path -Path \"$1\" -PathType Container) { Remove-Item -Recurse -Force -Path \"$1\" -ErrorAction SilentlyContinue | Out-Null } }", \
		@rm -rf "$(1)" \
	)
endef

# MKDIR: 创建指定目录（若已存在则忽略）
# 参数 $(1): 待创建的目录路径
define MKDIR
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { if (-not (Test-Path -Path \"$1\")) { New-Item -ItemType Directory -Path \"$1\" -Force -ErrorAction SilentlyContinue | Out-Null } }", \
		@mkdir -p "$(1)" \
	)
endef

# COPY_CONTENTS: 复制源目录中的所有内容到目标目录（保持目录结构，包含子文件）
# 参数 $(1): 源目录路径 (如 $(FRONTEND_BUILD_DIR))
# 参数 $(2): 目标目录路径 (如 $(FRONTEND_OUTPUT_TO))
define COPY_CONTENTS
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { Copy-Item -Path \"$1\\*\" -Destination \"$2\" -Recurse -Force -ErrorAction SilentlyContinue | Out-Null }", \
		@cp -R "$(1)/." "$(2)/" \
	)
endef

# ==============================================================================
# 3. 默认目标与帮助文档
# ==============================================================================

# 默认执行完整构建
all: build

# 显示命令使用说明与帮助菜单
help:
	@echo "========================================================================"
	@echo "                   jing-sync Makefile 使用指南                          "
	@echo "========================================================================"
	@echo " [常用构建命令]"
	@echo "   make build           一键完整构建：先构建前端，再打包嵌入至后端二进制"
	@echo "   make build-frontend  仅构建前端代码并同步到 public/web (快捷方式: make bf)"
	@echo "   make build-backend   仅构建 Go 后端二进制程序至 bin/ (快捷方式: make bb)"
	@echo ""
	@echo " [本地开发与运行]"
	@echo "   make dev-frontend    启动前端 Vite 开发服务器（支持热重载，快捷方式: make df）"
	@echo "   make dev-backend     以源码方式启动 Go 后端（go run，快捷方式: make db）"
	@echo "   make run             编译并直接运行后端二进制程序 (同 make run-backend)"
	@echo ""
	@echo " [依赖管理与测试]"
	@echo "   make frontend-deps   强制安装/重新安装前端 node_modules 依赖"
	@echo "   make test            运行后端全部单元测试"
	@echo "   make test-race       运行后端竞态分析测试 (go test -race)"
	@echo "   make test-frontend   运行前端单元测试 (vitest)"
	@echo "   make test-all        运行全部测试 (前端 + 后端)"
	@echo ""
	@echo " [Docker 容器管理]"
	@echo "   make docker-build    基于 Dockerfile 构建多阶段精简容器镜像"
	@echo "   make docker-run      运行 Docker 容器实例（挂载本地 data 存储目录）"
	@echo ""
	@echo " [清理维护]"
	@echo "   make clean           清理构建生成目录 (bin/, frontend/dist/, public/web/)"
	@echo "   make clean-all       深度清理 (额外清理 frontend/node_modules 及日志)"
	@echo "========================================================================"

# ==============================================================================
# 4. 前端任务 (Vue 3 + Vite + TypeScript)
# ==============================================================================

# 依赖自检（内部依赖，仅在 node_modules 缺失时自动安装，避免每次构建重复安装导致阻塞或二次确认）
check-frontend-deps:
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { if (-not (Test-Path -Path \"$(FRONTEND_DIR)/node_modules\")) { Write-Host '>> 正在安装前端依赖 ($(NPM))...'; cd $(FRONTEND_DIR); $(NPM) install } }", \
		@if [ ! -d "$(FRONTEND_DIR)/node_modules" ]; then \
			echo ">> 检测到前端 node_modules 缺失，正在使用 $(NPM) 安装依赖..."; \
			cd $(FRONTEND_DIR) && $(NPM) install; \
		fi \
	)

# 强制安装/更新前端依赖
frontend-deps:
	@echo ">> 正在安装/更新前端依赖 ($(NPM))..."
	@cd $(FRONTEND_DIR) && $(NPM) install

# 构建前端生产产物并同步拷贝至后端 embed 嵌入目录
build-frontend: check-frontend-deps
	@echo ">> 开始构建前端项目 (Vite)..."
	@cd $(FRONTEND_DIR) && $(NPM) run build
	@echo ">> 前端构建完成，输出目录: $(FRONTEND_BUILD_DIR)"
	@echo ">> 正在同步前端产物至后端静态资源目录: $(FRONTEND_OUTPUT_TO)"
	$(call RMDIR,$(FRONTEND_OUTPUT_TO))
	$(call MKDIR,$(FRONTEND_OUTPUT_TO))
	$(call COPY_CONTENTS,$(FRONTEND_BUILD_DIR),$(FRONTEND_OUTPUT_TO))
	@echo ">> 前端产物同步完成。"

# 启动前端本地开发热重载服务器 (端口一般默认为 5173)
dev-frontend: check-frontend-deps
	@echo ">> 启动前端开发服务器 (Vite dev)..."
	@cd $(FRONTEND_DIR) && $(NPM) run dev

# 快捷方式别名
bf: build-frontend
df: dev-frontend

# ==============================================================================
# 5. 后端任务 (Go + Gin + SQLite)
# ==============================================================================

# 确保 public/web 目录存在且不为空（Go 1.16+ go:embed 特性要求目标目录必须包含匹配文件，否则报错）
ensure-web-dir:
	$(if $(filter windows,$(PLATFORM)), \
		@$(POWERSHELL_PREFIX) "& { if (-not (Test-Path -Path \"$(FRONTEND_OUTPUT_TO)\")) { New-Item -ItemType Directory -Path \"$(FRONTEND_OUTPUT_TO)\" -Force | Out-Null }; if ((Get-ChildItem -Path \"$(FRONTEND_OUTPUT_TO)\" -Force | Measure-Object).Count -eq 0) { New-Item -ItemType File -Path \"$(FRONTEND_OUTPUT_TO)/.gitkeep\" -Force | Out-Null; Write-Host '  [提示] public/web 为空，已自动创建占位文件以满足 Go embed 编译要求' } }", \
		@mkdir -p "$(FRONTEND_OUTPUT_TO)"; \
		if [ -z "$$(ls -A "$(FRONTEND_OUTPUT_TO)" 2>/dev/null)" ]; then \
			touch "$(FRONTEND_OUTPUT_TO)/.gitkeep"; \
			echo "  [提示] public/web 为空，已自动创建占位文件以满足 Go embed 编译要求"; \
		fi \
	)

# 构建 Go 后端二进制可执行文件
build-backend: ensure-web-dir
	@echo ">> 开始构建 Go 后端程序..."
	$(call RMDIR,$(BACKEND_BIN_DIR))
	$(call MKDIR,$(BACKEND_BIN_DIR))
	$(GO) build $(GO_BUILD_FLAGS) -o $(BACKEND_OUTPUT) $(GO_MAIN_DIR)/$(GO_MAIN_FILE)
	@echo ">> 后端构建成功，二进制文件位于: $(BACKEND_OUTPUT)"

# 编译并运行后端程序
run-backend: build-backend
	@echo ">> 正在启动已构建的后端程序..."
	$(BACKEND_OUTPUT)

# 快捷方式别名：直接运行
run: run-backend

# 以开发模式直接运行后端源码（若需要文件变动自动重启，建议搭配 air 工具）
dev-backend: ensure-web-dir
	@echo ">> 以开发模式启动 Go 后端程序 (go run)..."
	$(GO) run $(GO_MAIN_DIR)/$(GO_MAIN_FILE)

# 快捷方式别名
bb: build-backend
db: dev-backend

# ==============================================================================
# 6. 整合全栈构建
# ==============================================================================

# 一键完整打包前后端：先构建前端生成静态资源，再将其嵌入 Go 后端二进制文件中
build: build-frontend build-backend
	@echo "========================================================================"
	@echo "           jing-sync 全栈构建成功！可执行文件: $(BACKEND_OUTPUT)        "
	@echo "========================================================================"

# ==============================================================================
# 7. 测试任务
# ==============================================================================

# 运行后端单元测试
test:
	@echo ">> 运行后端 Go 单元测试..."
	@$(GO) test ./...

# 运行后端数据竞态检测测试
test-race:
	@echo ">> 运行后端竞态检测 (go test -race)..."
	@$(GO) test -race ./...

# 运行前端测试
test-frontend: check-frontend-deps
	@echo ">> 运行前端单元测试 (vitest)..."
	@cd $(FRONTEND_DIR) && $(NPM) run test:unit

# 运行全部测试套件
test-all: test test-frontend

# ==============================================================================
# 8. Docker 容器化支持
# ==============================================================================

# 构建 Docker 镜像
docker-build:
	@echo ">> 开始构建 Docker 镜像: $(DOCKER_IMAGE_NAME)..."
	docker build -t $(DOCKER_IMAGE_NAME) .
	@echo ">> Docker 镜像构建完成。"

# 运行 Docker 容器实例
docker-run:
	@echo ">> 启动 Docker 容器 $(DOCKER_IMAGE_NAME) (端口映射 $(DOCKER_PORT):8888)..."
	docker run -d --name $(DOCKER_IMAGE_NAME) -p $(DOCKER_PORT):8888 -v "$$(pwd)/data:/app/data" $(DOCKER_IMAGE_NAME)
	@echo ">> 容器已在后台运行，访问地址: http://127.0.0.1:$(DOCKER_PORT)/"

# ==============================================================================
# 9. 清理维护
# ==============================================================================

# 清理常规构建产物目录
clean:
	@echo ">> 正在清理构建产物与临时文件..."
	$(call RMDIR,$(BACKEND_BIN_DIR))
	$(call RMDIR,$(FRONTEND_BUILD_DIR))
	$(call RMDIR,$(FRONTEND_OUTPUT_TO))
	@$(GO) clean
	@echo ">> 清理完成。"

# 深度清理（包含前端依赖包与运行时日志数据，供环境重置使用）
clean-all: clean
	@echo ">> 正在执行深度清理 (依赖与日志)..."
	$(call RMDIR,$(FRONTEND_DIR)/node_modules)
	$(call RMDIR,data/logs)
	@echo ">> 深度清理完成。"
