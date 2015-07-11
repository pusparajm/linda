package armenia

import (
	"bytes"
	"strings"
	"unicode"

	"golang.org/x/exp/utf8string"
)

var EasternLanguage = map[rune]string{
	'ա': "a",
	'բ': "b",
	'գ': "g",
	'դ': "d",
	'ե': "e",
	'զ': "z",
	'է': "e",
	'ը': "y",
	'թ': "t",
	'ժ': "jh",
	'ի': "i",
	'լ': "l",
	'խ': "x",
	'ծ': "c",
	'կ': "k",
	'հ': "h",
	'ձ': "d",
	'ղ': "gh",
	'ճ': "tw",
	'մ': "m",
	'յ': "y",
	'ն': "n",
	'շ': "sh",
	'ո': "o",
	'չ': "ch",
	'պ': "p",
	'ջ': "j",
	'ռ': "r",
	'ս': "s",
	'վ': "v",
	'տ': "t",
	'ր': "r",
	'ց': "c",
	'ւ': "w",
	'փ': "p",
	'ք': "q",
	'օ': "o",
	'ֆ': "f",
	'և': "ev",
}

var WesternLanguage = map[rune]string{
	'ա': "a",
	'բ': "p",
	'գ': "k",
	'դ': "t",
	'ե': "e",
	'զ': "z",
	'է': "e",
	'ը': "y",
	'թ': "t",
	'ժ': "zh",
	'ի': "i",
	'լ': "l",
	'խ': "dz",
	'ծ': "g",
	'կ': "g",
	'հ': "h",
	'ձ': "tz",
	'ղ': "gh",
	'ճ': "j",
	'մ': "m",
	'յ': "h",
	'ն': "n",
	'շ': "sh",
	'ո': "o",
	'չ': "ch",
	'պ': "b",
	'ջ': "ch",
	'ռ': "r",
	'ս': "s",
	'վ': "v",
	'տ': "d",
	'ր': "r",
	'ց': "c",
	'ւ': "w",
	'փ': "p",
	'ք': "q",
	'օ': "o",
	'ֆ': "f",
	'և': "ev",
}

// Source: http://github.com/mdigger/translit
func Translit(text string, translitMap map[rune]string) string {
	var result bytes.Buffer
	utf8text := utf8string.NewString(text)
	length := utf8text.RuneCount()
	for index := 0; index < length; index++ {
		runeValue := utf8text.At(index)
		switch str, ok := translitMap[unicode.ToLower(runeValue)]; {
		case !ok:
			result.WriteRune(runeValue)
		case str == "":
			continue
		case unicode.IsUpper(runeValue):
			// Если следующий или предыдущий символ тоже заглавная буква, то все буквы строки
			// заглавные. Иначе, заглавная только первая буква.
			if (length > index+1 && unicode.IsUpper(utf8text.At(index+1))) ||
				(index > 0 && unicode.IsUpper(utf8text.At(index-1))) {
				str = strings.ToUpper(str)
			} else {
				str = strings.Title(str)
			}
			fallthrough
		default:
			result.WriteString(str)
		}
	}
	return result.String()
}
