# Soporte API

API REST del sistema de soporte técnico en Go con arquitectura hexagonal.

## Características

- Arquitectura hexagonal con separación de dominio, casos de uso, adaptadores e HTTP.
- Configuración por entorno (development, production).
- Docker de desarrollo y producción.
- Healthchecks, métricas Prometheus, middlewares de seguridad y CI.
- Migraciones SQL versionadas.
- Modo degradado si la base de datos no está disponible.

## Tecnologías del proyecto

- Go `1.25.5`: lenguaje principal.
- Gin: framework HTTP para rutas, handlers y middlewares.
- GORM: ORM para acceso a base de datos.
- PostgreSQL: base principal pensada para uso real.
- SQLite: soporte auxiliar para tests y desarrollo.
- Prometheus: métricas HTTP vía `/metrics`.
- Swagger UI / OpenAPI: documentación y prueba manual de endpoints.
- Docker: empaquetado y ejecución en contenedores.
- GitHub Actions: CI con tests, lint, seguridad y build.
- `golangci-lint`: control de calidad estática con 20+ linters.
- `gosec`: análisis de seguridad de código.
- `govulncheck`: detección de vulnerabilidades en dependencias.
- `trivy`: escaneo de vulnerabilidades en imágenes Docker.

## Librerías usadas

Librerías directas del proyecto:

- `github.com/gin-gonic/gin`: capa HTTP.
- `github.com/gin-contrib/gzip`: compresión gzip de respuestas.
- `gorm.io/gorm`: base del ORM.
- `gorm.io/driver/postgres`: driver PostgreSQL.
- `gorm.io/driver/sqlite`: driver SQLite para tests y build no productivo.
- `github.com/prometheus/client_golang`: métricas Prometheus.
- `github.com/swaggo/http-swagger`: Swagger UI en desarrollo.
- `github.com/go-playground/validator/v10`: validación de structs y campos.
- `go.uber.org/mock`: generación y uso de mocks en tests.

Librerías relevantes que aparecen por dependencia:

- `github.com/jackc/pgx/v5`: driver real usado por PostgreSQL.
- `github.com/mattn/go-sqlite3`: soporte SQLite.
- `github.com/prometheus/common`: soporte interno de métricas.
- `github.com/swaggo/files`: assets de Swagger.

## Arquitectura usada

La base sigue una idea hexagonal simple:

- `core`: reglas del negocio y contratos.
- `application`: casos de uso.
- `adapters`: implementaciones técnicas.
- `delivery`: entrada HTTP.
- `app`: ensamblaje de dependencias.

## Orden recomendado para crear un proyecto nuevo desde esta base

Si copias esta estructura para otro proyecto, el orden recomendado es este:

1. Definir entidades del dominio en `internal/core/domain`.
2. Definir puertos/interfaces en `internal/core/ports`.
3. Crear casos de uso en `internal/application/services`.
4. Crear modelos GORM en `internal/adapters/repository/models`.
5. Crear repositorios GORM en `internal/adapters/repository/repository`.
6. Crear contratos HTTP en `internal/delivery/http/contracts`.
7. Crear handlers HTTP en `internal/delivery/http/handlers`.
8. Registrar rutas en `internal/delivery/http/routes`.
9. Crear migraciones SQL en `db/migrations`.
10. Documentar endpoints en OpenAPI.
11. Validar con `go test`, `go vet`, lint y Docker.

## Qué hace cada carpeta

[`cmd/api`](cmd/api)

Punto de entrada de la API. Solo arranca configuración, logger y aplicación.

[`cmd/migrate`](cmd/migrate)

Punto de entrada del comando de migraciones SQL. Sirve para `up`, `down` y `status`.

[`internal/app`](internal/app)

Bootstrap de la aplicación. Conecta config, base de datos, router y `http.Server`.

[`internal/config`](internal/config)

Carga variables de entorno, valida configuración y resuelve archivos `.env`.

[`internal/core/domain`](internal/core/domain)

Entidades puras del negocio y errores de aplicación. Aquí no debería entrar Gin ni GORM.

[`internal/core/ports`](internal/core/ports)

Interfaces del dominio. Definen qué necesita la aplicación de sus adaptadores.

[`internal/application/services`](internal/application/services)

Casos de uso. Aquí está la lógica de negocio y validación principal.

[`internal/adapters/repository/database`](internal/adapters/repository/database)

Apertura, cierre y selección de dialector para base de datos.

[`internal/adapters/repository/models`](internal/adapters/repository/models)

Modelos GORM que representan tablas y relaciones de persistencia.

[`internal/adapters/repository/repository`](internal/adapters/repository/repository)

Implementaciones concretas de los puertos usando GORM.

[`internal/adapters/repository/migrations`](internal/adapters/repository/migrations)

Runner de migraciones SQL versionadas y soporte de `AutoMigrate`.

[`internal/delivery/http/contracts`](internal/delivery/http/contracts)

Requests, query params y responses HTTP. Aquí viven los DTOs del transporte.

[`internal/delivery/http/handlers`](internal/delivery/http/handlers)

Handlers REST. Hacen bind, llaman servicios y devuelven JSON.

[`internal/delivery/http/middlewares`](internal/delivery/http/middlewares)

Middlewares HTTP: request_id, recovery, gzip, CORS, security headers, size limit, content-type validation, rate limit, query timeout, metrics y logging estructurado.

[`internal/delivery/http/routes`](internal/delivery/http/routes)

Registro central de endpoints y conexión entre handlers y router Gin.

[`configs`](configs)

Archivos de entorno para local, Docker y ejemplo base.

[`build`](build)

Definición de imagen Docker.

[`deployments`](deployments)

Archivos de despliegue y `docker-compose`.

[`db/migrations`](db/migrations)

Migraciones SQL versionadas del esquema.

[`/.github/workflows`](.github/workflows)

Pipelines de CI.

## Configuración de entorno

Archivos disponibles:

- [`.env.local`](configs/.env.local): para correr desde tu máquina con `make run`.
- [`.env.development`](configs/.env.development): para correr en contenedor de desarrollo.
- [`.env.production`](configs/.env.production): configuración optimizada para producción.
- [`.env.example`](configs/.env.example): ejemplo base de referencia.
- [`.env`](.env): override local en la raíz si quieres sobreescribir valores.

Variables más importantes:

- `DB_ENABLED=true`: activa la base.
- `DB_DRIVER=postgres|sqlite`: selecciona driver.
- `DB_DSN=...`: conexión a base.
- `DB_FAIL_FAST=false`: deja que la API arranque en modo degradado si la base falla.
- `DB_QUERY_TIMEOUT=5s`: timeout para queries de base de datos.
- `DB_CONN_MAX_LIFETIME=30m`: tiempo máximo de vida de una conexión.
- `DB_CONN_MAX_IDLE_TIME=5m`: tiempo máximo que una conexión puede estar sin usarse.
- `DB_SLOW_QUERY_THRESHOLD=200ms`: umbral para loguear queries lentas.
- `DOCS_ENABLED=true|false`: activa o apaga docs.
- `SECURITY_CORS_ALLOW_ALL=true|false`: permite todos los orígenes o solo los de la lista.
- `SECURITY_CORS_ALLOWED_ORIGINS=https://app.ejemplo.com`: orígenes permitidos (separados por coma).
- `SECURITY_TRUSTED_PROXIES=`: IPs del Load Balancer (vacío = sin proxy, usa IP TCP directa).
- `SECURITY_RATE_LIMIT_REQUESTS=60`: requests por ventana.
- `SECURITY_RATE_LIMIT_WINDOW=1m`: ventana del rate limit.
- `SECURITY_REQUEST_SIZE_LIMIT_BYTES=1048576`: límite de body.

## Cómo ejecutar

Ejecutar en local:

```bash
make run
```

La API levanta por defecto en `http://localhost:8080`.

Ejecutar con Docker (Desarrollo):

```bash
make docker-up
```

Ejecutar con Docker (Producción):

```bash
make docker-up-prod
```

Ver logs:

```bash
make docker-logs        # Desarrollo
make docker-logs-prod   # Producción
```

Bajar el stack:

```bash
make docker-down        # Desarrollo
make docker-down-prod   # Producción
```

## Docker

[`build/Dockerfile.dev`](build/Dockerfile.dev) y [`build/Dockerfile.prod`](build/Dockerfile.prod)

Qué hacen:

- Build multi-stage para imágenes optimizadas.
- Desarrollo: Alpine con herramientas de debug y healthcheck.
- Producción: imagen `scratch` ultra-ligera (~15MB).
- Usuario non-root (UID 65532 en prod, UID 1000 en dev).
- Binario compilado con `-trimpath` y `-ldflags="-s -w"`.
- Build de producción sin SQLite ni Swagger (`-tags production`).

Seguridad en docker-compose:

- `security_opt: no-new-privileges` (previene escalación de privilegios).
- `cap_drop: ALL` (quita capabilities de Linux).
- `read_only: true` (filesystem de solo lectura).
- Límites de recursos (CPU y memoria).
- Healthchecks configurados.
- Logging controlado (max-size, max-file).

Comandos:

```bash
make docker-build-dev
make docker-build-prod
```

## Migraciones

Comandos disponibles:

```bash
make migrate-status
make migrate-up
make migrate-down
```

O directamente:

```bash
go run ./cmd/migrate status
go run ./cmd/migrate up
go run ./cmd/migrate down -steps=1
```

Qué hace el migrador:

- lee `*.up.sql` y `*.down.sql`
- crea la tabla `schema_migrations`
- aplica migraciones pendientes
- permite rollback por pasos
- muestra estado de migraciones aplicadas y pendientes

## Observabilidad y operación

Incluye:

- `GET /health`: estado general de la API.
- `GET /livez`: liveness probe para Kubernetes.
- `GET /readyz`: readiness probe (verifica conexión a DB).
- `GET /metrics`: métricas Prometheus (latencias, requests, etc.).
- Logging estructurado con `slog` (JSON en producción, texto en desarrollo).
- Logs sanitizados (no exponen datos sensibles).
- Shutdown ordenado con timeout configurable.
- Timeouts HTTP configurables (read, write, idle).
- Timeout de queries de base de datos.
- Modo degradado si la base no conecta (`DB_FAIL_FAST=false`).

## Seguridad HTTP incluida

Incluye:

- CORS configurable: modo wildcard o lista de orígenes permitidos con `Vary: Origin`.
- Trusted proxies configurables: evita spoofing de IP vía `X-Forwarded-For`.
- Rate limit por IP con token bucket: headers `RateLimit-Limit`, `RateLimit-Remaining`, `RateLimit-Reset` y `Retry-After` en 429.
- Límite de tamaño de request.
- Validación de Content-Type (solo `application/json`).
- Security headers: HSTS, CSP estricta, X-Frame-Options, X-Content-Type-Options, Permissions-Policy.
- Request ID único por petición.
- Recovery middleware con logs sanitizados (no expone detalles del panic al cliente).
- Compresión gzip automática.
- Timeout configurable de queries con respuesta 504.
- Mensajes de validación sin exponer nombres internos de structs.
- Logs sanitizados (no exponen query params sensibles).
- Errores de DB sanitizados (no exponen estructura de tablas).
- `context.DeadlineExceeded` en DB devuelve 503 en lugar de 500.

## Calidad automática

Incluye:

- `go test ./...` con race detection y cobertura.
- `go vet ./...`
- `golangci-lint` con 20+ linters de seguridad y calidad.
- `gosec` para análisis de seguridad de código.
- `govulncheck` para detectar vulnerabilidades en dependencias.
- `trivy` para escanear vulnerabilidades en imágenes Docker.

Pipeline de CI (`.github/workflows/ci.yml`):

1. **test**: ejecuta tests con race detection y cobertura.
2. **lint**: ejecuta golangci-lint.
3. **security**: ejecuta govulncheck y gosec.
4. **docker-build**: construye imagen y la escanea con trivy.

## Documentación HTTP

En producción, el build no incluye Swagger UI.

### URLs de documentación (solo en desarrollo con `DOCS_ENABLED=true`)

| URL | Descripción |
|-----|-------------|
| `http://localhost:8080/swagger/index.html` | Swagger UI — interfaz visual para explorar y probar endpoints |
| `http://localhost:8080/openapi.json` | Esquema OpenAPI en formato JSON |

## Orden recomendado para agregar un módulo nuevo

Para agregar un recurso nuevo, el orden recomendado es:

1. Crear entidad en `internal/core/domain`.
2. Crear puerto en `internal/core/ports`.
3. Crear tipos del servicio en `internal/application/services`.
4. Crear servicio/caso de uso en `internal/application/services`.
5. Crear modelo GORM en `internal/adapters/repository/models`.
6. Crear repositorio GORM en `internal/adapters/repository/repository`.
7. Crear DTOs en `internal/delivery/http/contracts`.
8. Crear handler HTTP en `internal/delivery/http/handlers`.
9. Registrar rutas en `internal/delivery/http/routes`.
10. Crear tests del servicio.
11. Crear tests del handler.
12. Agregar migración SQL.
13. Actualizar OpenAPI/Swagger.

## Pendiente

1. Autenticación y autorización (JWT, roles).
2. Rate limit distribuido con Redis (para múltiples instancias).
3. Gestión segura de secretos (vault, secrets manager).
4. Cache con Redis para respuestas frecuentes.
5. OpenTelemetry para tracing distribuido.
6. Audit log de operaciones (quién hizo qué y cuándo).
