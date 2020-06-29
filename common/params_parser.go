package common

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultLimitParam          = "limit"
	defaultOffsetParam         = "offset"
	defaultOrderByParam        = "order_by"
	defaultPrevPageValuesParam = "prev_page_values"
)

type QueryParams = map[string][]string

type OrderByParseFunc = func(params QueryParams, orderByParam string) (orderByParams OrderByParams, err error)

type PrevPageValuesParseFunc = func(params QueryParams, prevPageValueParam string) (prevPageValues PrevPageValues, err error)

type PaginatorParamsParser interface {
	ParseQuery(query string) (params PaginatorParams, err error)
}

type defaultPaginatorParamsParser struct {
	LimitParam          string
	OffsetParam         string
	OrderByParam        string
	PrevPageValuesParam string

	OrderByParse        OrderByParseFunc
	PrevPageValuesParse PrevPageValuesParseFunc
}

func NewPaginatorParamsParser(limitParam, offsetParam, orderByParam, prevPageValuesParam string,
	orderByParse OrderByParseFunc, prevPageValuesParse PrevPageValuesParseFunc) *defaultPaginatorParamsParser {
	orderByParseFunc := DefaultOrderByParse
	if orderByParse != nil {
		orderByParseFunc = orderByParse
	}
	prevPageValuesParseFunc := DefaultPrevPageValuesParse
	if prevPageValuesParse != nil {
		prevPageValuesParseFunc = prevPageValuesParse
	}

	return &defaultPaginatorParamsParser{
		LimitParam:          limitParam,
		OffsetParam:         offsetParam,
		OrderByParam:        orderByParam,
		PrevPageValuesParam: prevPageValuesParam,
		OrderByParse:        orderByParseFunc,
		PrevPageValuesParse: prevPageValuesParseFunc,
	}
}

func NewDefaultPaginatorParamsParser() *defaultPaginatorParamsParser {
	return &defaultPaginatorParamsParser{
		LimitParam:          defaultLimitParam,
		OffsetParam:         defaultOffsetParam,
		OrderByParam:        defaultOrderByParam,
		PrevPageValuesParam: defaultPrevPageValuesParam,
		OrderByParse:        DefaultOrderByParse,
		PrevPageValuesParse: DefaultPrevPageValuesParse,
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
	orderBy, err := p.OrderByParse(values, p.OrderByParam)
	if err != nil {
		return
	}

	prevPage, err := p.PrevPageValuesParse(values, p.PrevPageValuesParam)

	return PaginatorParams{
		Limit:    limit,
		Offset:   offset,
		OrderBy:  orderBy,
		PrevPage: prevPage,
	}, nil
}

func DefaultOrderByParse(params QueryParams, orderByParam string) (orderByParams OrderByParams, err error) {
	orderBy := []string{}
	if len(params[orderByParam]) > 0 {
		orderBy = strings.Split(params[orderByParam][0], ",")
	} else {
		return
	}

	for _, p := range orderBy {
		x := strings.Split(p, "-")
		field := x[0]
		order := "asc"
		if len(x) > 1 {
			order = x[1]
		}

		orderByParams = append(orderByParams, &OrderByParam{
			Field: field,
			Order: order,
		})
	}

	return
}

func DefaultPrevPageValuesParse(params QueryParams, prevPageValueParam string) (prevPageValues PrevPageValues, err error) {
	prevPageValues = PrevPageValues{}
	if len(params[prevPageValueParam]) > 0 {
		err = json.Unmarshal([]byte(params[prevPageValueParam][0]), &prevPageValues)
		if err != nil {
			return
		}
	}

	return
}
