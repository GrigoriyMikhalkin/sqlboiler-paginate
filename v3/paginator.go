package paginate

import (
	"fmt"

	"github.com/grigoriymikhalkin/sqlboiler-paginate/common"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func PaginationQueryMods(params common.PaginatorParams) []qm.QueryMod {
	if params.Offset > 0 && len(params.OrderBy) < 1 {
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
	for _, param := range params.OrderBy {
		// update order by
		if orderByQuery != "" {
			orderByQuery += ", "
		}
		orderByQuery += param.Field + " " + param.Order

		// update where filter
		var sign string
		var mod qm.QueryMod
		if v, ok := params.PrevPage[param.Field]; ok {
			switch param.Order {
			case "asc":
				sign = ">"
			case "desc":
				sign = "<"
			}

			value := fmt.Sprintf("%v", v)
			if len(mods) > 0 {
				mod = qm.And(param.Field+sign+"?", value)
			} else {
				mod = qm.Where(param.Field+sign+"?", value)
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
