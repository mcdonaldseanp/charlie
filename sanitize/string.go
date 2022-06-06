package sanitize

import "strings"

func ReplaceAllNewlines(origin string) string {
	var replacer = strings.NewReplacer(
		"\r\n", "",
		"\r", "",
		"\n", "",
	)
	return replacer.Replace(origin)
}

func ReplaceAllSpaces(origin string) string {
	var replacer = strings.NewReplacer(
		"\r\n", "",
		"\r", "",
		"\n", "",
		" ", "",
		"\t", "",
	)
	return replacer.Replace(origin)
}
