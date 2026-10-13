package compare

import "github.com/algrvvv/git-go/domain"

func (s *Searcher) textSearch(branch *domain.GitBranch) float64 {
	s.log.Debug("compare branch meta", "query", s.query, "branch-meta", branch.Meta)

	// если точное совпадение
	if s.query == branch.Meta {
		return 1.0
	}

	// левенштейн и вся хуйня
	return 0
}
