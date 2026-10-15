package compare

import (
	"log/slog"

	"github.com/algrvvv/git-go/domain"
)

func (s *Searcher) freshnessSearch(branch *domain.GitBranch) float64 {
	var offset float64
	for i, b := range s.branches {
		if b.Name == branch.Name {
			offset = float64(i)
			break
		}
	}

	maxValue := float64(len(s.branches))
	result := 1.0 - offset/maxValue
	result = round(result)

	s.log.Debug(
		"Freshness search result: %f",
		slog.String("branch", branch.Name),
		slog.Float64("offset", offset),
		slog.Float64("maxValue", maxValue),
		slog.Float64("result", result),
	)

	return result
}
