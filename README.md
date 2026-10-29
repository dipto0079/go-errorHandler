# go-errorHandler

A small Go package that stores application errors in PostgreSQL. It also ships a gRPC server for the same operations.

## Install

```bash
go get github.com/dipto0079/go-errorHandler
```

## Usage

```go
import errorhandler "github.com/dipto0079/go-errorHandler"

h, err := errorhandler.New(
	"user=postgres password=secret host=localhost port=5432 dbname=mydb sslmode=disable",
	errorhandler.WithEnv("production"),
	errorhandler.WithService("api"),
)
if err != nil {
	log.Fatal(err)
}
defer h.Close()

if err := h.Migrate(ctx); err != nil { // creates/upgrades the table, safe to call every start
	log.Fatal(err)
}
```

### Record errors

```go
id, err := h.Capture(ctx, err)                              // stores message + stack trace
id, err := h.Capture(ctx, err, errorhandler.Code("DB_DOWN")) // custom code
```

The same error (same code, environment and message) that is still open is stored once and its
`OccurrenceCount` goes up. Once you `Resolve` it, the next occurrence opens a new record.
If the error has an `ErrorCode() string` method, that is used as the code.

### Record panics in a web server

```go
http.ListenAndServe(":8080", h.Middleware(mux)) // net/http, gorilla/mux, chi
```

### Notifications

Notifiers run on the first occurrence of an error only.

```go
errorhandler.New(dsn,
	errorhandler.WithNotifier(errorhandler.WebhookNotifier("https://hooks.slack.com/...")), // Slack/Discord/any JSON webhook
	errorhandler.WithNotifier(errorhandler.SMTPNotifier("smtp.example.com:587", auth, "from@x.com", "ops@x.com")),
	errorhandler.WithNotifier(func(ctx context.Context, e errorhandler.Error) error { return nil }), // your own
)
```

### Read and manage

```go
e, err := h.Get(ctx, id) // errors.Is(err, errorhandler.ErrNotFound) when missing
open, err := h.List(ctx, errorhandler.ListOptions{
	Status: errorhandler.StatusOpen, Env: "production", Since: time.Now().Add(-24 * time.Hour),
	Limit: 20, Offset: 0,
})
err = h.Resolve(ctx, id)
err = h.Delete(ctx, id, "production") // soft delete, hidden from List unless IncludeDeleted
```

| Method | Description |
| --- | --- |
| `Capture(ctx, err, ...CaptureOption)` | Record a Go error with stack trace, de-duplicated |
| `Middleware(next)` | `http.Handler` middleware that records panics and answers 500 |
| `Create(ctx, Error)` | Store a record exactly as given |
| `Get(ctx, id)` | Fetch one record |
| `List(ctx, ListOptions)` | Newest first; filter by env, status, code, time; paginate |
| `Resolve(ctx, id)` | Mark as resolved |
| `Delete(ctx, id, env)` | Soft delete, records the environment that deleted it |
| `Migrate(ctx)` | Create/upgrade the table |
| `Close()` | Close the database connection |

Already have a `*sql.DB`? Use `errorhandler.NewWithDB(db)` instead of `New`.
A runnable example is in [examples/http](examples/http).

## Database setup

`Migrate(ctx)` creates the table for you. The SQL is also available in [errorHandler/migrations/sql](errorHandler/migrations/sql). You can run it by hand or with the bundled goose runner:

```bash
cp errorHandler/env/sample.config errorHandler/env/config   # then edit the [database] section
cd errorHandler
go run ./migrations up
```

Other commands: `down`, `redo`, `status`, `version`, `create NAME sql`.

## gRPC server

```bash
go run ./errorHandler
```

It reads `errorHandler/env/config` and listens on `server.host:server.port`.

To regenerate the gRPC code from [gunk/v1/errorHandler/errorHandler.gunk](gunk/v1/errorHandler/errorHandler.gunk):

```bash
cd gunk && ./generate.sh
```

## Tests

```bash
cd errorHandler/storage/postgres
DATABASE_CONNECTION="user=postgres password=secret host=localhost port=5432 sslmode=disable" go test -v
```

## License

MIT, see [LICENSE](LICENSE).
