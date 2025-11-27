package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// TrimSpace trims the spaces of string
func TrimSpace(name string) string {
	//name not for accepting full space
	name = strings.TrimSpace(name)
	//name not accepting double spaces between two words &
	//for not allow to save same two words with single space
	reg := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	name = reg.ReplaceAllString(name, " ")

	return name
}

// ToCamelCase to convert first letter upper case
func ToCamelCase(str string) string {
	return strings.Title(strings.ToLower(str))
}

// StringAppendUseCommaSep append string and based on the string length append comma seprator or not
func StringAppendUseCommaSep(str string, appendVal string) string {
	if len(str) > 0 {
		str = str + ", "
	}
	return str + appendVal
}

// StringToInt64 convert string to int64
func StringToInt64(str string) (int64, bool) {
	res, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return res, true
	}
	return res, false
}

func FormatInt32ToString(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}
