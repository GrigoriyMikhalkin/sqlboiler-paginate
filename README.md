# **sqlboiler-paginate**

## Description
Pagination package for [sqlboiler](https://github.com/volatiletech/sqlboiler). Supports `sqlboiler` versions 3/4. For example, to install this lib for `sqlboiler v4`, run:
```bash
go get github.com/grigoriymikhalkin/sqlboiler-paginate/v4
```

Package consists of:
  - `PaginatorParamsParser` -- default implementation parses query params.
  - `PaginationQueryMods` function -- based on provided `PaginatorParams` produces slice of [QueryMods]().

`PaginationQueryMods` currently implements two pagination methods: [offset pagination](https://developer.box.com/guides/api-calls/pagination/offset-based/) and [keyset pagination](https://use-the-index-luke.com/no-offset)

## Usage

### PaginatorParams parser
Package provides customizable implementation of `PaginatorParams` parser.

### Offset Pagination
If `Offset` field in `PaginatorParams`> 0, [offset pagination](https://developer.box.com/guides/api-calls/pagination/offset-based/) is used:
```go

```

### Keyset pagination
Otherwise, `PaginationQueryMods` uses [keyset pagination](https://use-the-index-luke.com/no-offset):
```go
```

### Examples
```go
```

[simple example]()

## License
[MIT](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/LICENSE)
