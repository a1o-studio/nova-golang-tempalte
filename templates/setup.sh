#!/bin/bash
# Setup script for initializing the Go service template
# 用于 degit 克隆后初始化项目

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 默认值
SERVICE_NAME=$(basename "$PWD")
MODULE_PATH="github.com/a1ostudio"

# 显示帮助
show_help() {
  echo "使用方法: ./setup.sh [选项]"
  echo ""
  echo "选项:"
  echo "  --name=<name>       服务名称 (默认: 当前目录名)"
  echo "  --module=<module>   Go module 路径 (默认: github.com/a1ostudio)"
  echo "  -h, --help          显示帮助信息"
  echo ""
  echo "示例:"
  echo "  ./setup.sh --name=my-api --module=github.com/myorg"
  echo ""
}

# 解析参数
for arg in "$@"; do
  case $arg in
    --name=*)
      SERVICE_NAME="${arg#*=}"
      shift
      ;;
    --module=*)
      MODULE_PATH="${arg#*=}"
      shift
      ;;
    -h|--help)
      show_help
      exit 0
      ;;
    *)
      echo -e "${RED}❌ 未知参数: $arg${NC}"
      show_help
      exit 1
      ;;
  esac
done

FULL_MODULE="${MODULE_PATH}/${SERVICE_NAME}"

echo -e "${GREEN}📦 开始初始化 Go 服务: ${SERVICE_NAME}${NC}"
echo -e "${GREEN}📝 Go Module: ${FULL_MODULE}${NC}"
echo ""

# 1. 更新 go.mod
echo "🔧 更新 go.mod..."
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  sed -i '' "s|module.*|module ${FULL_MODULE}|g" go.mod
else
  # Linux
  sed -i "s|module.*|module ${FULL_MODULE}|g" go.mod
fi
echo -e "${GREEN}✅ go.mod 更新完成${NC}"

# 2. 更新所有 .go 文件的 import 路径
echo "🔄 更新 Go 文件的 import 路径..."
if [[ "$OSTYPE" == "darwin"* ]]; then
  find . -name "*.go" -type f -exec sed -i '' "s|github.com/a1ostudio/nova|${FULL_MODULE}|g" {} +
else
  find . -name "*.go" -type f -exec sed -i "s|github.com/a1ostudio/nova|${FULL_MODULE}|g" {} +
fi
echo -e "${GREEN}✅ Import 路径更新完成${NC}"

# 3. 更新 Makefile
echo "⚙️  更新 Makefile..."
if [ -f "Makefile" ]; then
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|/nova?sslmode|/${SERVICE_NAME}?sslmode|g" Makefile
    sed -i '' "s|nova:nova@pg|${SERVICE_NAME}:${SERVICE_NAME}@pg|g" Makefile
  else
    sed -i "s|/nova?sslmode|/${SERVICE_NAME}?sslmode|g" Makefile
    sed -i "s|nova:nova@pg|${SERVICE_NAME}:${SERVICE_NAME}@pg|g" Makefile
  fi
  echo -e "${GREEN}✅ Makefile 更新完成${NC}"
fi

# 4. 生成 README.md
echo "📖 生成 README..."
cat > README.md << 'EOF'
# {{SERVICE_NAME}}

基于 [nova-golang-template](https://github.com/a1o-studio/nova-golang-tempalte) 创建的 Go 服务项目。

## 快速开始

```bash
# 1. 复制配置文件
cp app.env.example app.env
cp docker-compose.yml.example docker-compose.yml

# 2. 修改 app.env 和 docker-compose.yml 中的配置
# 根据实际需求修改数据库名称、密码等

# 3. 安装依赖
go mod tidy

# 4. 启动数据库和 Redis
make dcup

# 5. 运行数据库迁移
make migrateup

# 6. 启动开发服务器
make dev
```

## API 文档

访问: http://localhost:4000/swagger/index.html

## 常用命令

```bash
make swag         # 生成 Swagger 文档
make sqlc         # 生成 SQLC 代码
make test         # 运行测试
make fmt          # 格式化代码
make dcdown       # 停止 Docker 服务
```
EOF

# 替换 README 中的占位符
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' "s|{{SERVICE_NAME}}|${SERVICE_NAME}|g" README.md
else
  sed -i "s|{{SERVICE_NAME}}|${SERVICE_NAME}|g" README.md
fi

echo -e "${GREEN}✅ README.md 生成完成${NC}"

echo ""
echo -e "${GREEN}✅ 初始化完成！${NC}"
echo ""
echo "📋 后续步骤:"
echo "  1. cp app.env.example app.env"
echo "  2. cp docker-compose.yml.example docker-compose.yml"
echo "  3. 修改 app.env 和 docker-compose.yml 配置"
echo "  4. go mod tidy        # 整理依赖"
echo "  5. make dcup          # 启动数据库和 Redis"
echo "  6. make migrateup     # 运行数据库迁移"
echo "  7. make dev           # 启动开发服务器"
echo ""
echo "📚 API 文档: http://localhost:4000/swagger/index.html"
echo ""
