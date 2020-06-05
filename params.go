package paginate

import (
  "encoding/json"
  "net/url"
  "strings"
  "strconv"
)

const (
  defaultLimitParam = "limit"
  defaultOrderByParam = "order_by"
  defaultPrevPageValuesParam = "prev_page_values"
)

type PaginatorParams struct {
  Limit   int
  Offset  int
  OrderBy map[string]string
}

type PaginatorParamsParser interface {
  ParseQuery(query string) (params *PaginatorParams, err error)
}

type OrderByParseFunc = func(queryParams map[string][]string, orderByParam, prevPageValueParam string) (orderBy map[string]string, err error)

type paginatorParamsParserImpl struct {
  LimitParam          string
  OffsetParam         string
  OrderByParam        string
  PrevPageValuesParam string

  OrderByParse OrderByParseFunc
}

func NewPaginatorParamsParser(limitParam, orderByParam, prevPageValuesParam string, orderByParse OrderByParseFunc) *paginatorParamsParserImpl {
  parseFunc := DefaultOrderByParse
  if orderByParse != nil {
    parseFunc = orderByParse
  }

  return &paginatorParamsParserImpl{
    LimitParam: limitParam,
    OrderByParam: orderByParam,
    PrevPageValuesParam: prevPageValuesParam,
    OrderByParse: parseFunc,
  }
}

func NewDefaultPaginatorParamsParser() *paginatorParamsParserImpl {
  return &paginatorParamsParserImpl{
    LimitParam: defaultLimitParam,
    OrderByParam: defaultOrderByParam,
    PrevPageValuesParam: defaultPrevPageValuesParam,
    OrderByParse: DefaultOrderByParse,
  }
}

func (p *paginatorParamsParserImpl) ParseQuery(query string) (params *PaginatorParams, err error) {
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

  // parsing prev page values
  orderBy, err := p.OrderByParse(values, p.OrderByParam, p.PrevPageValuesParam)
  if err != nil {
    return
  }

  return &PaginatorParams{
    Limit: limit,
    OrderBy: orderBy,
  }, nil
}

func DefaultOrderByParse(queryParams map[string][]string, orderByParam, prevPageValueParam string) (orderByValues map[string]string, err error) {
  orderBy := []string{}
  if len(queryParams[orderByParam]) > 0 {
    orderBy = strings.Split(queryParams[orderByParam][0], ",")
  }

  prevPageValues := map[string]string{}
  if len(queryParams[prevPageValueParam]) > 0 {
    err := json.Unmarshal([]byte(queryParams[prevPageValueParam][0]), &prevPageValues)
    if err != nil {
      return nil, err
    }
  }

  // filter order params
  orderByValues = map[string]string{}
  for _, orderParam := range orderBy {
    parts := strings.Split(orderParam, "-")

    if len(parts) > 1 {
      orderByValues[strings.Join(parts, " ")] = prevPageValues[parts[0]]
    } else {
      orderByValues[parts[0]] = prevPageValues[parts[0]]
    }
  }

  return
}
