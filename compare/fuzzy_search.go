package compare

import (
	"log/slog"
	"unicode/utf8"
)

const minLengthThreshold = 32

func (s *Searcher) compareFuzzy(meta string) float64 {
	diff := float64(max(len(meta), len(s.query)))
	if diff == 0 {
		return 0
	}

	distance := float64(ComputeDistance(meta, s.query))
	result := 1.0 - (distance / diff)
	result = round(result)

	s.log.Debug(""+
		"fuzzy-search",
		slog.String("meta", meta),
		slog.String("query", s.query),
		slog.Float64("distance", distance),
		slog.Int("len-meta", len(meta)),
		slog.Int("len-query", len(s.query)),
		slog.Float64("result", result),
	)

	return result
}

func ComputeDistance(a, b string) int {
	if len(a) == 0 {
		return utf8.RuneCountInString(b)
	}

	if len(b) == 0 {
		return utf8.RuneCountInString(a)
	}

	if a == b {
		return 0
	}

	s1 := []rune(a)
	s2 := []rune(b)

	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}

	s1, s2 = trimLongestCommonSuffix(s1, s2)
	s1, s2 = trimLongestCommonPrefix(s1, s2)

	lenS1 := len(s1)
	lenS2 := len(s2)

	var x []uint16
	if lenS1+1 > minLengthThreshold {
		x = make([]uint16, lenS1+1)
	} else {
		x = make([]uint16, minLengthThreshold)
		x = x[:lenS1+1]
	}

	for i := 1; i < len(x); i++ {
		x[i] = uint16(i)
	}

	_ = x[lenS1]
	y := x[1:]
	y = y[:lenS1]

	for i := range lenS2 {
		prev := uint16(i + 1)
		for j := range lenS1 {
			current := x[j]
			if s2[i] != s1[j] {
				current = min(x[j], prev, y[j]) + 1
			}
			x[j] = prev
			prev = current
		}
		x[lenS1] = prev
	}
	return int(x[lenS1])
}

func trimLongestCommonSuffix(a, b []rune) ([]rune, []rune) {
	m := min(len(a), len(b))
	a2 := a[len(a)-m:]
	b2 := b[len(b)-m:]
	i := len(a2)
	b2 = b2[:i]
	for ; i > 0 && a2[i-1] == b2[i-1]; i-- {
	}
	return a[:len(a)-len(a2)+i], b[:len(b)-len(b2)+i]
}

func trimLongestCommonPrefix(a, b []rune) ([]rune, []rune) {
	var i int
	for m := min(len(a), len(b)); i < m && a[i] == b[i]; i++ {
	}
	return a[i:], b[i:]
}
