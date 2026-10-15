package compare

import "math"

func round(x float64) float64 {
	return math.Round(x*1000) / 1000
}
