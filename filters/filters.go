package filters

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/fiam/gounidecode/unidecode"
)

type Filter func(input string) string

var filters = map[string]Filter{
	"base64": func(input string) string {
		data := []byte(input)
		return base64.StdEncoding.EncodeToString(data)
	},
	"md5": func(input string) string {
		hash := md5.Sum([]byte(input))
		return hex.EncodeToString(hash[:])
	},
	"translit": func(input string) string {
		return unidecode.Unidecode(input)
	},
	"uppercase": func(input string) string {
		return strings.ToUpper(input)
	},
}

// Creates new Filter instance
func New(id string) Filter {
	if filter, ok := filters[id]; ok {
		return filter
	} else {
		return nil
	}
}
