package request

import (
	"net/http"
	"net/url"
	"strconv"
)

// QueryHelper represent helper to get query string data
type QueryHelper struct {
	r  *http.Request
	uv url.Values
}

// NewQueryHelper is a function to create query helper struct
func NewQueryHelper(r *http.Request) *QueryHelper {
	return &QueryHelper{r, r.URL.Query()}
}

// GetString to get query string value with string data type, return defValue if query url not found
func (q *QueryHelper) GetString(p string, defValue string) string {
	sv := q.uv.Get(p)
	if sv != "" {
		return sv
	}
	return defValue
}

// GetInt to get query string value with integer data type, return defValue if query url not found or invalid
func (q *QueryHelper) GetInt(p string, defValue int) int {
	sv := q.uv.Get(p)
	if sv != "" {
		if v, err := strconv.Atoi(sv); err == nil {
			return v
		}
	}
	return defValue
}
