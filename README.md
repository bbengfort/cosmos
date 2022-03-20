# Cosmos

**An API for an expansive space strategy and battle game.**

## Database

The database schema is migrated using [go-migrate](https://github.com/golang-migrate/migrate) and managed in the application using [go-bindata](https://github.com/kevinburke/go-bindata). Both of these command line tools must be installed in order to create and manage migrations.

To create a new migration:

```
$ migrate create -ext sql -dir pkg/db/migrations/ -seq "description_of_migration"
```

This should create two files in `pkg/db/migrations`: `0000N_description_of_migration.up.sql` and `0000N_description_of_migration.down.sql` which are used to apply SQL commands when migrating up and down versions respectively. Add the SQL to these files that is required for the migration, then generate the bindata as follows:

```
$ go generate ./pkg/db/...
```

This will create a `migrations.go` file in `pkg/db/schema` that contains the migration data. The migration can be applied using the `cosmos` binary as follows:

```
$ go run ./cmd/cosmos migrate
```

You can check the current version of the database schema with the `cosmos` binary as well:

```
$ go run ./cmd/cosmos schema
```