# Go Vue Admin Starter

这是基于 `Wei-Shaw/sub2api` 仓库抽离并重建后的通用后台管理基础骨架。

当前活跃代码已经不再面向 AI 网关产品，而是一个可复用、可扩展、可继续二次开发的后台 starter framework。

## 项目定位

适合作为后续各类管理后台项目的起点，重点提供基础设施层与通用后台能力。

已保留：

- Go 后端：`Gin + Ent`
- Vue 3 前端：`Vite + TypeScript + TailwindCSS`
- 前端内置 `zh / en` 多语言能力，并提供顶部语言切换入口
- PostgreSQL / Redis 接入
- JWT 鉴权
- RBAC 权限控制
- 统一响应结构与统一错误处理
- 审计日志
- 用户、角色、菜单、系统配置等最小 CRUD 示例
- 内置个人中心、资料更新、修改密码等基础账号自助能力

已从活跃骨架中剥离：

- AI provider 适配
- 上游账号池
- 转发 / 中继 / 网关逻辑
- API Key 分发
- 计费 / 配额 / 价格 / 订阅
- 模型调度 / 粘性会话 / 供应商路由

旧业务材料已隔离到 `_legacy_*` 目录，不再属于当前 starter 的运行部分。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 后端 | Go 1.25, Gin, Ent, Viper, JWT |
| 前端 | Vue 3, Vite 5, TypeScript, TailwindCSS, Pinia, Vue Router |
| 数据库 | PostgreSQL |
| 缓存 | Redis |
| 工程化 | Docker Compose |

## 当前内置模块

后端模块：

- `auth`
- `user`
- `role`
- `menu`
- `audit_log`
- `system_config`

前端页面：

- 登录页
- 仪表盘
- 用户管理
- 角色管理
- 菜单管理
- 系统配置
- 审计日志
- 个人中心
- 403 无权限页
- 404 页面不存在

## 目录结构

```text
backend/
  cmd/server/
  ent/
  internal/
    app/
    auth/
    config/
    domain/
      audit/
      auth/
      menu/
      role/
      system/
      user/
    infrastructure/
      cache/
      db/
    middleware/
    pkg/
    server/
    util/
    web/
      errorx/
      response/
  migrations/

frontend/
  src/
    api/
      modules/
    components/
      base/
      business/
    composables/
    i18n/
    layouts/
    router/
    stores/
    styles/
    types/
    utils/
    views/
      audit/
      auth/
      dashboard/
      exception/
      system/
      user/
```

## 本地启动

重要：后端依赖 PostgreSQL，默认健康检查也会校验 Redis，所以必须先启动 `postgres` 和 `redis`，否则后端初始化会直接失败。

推荐启动顺序：

1. 先启动 PostgreSQL 和 Redis
2. 执行数据库迁移与种子初始化
3. 启动后端服务
4. 启动前端开发服务器

### 0. 先启动 PostgreSQL 和 Redis

在仓库根目录执行：

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api
docker compose --env-file .env.example up -d postgres redis
```

检查依赖状态：

```bash
docker compose --env-file .env.example ps -a
```

预期结果：

- `admin-starter-postgres` 状态为 `healthy`
- `admin-starter-redis` 状态为 `healthy`

### 1. 启动后端

首次启动前先执行迁移：

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
```

如果 PostgreSQL 没启动，这一步会直接失败，并输出明确的数据库连接提示。

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server
```

默认地址：

- `http://localhost:18080`
- 健康检查：`GET /healthz`

### 2. 启动前端

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/frontend
npm install
npm run dev
```

默认地址：

- `http://localhost:5173`

Vite 开发服务器默认会把 `/api` 代理到 `http://localhost:18080`。

## Docker Compose 启动

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api
cp .env.example .env
docker compose up --build
```

默认地址：

- 前端：`http://localhost:5173`
- 后端：`http://localhost:18080`
- PostgreSQL：`localhost:5432`
- Redis：`localhost:6379`

## 数据库初始化

当前 starter 使用 Ent 在应用启动阶段自动建表与初始化种子数据。

只执行迁移与初始化：

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
```

默认演示管理员账号：

- 用户名：`admin`
- 密码：`admin123456`

已验证的本地启动顺序：

```bash
cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api
docker compose --env-file .env.example up -d postgres redis

cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/backend
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server

cd /Volumes/NewDisk/Github/Wei-Shaw/sub2api/frontend
npm install
npm run dev
```

后端启动后可用以下命令检查：

```bash
curl http://127.0.0.1:18080/healthz
```

预期返回：

```json
{"code":"success","message":"success","data":{"database":"up","redis":"up"}}
```

## 构建校验

后端：

```bash
cd backend
GOTOOLCHAIN=local GOSUMDB=off go build ./...
```

前端：

```bash
cd frontend
npm install
npm run build
```

## 环境变量

项目通过 Viper 读取环境变量。

常用变量包括：

- `SERVER_PORT`
- `DATABASE_HOST`
- `DATABASE_PORT`
- `DATABASE_USER`
- `DATABASE_PASSWORD`
- `DATABASE_NAME`
- `REDIS_ENABLED`
- `REDIS_ADDR`
- `JWT_SECRET`
- `SEED_ADMIN_USERNAME`
- `SEED_ADMIN_PASSWORD`

完整本地模板见 [`.env.example`](./.env.example)。

## 如何新增一个业务模块

可以直接参照现有模块的结构扩展。

1. 在 `backend/ent/schema` 下新增 Ent schema。
2. 生成 Ent 代码。
3. 在 `backend/internal/domain/<module>/` 下创建：
   - `entity`
   - `repository`
   - `service`
   - `handler`
4. 在 `backend/internal/app/app.go` 与 `backend/internal/server/router.go` 中注册模块。
5. 在 `frontend/src/api/modules` 中补充接口封装。
6. 在 `frontend/src/views/<module>` 中新增页面。
7. 在 `frontend/src/router/index.ts` 中注册路由。
8. 如果模块需要出现在菜单中，同时补充菜单种子和权限点。

## 说明

- 仓库名与 Go module path 仍沿用了原始 GitHub 仓库路径，便于在当前仓库内继续演进。
- 活跃代码与文案已经切换为通用后台 starter 语义。
- `_legacy_*` 目录仅作为本次重构过程中隔离的历史材料。
