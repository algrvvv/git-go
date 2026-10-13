package compare

import (
	"log/slog"
	"sync"

	"github.com/algrvvv/git-go/domain"
)

type Searcher struct {
	branches []*domain.GitBranch
	query    string

	wg         *sync.WaitGroup
	candidates chan *candidate

	log *slog.Logger
}

type candidate struct {
	branch *domain.GitBranch
	score  *Score
	m      sync.Mutex
}

func newCandidate() *candidate {
	return &candidate{
		branch: nil,
		score:  &Score{},
		m:      sync.Mutex{},
	}
}

func (c *candidate) setWithoutLocks(b *domain.GitBranch, s *Score) {
	c.branch = b
	c.score = s
}

func (c *candidate) set(b *domain.GitBranch, s *Score) {
	c.m.Lock()
	defer c.m.Unlock()
	c.branch = b
	c.score = s
}

func NewSearcher(log *slog.Logger, query string, branches []*domain.GitBranch) *Searcher {
	return &Searcher{
		query:      query,
		branches:   branches,
		log:        log,
		wg:         &sync.WaitGroup{},
		candidates: make(chan *candidate),
	}
}

func (s *Searcher) SearchByQuery() (branch *domain.GitBranch, err error) {
	for _, br := range s.branches {
		s.wg.Add(1)
		go s.rank(br)
	}

	go func() {
		s.wg.Wait()
		close(s.candidates)
	}()

	result := newCandidate()
	for cnd := range s.candidates {
		s.log.Debug(cnd.score.ToString())
		s.log.Debug(
			"total scores for branches in current moment",
			"result-total",
			result.score.Total(),
			"candidate-total",
			cnd.score.Total(),
		)

		if result.score.Total() < cnd.score.Total() {
			result.set(cnd.branch, cnd.score)
		}
	}

	return result.branch, nil
}

func (s *Searcher) rank(branch *domain.GitBranch) {
	defer s.wg.Done()

	score := NewScore(branch.Name)
	score.BranchNameScore = s.branchNameSearch(branch)
	score.TextScore = s.textSearch(branch)

	c := newCandidate()
	c.setWithoutLocks(branch, score)
	s.candidates <- c
}
