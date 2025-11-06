.PHONY: swag dev sqlc test fmt lint

DB_URL=postgresql://{{username}}:{{password}}@localhost:5434/{{db_name}}?sslmode=disable
MIGRATE_DIR=db/migration
MIGRATE_CMD=migrate -path $(MIGRATE_DIR) -database "$(DB_URL)" -verbose

# 生成 Swagger 文档
swag:
	@echo "Generating swagger docs..."
	@swag fmt
	@swag init -g cmd/main.go -o ./docs --parseDependency

# 启动 Dev
dev: swag
	@echo "Starting dev server..."
	@go run cmd/main.go

migrateup:
	@echo "Running migrations up..."
	@if [ -n "$(step)" ]; then \
		$(MIGRATE_CMD) up $(step); \
	else \
		$(MIGRATE_CMD) up; \
	fi

migratedown:
	@echo "Running migrations down..."
	@if [ -n "$(step)" ]; then \
		$(MIGRATE_CMD) down $(step); \
	else \
		$(MIGRATE_CMD) down; \
	fi

## Create a new migration: make migratecreate name=create_users_table
migratecreate:
	@if [ -n "$(name)" ]; then \
		migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name); \
	else \
		echo "Please provide migration name: make migratecreate name=add_users"; \
	fi

sqlc:
	@echo "Running sqlc generate..."
	@sqlc generate
	@echo "Sqlc generate done."

dcup:
	@echo "Starting docker-compose..."
	@docker-compose --env-file app.env up -d
	@echo "docker-compose started."

dcdown:
	@echo "Stopping docker-compose..."
	@docker-compose --env-file app.env down
	@echo "docker-compose stopped."

test:
	@echo "Running tests..."
	@go test -v -cover -short ./...
	@echo "Tests completed."

fmt:
	@go fmt ./...
