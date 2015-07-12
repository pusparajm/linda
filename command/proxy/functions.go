package proxy

import (
	"net/url"
)

// Predefined html/template helpers
var FuncMap = map[string]interface{}{
	"query_unescape": func(a string) string { i, _ := url.QueryUnescape(a); return i },
}
