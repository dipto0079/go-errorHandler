# go-errorHandler

```bash
go get github.com/dipto0079/go-errorHandler
```

```go
import errorhandler "github.com/dipto0079/go-errorHandler"

h, err := errorhandler.New("user=postgres password=... host=localhost port=5432 dbname=mydb sslmode=disable")
id, err := h.Create(ctx, errorhandler.Error{ErrorCode: "E100", ErrorDetails: "boom", EnvType: "production", CreatedBy: "api"})
```

Create the table first with `errorHandler/migrations/sql/*.sql`.

##gunk get
<hr/>
https://github.com/gunk/gunk
<hr/>
##generate
<hr/>
./generate.sh
<hr/>
##goose add
<hr/>
go get github.com/pressly/goose
<hr/>

##migrations
<hr/>
go run migrations/migrate.go create create_Name_tabls sql
<hr/>
go run migrations/migrate.go up
<hr/>
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version


##test
<hr/>
<hr/>
cd errorHandler/storage/postgres
<hr/>
DATABASE_CONNECTION="user=postgres password=dipto host=localhost port=5432 sslmode=disable" go test -v