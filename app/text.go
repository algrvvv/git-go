package app

import "strings"

// trimLastWord удаляет последнее слово (ctrl+w в поле фильтра).
func trimLastWord(s string) string {
	s = strings.TrimRight(s, " ")
	if s == "" {
		return s
	}
	if i := strings.LastIndex(s, " "); i >= 0 {
		return s[:i]
	}
	return ""
}

// truncateRunes обрезает строку до n рун (не байт).
func truncateRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}
