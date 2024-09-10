# go-sqlc

This project is for testing how `sqlc` works and how useful it is. It includes some utility functions that are worth adding when applied to actual projects.

- https://github.com/sqlc-dev/sqlc
- https://github.com/jackc/pgx

```bash
go generate ./...
go run .
```

## Impressions
- It might be useful for projects that are not only simple but also involve many complex queries (e.g., a statistics server).
- I was looking forward to performing unit tests with raw queries without writing any test code, but the macros (`sqlc.embed(table_name)`, etc.) are getting in the way.
- In projects where queries change frequently depending on application logic, such as in game servers, it seems less convenient compared to using a general query builder ORM.
