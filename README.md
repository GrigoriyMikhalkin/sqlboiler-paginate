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
 * orderByParse -- function [OrderByParseFunc](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/common/params_parser.go#L19) that parses ordering parameters, by default [DefaultOrderByParse](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/common/params_parser.go#L113) function is used.
 * prevPageValuesParse -- function [PrevPageValuesParseFunc](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/common/params_parser.go#L21), that will parse last page values, by default [DefaultPrevPageValuesParse](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/common/params_parser.go#L138) function is used.
 
If you want to use default values, you can initialize parser by calling function `paginate.NewDefaultPaginatorParamsParser()`.

Parser provides single function, to parse url query:
```go
ParseQuery(query string) (params PaginatorParams, err error)
```
Usage example:
```go
params, _ := parser.ParseQuery(`/users?limit=5&order_by=id,name-desc&prev_page_values={"id":5}`)
> params
PaginatorParams {
  Limit: 5,
  Offset: -1,
  OrderBy: []*OrderByParam{
    { Field: "id", Order: "asc" },
    { Field: "name", Order: "desc" }
  },
  PrevPage: map[string]interface{}{
    "id": 5,  
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
If previous values for order columns unset, [offset pagination](https://developer.box.com/guides/api-calls/pagination/offset-based/) is used.
Example request:
```http request
GET /users?limit=5&offset=5
```

### Keyset pagination
If previous values for order columns are set, `PaginationQueryMods` uses [keyset pagination](https://use-the-index-luke.com/no-offset).
Example request:
```http request
GET /users?limit=5&order_by=id-asc,name-desc&prev_page_values={%22id%22:5}
```


### Examples
```go
import (
  ...
  "github.com/grigoriymikhalkin/sqlboiler-paginate/common"
  "github.com/grigoriymikhalkin/sqlboiler-paginate/v4"

  "your-project/sqlboiler-models"
)
func paginate(query) []*models.User {
  parser := common.NewDefaultPaginatorParamsParser()
  params, err := parser.ParseQuery(query)
  if err != nil {
    panic(err)
  }
  mods := paginate.PaginationQueryMods(params)
  users, err := models.Users(mods...).All(...)
  return users
}
```

[net/http example](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/examples/simple_net_example.go)

[gofiber example]()

## License
[MIT](https://github.com/GrigoriyMikhalkin/sqlboiler-paginate/blob/master/LICENSE)
