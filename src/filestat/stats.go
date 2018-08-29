package filestat

import (
	"math"
	"sort"
)

type statDatum struct {
	val   int
	count int
}

type statData []statDatum

func (data statData) N() int {
	n := 0
	for _, datum := range data {
		n += datum.count
	}
	return n
}

func Std(data statData) float64 {
	mean := Mean(data)
	varTot, n := 0.0, 0.0
	for _, datum := range data {
		varTot += float64(datum.count) * math.Pow(mean-float64(datum.val), 2.0)
		n += float64(datum.count)
	}
	return math.Sqrt(varTot / (n - 1))
}

func Mean(data statData) float64 {
	total, n := 0.0, data.N()
	for _, datum := range data {
		total += float64(datum.val * datum.count)
	}
	return (total / float64(n))
}

func Median(data statData) float64 {
	n := data.N()
	vals := make([]int, n)
	idx := 0
	for _, datum := range data {
		for i := 0; i < datum.count; i++ {
			vals[idx] = datum.val
			idx++
		}
	}
	sort.Ints(vals)
	if Even(n) {
		return float64(vals[LowerIdx(n)]+vals[UpperIdx(n)]) / 2.0
	}
	return float64(vals[LowerIdx(n)])
}

func UpperIdx(n int) int {
	return int(math.Ceil(float64(n) / 2.0))
}

func LowerIdx(n int) int {
	return int(math.Floor(float64(n) / 2.0))
}

func Even(number int) bool {
	return number%2 == 0
}
