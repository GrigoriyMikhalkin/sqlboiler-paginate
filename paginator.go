package paginate

import (
  "strings"

  "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Paginator struct {}

func (p *Paginator) PaginationQueryMods(params *PaginatorParams) []qm.QueryMod {
  var mods []qm.QueryMod

  orderByQuery := ""
  for param, val := range params.OrderBy {
    if val != "" {
      parts := strings.Split(param, " ")
      sign := ">"
      if len(parts) > 1 {
        if parts[1] == "desc" {
          sign = "<"
        }
      }

      if len(mods) > 0 {
        mods = append(mods, qm.And(parts[0] + sign + "?", val))
      } else {
        mods = append(mods, qm.Where(parts[0] + sign + "?", val))
      }
    }

    if orderByQuery != "" {
      orderByQuery +=", "
    }
    orderByQuery += param
  }

  if orderByQuery != "" {
    mods = append(mods, qm.OrderBy(orderByQuery))
  }

  if params.Limit > -1 {
    mods = append(mods, qm.Limit(params.Limit))
  }

  return mods
}
