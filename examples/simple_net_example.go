package main

import (
	"fmt"
	"net/http"

	"github.com/grigoriymikhalkin/sqlboiler-paginate/common"
	"github.com/grigoriymikhalkin/sqlboiler-paginate/v4"
)

func main() {
	http.HandleFunc("/users", users)
	http.ListenAndServe(":8000", nil)
}

/*
 * Requests for testing:
 *   - Limit & Offset:
 *     /users?limit=5&offset=5
 *   - Order by id ASC, name DESC and filter records where id > 5:
 *     /users?limit=5&order_by=id-asc,name-desc&prev_page_values={%22id%22:5}
 */
func users(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Query: %s\n", req.URL.RawQuery)

	// Parse params
	parser := common.NewDefaultPaginatorParamsParser()
	params, err := parser.ParseQuery(req.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Params: %+v\n", params)

	// Generate pagination query mods
	mods := paginate.PaginationQueryMods(params)
	fmt.Printf("Mods: %+v", mods)

	// Query DB
	// ...
}