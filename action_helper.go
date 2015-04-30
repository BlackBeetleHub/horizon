package horizon

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
)

const (
	DefaultPageSize = 10
)

// ActionHelper wraps the goji context and provides helper functions
// to make defining actions easier.
//
// Notably, this object provides a means of more simply extracting information
// from the Env and URLParams.  Any call to the Get* methods (GetInt, GetString, etc.)
// that fails will populate the Err field and subsequent calls to Get* will be no
// ops.  This allows the simpler pattern:
//
//	ah = &ActionHelper{C:c}
//	id := ah.GetInt("id")
//	order := ah.GetString("order")
//
//	if ah.Err() != nil {
//	  // write an error response here
//	  return
//	}
//
type ActionHelper struct {
	c   web.C
	r   *http.Request
	err error
}

func (a *ActionHelper) Err() error {
	return a.err
}

func (a *ActionHelper) App() *App {
	return a.c.Env["app"].(*App)
}

// Gets a string from either the URLParams or query string.
// This method prioritizes the URLParams over the Query.
//
// TODO: Add form support, prioritized over query
func (a *ActionHelper) GetString(name string) string {
	if a.err != nil {
		return ""
	}

	fromUrl, ok := a.c.URLParams[name]

	if ok {
		return fromUrl
	}

	return a.r.URL.Query().Get(name)
}

func (a *ActionHelper) GetInt64(name string) int64 {
	if a.err != nil {
		return 0
	}

	asStr := a.GetString(name)

	if asStr == "" {
		return 0
	}

	asI64, err := strconv.ParseInt(asStr, 10, 64)

	if err != nil {
		a.err = err
		return 0
	}

	return asI64
}

func (a *ActionHelper) GetInt32(name string) int32 {
	if a.err != nil {
		return 0
	}

	asStr := a.GetString(name)

	if asStr == "" {
		return 0
	}

	asI64, err := strconv.ParseInt(asStr, 10, 32)

	if err != nil {
		a.err = err
		return 0
	}

	return int32(asI64)
}

func (a *ActionHelper) GetPagingParams() (after string, order string, limit int32) {
	if a.err != nil {
		return
	}

	after = a.GetString("after")
	order = a.GetString("order")
	limit = a.GetInt32("limit")

	if limit == 0 {
		limit = DefaultPageSize
	}

	if order == "" {
		order = "asc"
	}

	return
}