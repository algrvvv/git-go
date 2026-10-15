package compare

import (
	"fmt"
)

type Score struct {
	Branch string

	BranchNameScore  float64
	TextScore        float64
	FreshnessScore   float64
	ActivityScore    float64
	BranchStateScore float64
}

func NewScore(name string) *Score {
	return &Score{
		Branch: name,
	}
}

func (s *Score) ToString() string {
	return fmt.Sprintf(
		"branch score (%s):\nBranchName: %0.2f\nTextScore: %0.2f\nFreshnessScore: %0.2f\nTotal: %0.2f\n",
		s.Branch,
		s.BranchNameScore,
		s.TextScore,
		s.FreshnessScore,
		s.Total(),
	)
}

func (s *Score) Total() float64 {
	score := 0.60*s.TextScore +
		0.20*s.BranchNameScore +
		0.10*s.FreshnessScore +
		0.10*s.ActivityScore

	return round(score)
}
