# liveplan_backend_go

Proyecto backend en Go para Liveplan. Contiene endpoints básicos y conexión a Postgres.

Endpoints:
- GET /health -> 200 OK
- GET /users -> lista usuarios
- POST /users -> crea usuario {"name":"...","email":"..."}

Credenciales Postgres están hardcodeadas en `internal/db/db.go`.

Para correr:

1. Asegúrate de tener Go instalado (>=1.20).
2. Ten una base Postgres corriendo en las credenciales por defecto en `internal/db/db.go`.
3. Desde la carpeta del proyecto:

```bash
cd liveplan_backend_go
go build
./liveplan_backend_go
```

Usar Docker Compose (levanta Postgres para desarrollo):

```bash
cd liveplan_backend_go
docker compose up -d
# Espera a que Postgres esté listo (puedes revisar logs con: docker compose logs -f db)
go build
./liveplan_backend_go
```
# Liveplan_Go
