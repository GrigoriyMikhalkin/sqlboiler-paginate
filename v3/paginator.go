package paginate

import (
  "github.com/grigoriymikhalkin/sqlboiler-paginate/common"
  "github.com/volatiletech/sqlboiler/queries/qm"
)

func PaginationQueryMods(params common.PaginatorParams) []qm.QueryMod {
  if params.Offset > 0 && params.OrderBy == nil {
    return offsetPagination(params)
  }

  return keysetPagination(params)
}

func offsetPagination(params common.PaginatorParams) []qm.QueryMod {
  return []qm.QueryMod{
    qm.Offset(params.Offset),
    qm.Limit(params.Limit),
  }
}

func keysetPagination(params common.PaginatorParams) []qm.QueryMod {
  var mods []qm.QueryMod

  orderByQuery := ""
  for field, param := range params.OrderBy {
    // update order_by
    if orderByQuery != "" {
      orderByQuery += ", "
    }
    orderByQuery += field + " " + param.Order

    // update where filter
    var sign string
    var mod qm.QueryMod
    if param.LastValue != "" {
      switch param.Order {
      case "asc":
        sign = ">"
      case "desc":
        sign = "<"
      }

      if len(mods) > 0 {
        mod = qm.And(field + sign + "?", param.LastValue)
      } else {
        mod = qm.Where(field + sign + "?", param.LastValue)
      }

      mods = append(mods, mod)
    }
  }

  if orderByQuery != "" {
    mods = append(mods, qm.OrderBy(orderByQuery))
  }

  if params.Limit > -1 {
    mods = append(mods, qm.Limit(params.Limit))
  }

  return mods
}
