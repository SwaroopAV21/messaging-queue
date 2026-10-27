# messaging-queue
Messaging queue using FIFO and to act as Kafka,

## PostgreSQL

The API uses Go's `database/sql` package with the pgx PostgreSQL driver. Set
these environment variables using your environment or secret manager; do not
commit credentials:

- `DB_URL`: PostgreSQL server hostname (for example, `localhost`).
- `DB_PORT`: PostgreSQL server port (for example, `5432`).
- `POSTGRES_USER` and `POSTGRES_PASS`: PostgreSQL credentials.
- `DB_NAME` and `POSTGRES_DB`: the same target database name. If both are set,
  they must match; either variable can provide the name.

`repository.InitDB` connects to the standard `postgres` maintenance database,
creates the target database if missing, connects to it, and creates the
`responses` table if needed. The resulting handle is available as
`repository.DB`; the application entry point calls `InitDB` once at startup.

After mapping a profile to a response, persist it with:

```go
response, ok := (&models.Response{}).MaptoResponse(&profile, request)
if !ok {
    return errors.New("profile response is nil")
}
if err := repository.SaveToDB(ctx, response); err != nil {
    return err
}
```

`InitializeSchema` creates the `responses` table if needed. The profile payload
is stored as PostgreSQL `JSONB`.

## GET proxy

Set `BASE_URL` to the base URL of the service that provides project,
experience, and skill data. For local development this can be
`http://localhost:8080`. In Docker Compose, use the upstream service's Compose
service name instead of `localhost` (for example, `http://profile-service:8080`).

The HTTP server listens on port `8081` and proxies these GET routes:

- `/api/projects`
- `/api/projects/featured`
- `/api/projects/{id}` (positive integer IDs)
- `/api/experiences`
- `/api/skill`

## Set this variables


```bash
# Linux:
    $env:BASE_URL = "http://localhost:8080"
    $env:DB_URL = "localhost"
    $env:DB_PORT = "5432"
    $env:DB_NAME = "messaging_queue"
    $env:POSTGRES_DB = "messaging_queue"
    $env:POSTGRES_USER = "your_postgres_user"
    $env:POSTGRES_PASS = "your_postgres_password"

#Windows:
    export BASE_URL = "http://localhost:8080"
    export DB_URL = "localhost"
    export DB_PORT = "5432"
    export DB_NAME = "messaging_queue"
    export POSTGRES_DB = "messaging_queue"
    export POSTGRES_USER = "your_postgres_user"
    export POSTGRES_PASS = "your_postgres_password"
---
```bash