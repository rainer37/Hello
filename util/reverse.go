package util

import (
	"strings"
)
// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
//The full pathname of the directory that contains Daikon

func Get_PreSlash(s string) (string,string) {

	if strings.Index(s, "/") == -1 {
		return "",""
	}

	prefix := string(s[0:strings.LastIndex(s,"/")+1])
	after := string(s[strings.LastIndex(s,"/")+1:])
	return prefix, after
}