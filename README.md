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
Package provides customizable implementation of `PaginatorParams` parser. Function `NewPaginatorParamsParser(limitParam, offsetParam, orderByParam, prevPageValuesParam string, orderByParse OrderByParseFunc)` initializes parser with passed parameters:

 * limitParam -- name of limit parameter in query, by default `limit`
 * offsetParam -- name of offset parameter in query, by default `offset`
 * orderByParam --  name of ordering parameter in query, by default `order_by`
 * prevPageValuesParam -- previous values parameter, by default `prev_page_values`
 * orderByParse -- function that parses ordering parameter and last page values, by default function `DefaultOrderByParse`

If you want to use default values, you can initialize parser by calling function `paginate.NewDefaultPaginatorParamsParser()`.

Parser provides single function, to parse url query:
```go
ParseQuery(query string) (params PaginatorParams, err error)
```
Usage example:
```go
params, _ := parser.ParseQuery(`/users?limit=5&order_by=id&prev_page_values={"id":5}`)
> params
PaginatorParams {
  Limit: 5,
  Offset: -1,
  OrderBy: OrderByParams{
    "id": OrderByParam{
      Order: "asc",
      LastValue: "5",
    }
  }
}
```

Expected types:

* limit -- number
* offset -- number
* order_by -- string of format(`field-order,field2-order,etc`)
* prev_page_values -- urlencoded JSON object

Note: all previous page values converted to strings.

### Offset Pagination
If `Offset` field in `PaginatorParams`> 0, [offset pagination](https://developer.box.com/guides/api-calls/pagination/offset-based/) is used.

### Keyset pagination
Otherwise, `PaginationQueryMods` uses [keyset pagination](https://use-the-index-luke.com/no-offset).

### Examples
```go
import (
  ...
  "github.com/grigoriymikhalkin/sqlboiler-paginate/v4"

  "your-project/sqlboiler-models"
)
func paginate(query) []*models.User {
  parser := paginate.NewDefaultPaginatorParamsParser()
  params, err := parser.ParseQuery(query)
  if err != nil {
    panic(err)
  }
  mods := paginate.PaginationQueryMods(params)
  users, err := models.Users(mods...).All(...)
  return users
}
```

[simple example]()

## License
[MIT](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/LICENSE)
