package compare

import (
	"fmt"
)

type Score struct {
	Branch string

	BranchNameScore float64
	TextScore       float64
	FuzzyScore      float64
}

func NewScore(name string) *Score {
	return &Score{
		Branch: name,
	}
}

func (s *Score) ToString() string {
	return fmt.Sprintf(
		"branch score (%s):\nBranchName: %0.2f\nTextScore: %0.2f\nFuzzyScore: %0.2f\nTotal: %0.2f\n",
		s.Branch,
		s.BranchNameScore,
		s.TextScore,
		s.FuzzyScore,
		s.Total(),
	)
}

func (s *Score) Total() float64 {
	return s.BranchNameScore + +s.TextScore + s.FuzzyScore
}
