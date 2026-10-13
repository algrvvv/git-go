package compare

import (
	"sort"
	"strings"

	"github.com/algrvvv/git-go/domain"
)

// Item — ветка с весом совпадения с поисковым запросом.
// Чем больше Score, тем релевантнее ветка.
type Item struct {
	Branch *domain.GitBranch
	Score  float64
}

// FilterAndRank — фильтр + ранжирование (AND по токенам, имя важнее описания).
func FilterAndRank(branches []*domain.GitBranch, filter string) []Item {
	f := strings.ToLower(strings.TrimSpace(filter))
	if f == "" {
		out := make([]Item, 0, len(branches))
		for _, b := range branches {
			out = append(out, Item{Branch: b})
		}
		return out
	}

	tokens := strings.Fields(f)
	var out []Item
	for _, b := range branches {
		name := strings.ToLower(b.Name)
		meta := strings.ToLower(b.Meta)
		combined := name + " " + meta

		matched := true
		score := 0.0
		for _, tok := range tokens {
			if !strings.Contains(combined, tok) {
				matched = false
				break
			}
			switch {
			case name == tok:
				score += 100
			case strings.HasPrefix(name, tok):
				score += 30
			case strings.Contains(name, tok):
				score += 20
			}
			switch {
			case meta == tok:
				score += 15
			case strings.Contains(meta, " "+tok) || strings.HasPrefix(meta, tok):
				score += 8
			case strings.Contains(meta, tok):
				score += 4
			}
			// короткие имена веток (номера задач) — точное вхождение ценится
			if strings.Contains(name, tok) && len(tok) >= 3 {
				score += 2
			}
		}
		if matched {
			// ветки с описанием чуть выше пустых при равном скоре
			if b.Meta != "" {
				score += 0.5
			}
			out = append(out, Item{Branch: b, Score: score})
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Score > out[j].Score
	})
	return out
}
