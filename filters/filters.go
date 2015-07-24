package filters

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/fiam/gounidecode/unidecode"
)

const (
	// Base64 filter
	FilterBase64 = "base64"

	// MD5 filter
	FilterMd5 = "md5"

	// Transliteration filter
	FilterTranslit = "translit"

	// Uppercase filter
	FilterUppercase = "uppercase"
)

type Filter func(input string) string

var filters = map[string]Filter{
	FilterBase64: func(input string) string {
		data := []byte(input)
		return base64.StdEncoding.EncodeToString(data)
	},
	FilterMd5: func(input string) string {
		hash := md5.Sum([]byte(input))
		return hex.EncodeToString(hash[:])
	},
	FilterTranslit: func(input string) string {
		return unidecode.Unidecode(input)
	},
	FilterUppercase: func(input string) string {
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
