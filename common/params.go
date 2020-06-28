package common

type OrderByParam struct {
  Order string
  LastValue string
}

type OrderByParams = map[string]*OrderByParam

type PaginatorParams struct {
  Limit   int
  Offset  int
  OrderBy OrderByParams
}
