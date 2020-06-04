package main

import (
  "encoding/json"
  "net/url"
  "strings"
  "strconv"
)

const (
  defaultLimitParam = "limit"
  defaultOrderByParam = "order_by"
  defaultPrevPageValuesParamPrefix = "prev_page_values"
)

type PaginatorParams struct {
  Limit          int
  OrderBy        []string
  PrevPageValues map[string]string
}

type PaginatorParamsParser interface {
  ParseQuery(query string) (params *PaginatorParams, err error)
}

type PrevPageValuesParseFunc = func(queryParams map[string][]string, paramPrefix string) (prevPageValues map[string]string, err error)

type paginatorParamsParserImpl struct {
  LimitParam          string
  OrderByParam        string
  PrevPageValuesParamPrefix string

  PrevPageValuesParse PrevPageValuesParseFunc
}

func NewPaginatorParamsParser(limitParam, orderByParam, prevPageValuesParamPrefix string, prevPageValuesParse PrevPageValuesParseFunc) *paginatorParamsParserImpl {
  parseFunc := DefaultPrevPageValuesParse
  if prevPageValuesParse != nil {
    parseFunc = prevPageValuesParse
  }

  return &paginatorParamsParserImpl{
    LimitParam: limitParam,
    OrderByParam: orderByParam,
    PrevPageValuesParamPrefix: prevPageValuesParamPrefix,
    PrevPageValuesParse: parseFunc,
  }
}

func NewDefaultPaginatorParamsParser() *paginatorParamsParserImpl {
  return &paginatorParamsParserImpl{
    LimitParam: defaultLimitParam,
    OrderByParam: defaultOrderByParam,
    PrevPageValuesParamPrefix: defaultPrevPageValuesParamPrefix,
    PrevPageValuesParse: DefaultPrevPageValuesParse,
  }
}

func (p *PaginatorParamsParser) ParseQuery(query string) (params *PaginatorParams, err error) {
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

  // parsing order by list
  orderBy := []string{}
  if len(values[p.OrderByParam]) > 0 {
    orderBy = strings.Split(values[p.OrderByParam][0], ",")
  }

  // parsing prev page values
  prevPageValues, err := p.PrevPageValuesParse(paramsQuery, p.PrevPageValuesParamPrefix)
  if err != nil {
    return
  }

  return &PaginatorParams{
    Limit: values[p.LimitParam],
    OrderBy: orderBy,
    PrevPageValues: p.PrevPageValuesParse(paramsQuery),
  }, nil
}

func DefaultPrevPageValuesParse(queryParams map[string][]string, paramPrefix string) (prevPageValues map[string]string, err error) {
  prevPageValues = map[string]string{}
  if len(queryParams[paramPrefix]) > 0 {
    err := json.Unmarhsal(queryParams[paramPrefix][0], &prevPageValues)
    if err != nil {
      return nil, err
    }
  }

  return
}
