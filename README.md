# Nova Golang Template

一个功能完整、开箱即用的 Go 服务模板项目，集成了最佳实践和常用工具链。

## 特性

- ✅ **Gin** - 高性能 HTTP 框架
- ✅ **PostgreSQL + SQLC** - 类型安全的数据库操作
- ✅ **Redis** - 缓存和分布式锁
- ✅ **Swagger** - 自动生成 API 文档
- ✅ **JWT & Paseto** - 多种认证方式
- ✅ **Zap** - 结构化日志
- ✅ **Migrate** - 数据库迁移管理
- ✅ **Docker Compose** - 一键启动开发环境
- ✅ **Makefile** - 统一的命令管理

## 快速开始

### 使用 degit 克隆模板

直接从 GitHub 克隆模板到本地，不包含 git 历史：

```bash
# 1. 使用 degit 克隆模板
npx degit a1o-studio/nova-golang-tempalte/templates my-service

# 2. 进入目录
cd my-service

# 3. 运行初始化脚本（自动替换模块路径和服务名）
chmod +x setup.sh
./setup.sh --name=my-service --module=github.com/yourorg

# 4. 配置环境文件和 Docker Compose（需要手动创建和修改）
cp app.env.example app.env
cp docker-compose.yml.example docker-compose.yml
# 编辑 app.env 和 docker-compose.yml，修改测试数据库与容器配置
# 修改 Makefile DB_URL

# 5. 初始化依赖
go mod tidy

# 6. 启动开发环境
make dcup       # 启动 PostgreSQL 和 Redis
make migrateup  # 运行数据库迁移
make dev        # 启动服务
```

### setup.sh 参数说明

```bash
./setup.sh [选项]

选项:
  --name=<name>       服务名称（默认：当前目录名）
  --module=<module>   Go module 路径（默认：github.com/a1ostudio）
  -h, --help          显示帮助信息

示例:
  ./setup.sh --name=user-api --module=github.com/mycompany
```

**注意**: `setup.sh` 只会自动替换代码中的模块路径和服务名称，不会创建配置文件。你需要手动从 `.example` 文件复制并修改配置。

### 使用示例

```bash
# 基本用法
npm run new:go user-api

# 指定目标路径
npm run new:go user-api --path=apps/user-api

# 自定义 module 路径
npm run new:go user-api --module=github.com/mycompany
## setup.sh 做了什么？
```

初始化脚本会自动完成以下配置：

1. ✅ 更新 `go.mod` 中的 module 名称
2. ✅ 递归替换所有 `.go` 文件中的 import 路径
3. ✅ 创建 `app.env` 并替换服务名
4. ✅ 更新 `app.env.example` 中的默认值
5. ✅ 更新 `Makefile` 中的数据库连接配置
6. ✅ 更新 `docker-compose.yml` 中的服务名和容器名
7. ✅ 生成项目专属的 `README.md`

## 项目结构

``` plaintext
templates/                    # Go 服务模板
├── cmd/
│   └── main.go              # 应用入口
├── db/
│   ├── migration/           # 数据库迁移
│   ├── query/               # SQLC 查询
│   └── sqlc/                # 生成的代码
├── internal/
│   ├── config/              # 配置管理
│   ├── controller/          # HTTP 控制器
│   ├── middleware/          # 中间件
│   ├── model/               # 数据模型
│   ├── pkg/                 # 工具包
│   ├── server/              # HTTP 服务器
│   └── service/             # 业务逻辑
├── app.env.example          # 环境变量示例
├── docker-compose.yml       # Docker 配置
├── Makefile                 # Make 命令
├── sqlc.yaml                # SQLC 配置
├── setup.sh                 # 初始化脚本
└── go.mod                   # Go 模块
```

## 模板特性详解

### 数据库管理

- **SQLC** - 从 SQL 生成类型安全的 Go 代码
- **Migrate** - 数据库版本管理和迁移

```bash
# 创建迁移
make migratecreate name=add_users_table

# 执行迁移
make migrateup

# 回滚迁移
make migratedown step=1

# 生成 SQLC 代码
make sqlc
```

### API 文档

- **Swagger** - 自动生成 API 文档

```bash
# 生成 Swagger 文档
make swag

# 访问文档
# http://localhost:4000/swagger/index.html
```

### 开发工具

- **Docker Compose** - 本地开发环境
- **Makefile** - 统一命令接口

```bash
# 启动开发环境
make dcup

# 停止开发环境
make dcdown

# 运行测试
make test

# 格式化代码
make fmt
```

### 认证方式

模板内置两种 Token 方式：

- **JWT** - JSON Web Token
- **Paseto** - Platform-Agnostic Security Tokens

### 日志系统

- **Zap** - 高性能结构化日志
- 支持日志文件轮转
- 开发/生产环境不同配置

### 中间件

- CORS 跨域处理
- Rate Limiting 限流
- Recovery 恢复
- Timeout 超时控制

## 环境要求

- Go 1.25.3+
- PostgreSQL 17+
- Redis 8+
- Docker & Docker Compose (可选)

## 常见问题

### Q: 如何修改默认端口？

A: 编辑 `app.env` 文件中的 `APP_PORT` 变量。

### Q: 如何添加新的数据库表？

A:

1. 创建迁移文件：`make migratecreate name=add_xxx_table`
2. 编辑生成的 SQL 文件
3. 运行迁移：`make migrateup`
4. 在 `db/query` 添加查询
5. 生成代码：`make sqlc`

### Q: 如何更换 Go module 路径？

A: 使用 `setup.sh` 脚本时指定 `--module` 参数：

```bash
./setup.sh --name=my-service --module=github.com/yourcompany
```

或手动修改：

1. 编辑 `go.mod` 中的 `module` 行
2. 全局替换所有 `.go` 文件中的 import 路径
3. 运行 `go mod tidy`

## 贡献

欢迎提交 Issue 和 Pull Request！

## License

MIT
