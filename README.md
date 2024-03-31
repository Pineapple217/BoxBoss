# BoxBoss

## Dev Setup

### Dependencies

Latest version of Go and the following codegen tools.

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/a-h/templ/cmd/templ@latest
```

Air is optional but strongly recommended.

```sh
go install github.com/cosmtrek/air@latest
```

Code gen

```sh
sqlc generate
templ generate
```
