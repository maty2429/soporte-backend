# Variables
APP_NAME=soporte-api
CMD_DIR=./cmd/api
BIN_DIR=./bin
DOCKER_DIR=build
CONFIGS_DIR=configs

# Phony targets
.PHONY: help run test fmt tidy vet lint generate clean build-dev build-prod docker-build-dev docker-build-prod docker-up docker-down docker-up-prod docker-down-prod docker-logs-prod run-prod migrate-up migrate-down migrate-status

# Default target
help:
	@echo "Comandos disponibles:"
	@echo "  run               - Ejecutar la aplicación en modo desarrollo"
	@echo "  test              - Ejecutar todos los tests del proyecto"
	@echo "  fmt               - Formatear el código Go"
	@echo "  tidy              - Limpiar dependencias del módulo Go"
	@echo "  vet               - Ejecutar análisis estático (go vet)"
	@echo "  lint              - Ejecutar linter (golangci-lint o go vet como fallback)"
	@echo "  clean             - Eliminar binarios y archivos temporales"
	@echo "  build-dev         - Construir binario para desarrollo"
	@echo "  build-prod        - Construir binario optimizado para producción"
	@echo "  docker-build-dev  - Construir imagen Docker de desarrollo"
	@echo "  docker-build-prod - Construir imagen Docker de producción"
	@echo "  docker-up         - Levantar servicios con Docker Compose (dev)"
	@echo "  docker-down       - Detener servicios de Docker Compose"
	@echo "  docker-up-prod    - Levantar servicios con Docker Compose (prod)"
	@echo "  docker-down-prod  - Detener servicios de Docker Compose (prod)"
	@echo "  docker-logs-prod  - Ver logs de producción"
	@echo "  run-prod          - Ejecutar la aplicación en modo producción (local)"
	@echo "  migrate-up        - Ejecutar migraciones hacia adelante"
	@echo "  migrate-down      - Revertir última migración"
	@echo "  generate          - Regenerar mocks de los ports (go:generate)"
	@echo "  migrate-status    - Ver estado de las migraciones"

# Desarrollo Local
run:
	@echo "Iniciando aplicación..."
	go run $(CMD_DIR)

test:
	@echo "Ejecutando tests..."
	go test -race -v ./...

fmt:
	@echo "Formateando código..."
	go fmt ./...

tidy:
	@echo "Limpiando dependencias..."
	go mod tidy

vet:
	@echo "Analizando código (vet)..."
	go vet ./...

lint:
	@echo "Ejecutando linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint no está instalado. Usando go vet..."; \
		go vet ./...; \
	fi

generate:
	@echo "Regenerando mocks..."
	go generate ./internal/core/ports/...

clean:
	@echo "Limpiando binarios..."
	rm -rf $(BIN_DIR)

# Compilación
build-dev:
	@echo "Construyendo binario de desarrollo..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

build-prod:
	@echo "Construyendo binario de producción..."
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 go build -ldflags="-s -w" -tags production -o $(BIN_DIR)/$(APP_NAME)-prod $(CMD_DIR)

# Docker
docker-build-dev:
	@echo "Construyendo imagen Docker (dev)..."
	docker build -f $(DOCKER_DIR)/Dockerfile.dev -t $(APP_NAME):dev .

docker-build-prod:
	@echo "Construyendo imagen Docker (prod)..."
	docker build -f $(DOCKER_DIR)/Dockerfile.prod -t $(APP_NAME):prod .

docker-up:
	@echo "Levantando infraestructura (dev)..."
	docker compose -f deployments/docker-compose.yml up --build -d

docker-down:
	@echo "Bajando infraestructura..."
	docker compose -f deployments/docker-compose.yml down

docker-up-prod:
	@echo "Levantando infraestructura (prod)..."
	docker compose -f deployments/docker-compose.prod.yml up --build -d

docker-down-prod:
	@echo "Bajando infraestructura (prod)..."
	docker compose -f deployments/docker-compose.prod.yml down

docker-logs-prod:
	docker compose -f deployments/docker-compose.prod.yml logs -f

run-prod:
	@echo "Iniciando aplicación en modo producción..."
	go run -tags production $(CMD_DIR)

# Migraciones
migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-status:
	go run ./cmd/migrate status
