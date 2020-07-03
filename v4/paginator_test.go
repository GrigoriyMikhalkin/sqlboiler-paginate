package paginate

import (
	"reflect"
	"testing"

	"github.com/grigoriymikhalkin/sqlboiler-paginate/common"
)

func TestOffsetPagination(t *testing.T) {
	params := common.PaginatorParams{
		Limit:  5,
		Offset: 5,
	}
	qms := PaginationQueryMods(params)
	if len(qms) != 2 {
		t.Errorf("Expected QueryMods qnt is 2, instead got %d", len(qms))
	}

	qm1 := reflect.TypeOf(qms[0]).String()
	if qm1 != "qm.offsetQueryMod" {
		t.Errorf("First QueryMod expected type 'qm.offsetQueryMod', instead got %s", qm1)
	}

	qm2 := reflect.TypeOf(qms[1]).String()
	if qm2 != "qm.limitQueryMod" {
		t.Errorf("Second QueryMod expected type 'qm.limitQueryMod', instead got %s", qm2)
	}
}

func TestKeysetPagination(t *testing.T) {
	params := common.PaginatorParams{
		Limit:  5,
		Offset: 5,
		OrderBy: common.OrderByParams{
			{"id", "asc"},
			{"name", "desc"},
		},
		PrevPage: common.PrevPageValues{
			"id": 5,
		},
	}
	qms := PaginationQueryMods(params)
	if len(qms) != 3 {
		t.Errorf("Expected QueryMods qnt is 3, instead got %d", len(qms))
	}

	qm1 := reflect.TypeOf((qms[0])).String()
	if qm1 != "qmhelper.WhereQueryMod" {
		t.Errorf(
			"First QueryMod expected type 'qmhelper.WhereQueryMod', instead got %s", qm1)
	}

	qm2 := reflect.TypeOf((qms[1])).String()
	if qm2 != "qm.orderByQueryMod" {
		t.Errorf(
			"Second QueryMod expected type 'qm.orderByQueryMod', instead got %s", qm2)
	}

	qm3 := reflect.TypeOf((qms[2])).String()
	if qm3 != "qm.limitQueryMod" {
		t.Errorf(
			"Third QueryMod expected type 'qm.limitQueryMod', instead got %s", qm3)
	}
}
