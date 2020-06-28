package common

import (
	"testing"
)

func TestParseQuery(t * testing.T) {
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

  query = `/users?limit=5&order_by=id&prev_page_values={"id":5}`
  params, err = parser.ParseQuery(query)
  if err != nil {
    t.Errorf("Params parsing failed: %v", err)
  }
  if params.Limit != 5 {
    t.Errorf("Limit should be 5, instead go %d", params.Limit)
  }
  orderByParam, ok := params.OrderBy["id"]
  if !ok {
    t.Error("id field should be in params")
  }
  if orderByParam.Order != "asc" {
    t.Error("id param order should be 'asc'")
  }
  if orderByParam.LastValue != "5" {
    t.Errorf("Last id should be 5, got %s", orderByParam.LastValue)
  }
}
