package paginate

import (
  "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Paginator struct {}

func (p *Paginator) PaginationQueryMods(params PaginatorParams) []qm.QueryMod {
  mods := []qm.QueryMod{}

  return mods
}
