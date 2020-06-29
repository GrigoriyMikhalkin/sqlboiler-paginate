package common

import (
	"testing"
)

func TestParseQuery(t *testing.T) {
	parser := NewDefaultPaginatorParamsParser()

	query := "/endpoint?limit=5"
	params, err := parser.ParseQuery(query)
	if err != nil {
		t.Errorf("Params parsing failed: %v", err)
	}
	if params.Limit != 5 {
		t.Errorf("Limit should be 5, instead go %d", params.Limit)
	}

	query = "/endpoint?limit=5&offset=5"
	params, err = parser.ParseQuery(query)
	if err != nil {
		t.Errorf("Params parsing failed: %v", err)
	}
	if params.Offset != 5 {
		t.Errorf("Offset should be 5, instead go %d", params.Offset)
	}

	query = `/users?limit=5&order_by=id,name-desc&prev_page_values={"id":5}`
	params, err = parser.ParseQuery(query)
	if err != nil {
		t.Errorf("Params parsing failed: %v", err)
	}
	if params.Limit != 5 {
		t.Errorf("Limit should be 5, instead go %d", params.Limit)
	}
	if len(params.OrderBy) != 2 {
		t.Error("expected ordering by id and name fields")
	}

	// Check ordering
	idOrdering := params.OrderBy[0]
	nameOrdering := params.OrderBy[1]
	if idOrdering.Field != "id" {
		t.Error("expected ordering by id first")
	}
	if idOrdering.Order != "asc" {
		t.Error("id param order should be 'asc'")
	}
	if nameOrdering.Field != "name" {
		t.Error("expected ordering by name second")
	}
	if nameOrdering.Order != "desc" {
		t.Error("name param order should be 'desc'")
	}

	// Check previous page values
	if len(params.PrevPage) != 1 {
		t.Errorf("expecting 1 prev page value, got %d", len(params.PrevPage))
	}
	if params.PrevPage["id"] != float64(5) {
		t.Errorf("Last id should be 5, got %d", params.PrevPage["id"])
	}
}
