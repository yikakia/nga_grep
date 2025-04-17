package data

import (
	"github.com/samber/lo"
)

func GetMA_N(dots []Dot, n int) []float64 {
	var maNValues []float64
	values := lo.Map(dots, func(item Dot, index int) int {
		return item.Count
	})
	window := lo.Sum(values[:n-1])
	for i, v := range values[n-1:] {
		// i 是第n个
		window += v
		if i > 0 {
			window -= values[i]
		}

		maNValues = append(maNValues, float64(window)/float64(n))
	}
	return maNValues
}
