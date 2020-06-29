package common

type OrderByParam struct {
	Field string
	Order string
}

type OrderByParams = []*OrderByParam

type PrevPageValues = map[string]interface{}

type PaginatorParams struct {
	Limit    int
	Offset   int
	OrderBy  OrderByParams
	PrevPage PrevPageValues
}
