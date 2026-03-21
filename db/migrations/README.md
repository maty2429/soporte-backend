# Migrations

Este directorio contiene las migraciones SQL versionadas.

Formato esperado:

- `000001_nombre.up.sql`
- `000001_nombre.down.sql`

Comandos disponibles desde la raiz del proyecto:

- `make migrate-up`
- `make migrate-down`
- `make migrate-status`

Tambien puedes usar el comando directamente:

- `go run ./cmd/migrate up`
- `go run ./cmd/migrate down -steps=1`
- `go run ./cmd/migrate status`
