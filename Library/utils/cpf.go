package utils

import "strings"

func TrimCpfCnpj(doc string) string {
	s := strings.Replace(doc, ".", "", -1)
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, "-", "", -1)
	s = strings.Replace(s, "*", "", -1)
	return s
}
