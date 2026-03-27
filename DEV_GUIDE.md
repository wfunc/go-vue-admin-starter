# Go Vue Admin Starter 开发指南

## 目标

当前仓库的活跃部分是一个通用后台管理骨架，不再承载 AI 网关业务语义。

开发时优先遵守以下原则：

- 公共能力留在 starter 层
- 新业务以独立模块方式新增
- 不要把业务耦合重新塞回基础设施层
- 优先补齐实体、权限、路由、菜单、页面的完整闭环

## 本地依赖

- Go 1.25.x
- Node.js 20+
- PostgreSQL 15+
- Redis 7+

## 常用命令

### 后端

```bash
cd backend
GOTOOLCHAIN=local GOSUMDB=off go build ./...
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server
GOTOOLCHAIN=local GOSUMDB=off go run ./cmd/server --migrate-only
```

### 前端

```bash
cd frontend
npm install
npm run dev
npm run build
npm run typecheck
```

## 代理说明

如果拉取依赖遇到网络问题，可使用以下代理：

```bash
export https_proxy=http://127.0.0.1:7890
export http_proxy=http://127.0.0.1:7890
export all_proxy=socks5://127.0.0.1:7890
```

## 配置来源

后端通过 Viper 读取：

- 根目录或父级目录中的 `config.yaml`
- 环境变量

默认开发通常直接使用环境变量即可。

## 模块开发流程

1. 定义 Ent schema
2. 生成 Ent 代码
3. 新增 `entity / repository / service / handler`
4. 注册路由和权限
5. 在前端新增 API 模块与页面
6. 在启动种子中补充菜单与演示数据

## 审计与权限

- 受保护接口统一走 JWT 鉴权中间件
- 权限点通过 `RequirePermission` 控制
- 非 GET 的受保护接口默认进入审计日志

## legacy 目录说明

以下目录不是当前 starter 的活跃代码：

- `_legacy_deploy`
- `_legacy_docs`
- `backend/_legacy_*`
- `frontend/_legacy_src`

这些内容只作为历史隔离区保留，不应继续向当前骨架回灌业务语义。
