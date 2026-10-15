package compare

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/algrvvv/git-go/domain"
)

func (s *Searcher) textSearch(branch *domain.GitBranch) float64 {
	tokenSimilarity := s.compareTokens(branch.Meta)
	fuzzySimilarity := s.compareFuzzy(branch.Meta)
	prefixSimilarity := s.comparePrefix(branch.Meta)

	total := 0.45*tokenSimilarity +
		0.35*fuzzySimilarity +
		0.15*prefixSimilarity
	total = round(total)

	s.log.Debug(
		"text-search: compare branch result",
		slog.String("branch", branch.Name),
		slog.Float64("tokenSimilarity", tokenSimilarity),
		slog.Float64("fuzzySimilarity", fuzzySimilarity),
		slog.Float64("prefixSimilarity", prefixSimilarity),
		slog.Float64("total", total),
	)

	return total
}

func (s *Searcher) compareTokens(meta string) float64 {
	if meta == s.query {
		return 1.0
	}

	queryTokens := strings.Split(s.query, " ")
	metaTokens := strings.Split(meta, " ")

	var coincidence int
	for _, token := range queryTokens {
		if slices.Contains(metaTokens, token) {
			coincidence++
		}
	}

	return float64(coincidence / len(queryTokens))
}

func (s *Searcher) comparePrefix(meta string) float64 {
	queryTokens := strings.Split(s.query, " ")
	metaTokens := strings.Split(meta, " ")

	var coincidence int
	for i := range metaTokens {
		if i >= len(queryTokens) {
			break
		}

		if strings.Contains(metaTokens[i], queryTokens[i]) {
			coincidence++
			continue
		}

		break
	}

	return float64(coincidence / len(metaTokens))
}
