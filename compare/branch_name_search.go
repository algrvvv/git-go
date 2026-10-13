package compare

import "github.com/algrvvv/git-go/domain"

func (s *Searcher) branchNameSearch(branch *domain.GitBranch) float64 {
	s.log.Debug("compare branch names", "query", s.query, "branch-name", branch.Name)

	if s.query == branch.Name {
		return 1.0
	}

	return 0
}
