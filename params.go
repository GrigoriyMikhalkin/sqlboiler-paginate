package paginate

type OrderByParam struct {
  Order string
  LastValue interface{}
}

type OrderByParams = map[string]*OrderByParam

type PaginatorParams struct {
  Limit   int
  Offset  int
  OrderBy OrderByParams
}
