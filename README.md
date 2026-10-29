# go-errorHandler

A small Go package that stores application errors in PostgreSQL. It also ships a gRPC server for the same operations.

## Install

```bash
go get github.com/dipto0079/go-errorHandler
```

## Usage

```go
import errorhandler "github.com/dipto0079/go-errorHandler"

h, err := errorhandler.New("user=postgres password=secret host=localhost port=5432 dbname=mydb sslmode=disable")
if err != nil {
	log.Fatal(err)
}
defer h.Close()

if err := h.Migrate(ctx); err != nil { // creates the table if missing, safe to call every start
	log.Fatal(err)
}

id, err := h.Create(ctx, errorhandler.Error{
	ErrorCode:    "E100",
	ErrorDetails: "something went wrong",
	EnvType:      "production",
	CreatedBy:    "api",
})

e, err := h.Get(ctx, id) // errors.Is(err, errorhandler.ErrNotFound) when missing
all, err := h.List(ctx)
err = h.Delete(ctx, id, "production") // soft delete
```

| Method | Description |
| --- | --- |
| `Create(ctx, Error)` | Store an error, returns its id |
| `Get(ctx, id)` | Fetch one error |
| `List(ctx)` | List errors |
| `Delete(ctx, id, env)` | Soft delete, records the environment that deleted it |
| `Migrate(ctx)` | Create the table if it does not exist |
| `Close()` | Close the database connection |

Already have a `*sql.DB`? Use `errorhandler.NewWithDB(db)` instead of `New`.

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
