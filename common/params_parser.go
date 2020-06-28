package common

import (
	"fmt"
  "encoding/json"
  "net/url"
  "strings"
  "strconv"
)

const (
  defaultLimitParam = "limit"
  defaultOffsetParam = "offset"
  defaultOrderByParam = "order_by"
  defaultPrevPageValuesParam = "prev_page_values"
)

type QueryParams = map[string][]string

type OrderByParseFunc = func(queryParams QueryParams, orderByParam, prevPageValueParam string) (orderBy OrderByParams, err error)

type PaginatorParamsParser interface {
  ParseQuery(query string) (params PaginatorParams, err error)
}

type defaultPaginatorParamsParser struct {
  LimitParam          string
  OffsetParam         string
  OrderByParam        string
  PrevPageValuesParam string

  OrderByParse OrderByParseFunc
}

func NewPaginatorParamsParser(limitParam, offsetParam, orderByParam, prevPageValuesParam string, orderByParse OrderByParseFunc) *defaultPaginatorParamsParser {
  parseFunc := DefaultOrderByParse
  if orderByParse != nil {
    parseFunc = orderByParse
  }

  return &defaultPaginatorParamsParser{
    LimitParam: limitParam,
    OffsetParam: offsetParam,
    OrderByParam: orderByParam,
    PrevPageValuesParam: prevPageValuesParam,
    OrderByParse: parseFunc,
  }
}

func NewDefaultPaginatorParamsParser() *defaultPaginatorParamsParser {
  return &defaultPaginatorParamsParser{
    LimitParam: defaultLimitParam,
    OffsetParam: defaultOffsetParam,
    OrderByParam: defaultOrderByParam,
    PrevPageValuesParam: defaultPrevPageValuesParam,
    OrderByParse: DefaultOrderByParse,
  }
}

func (p *defaultPaginatorParamsParser) ParseQuery(query string) (params PaginatorParams, err error) {
  paramsQuery := query
  if strings.Contains(paramsQuery, "?") {
    paramsQuery = strings.SplitN(paramsQuery, "?", 2)[1]
  }
  values, err := url.ParseQuery(paramsQuery)
  if err != nil {
    return
  }

  // parsing limit
  limit := -1
  if len(values[p.LimitParam]) > 0 {
    limit, err = strconv.Atoi(values[p.LimitParam][0])
    if err != nil {
      return
    }
  }

  // parsing offset
  offset := -1
  if len(values[p.OffsetParam]) > 0 {
    offset, err = strconv.Atoi(values[p.OffsetParam][0])
    if err != nil {
      return
    }
  }

  // parsing order by params
  orderBy, err := p.OrderByParse(values, p.OrderByParam, p.PrevPageValuesParam)
  if err != nil {
    return
  }

  return PaginatorParams{
    Limit: limit,
    Offset: offset,
    OrderBy: orderBy,
  }, nil
}

func DefaultOrderByParse(queryParams QueryParams, orderByParam, prevPageValueParam string) (orderByParams OrderByParams, err error) {
  orderByFields := []string{}
  if len(queryParams[orderByParam]) > 0 {
    orderByFields = strings.Split(queryParams[orderByParam][0], ",")
  } else {
    return
  }

  prevPageValues := map[string]interface{}{}
  if len(queryParams[prevPageValueParam]) > 0 {
    err = json.Unmarshal([]byte(queryParams[prevPageValueParam][0]), &prevPageValues)
    if err != nil {
      return
    }
  }

  // filter order params
  orderByParams = OrderByParams{}
  for _, orderByField := range orderByFields {
    x := strings.Split(orderByField, "-")
    field := x[0]
    order := "asc"
    if len(x) > 1 {
      order = x[1]
    }

    orderByParams[field] = &OrderByParam{
      Order: order,
      LastValue: fmt.Sprintf("%v", prevPageValues[field]),
    }
  }

  return
}
