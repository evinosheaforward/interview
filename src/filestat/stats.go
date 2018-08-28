package filestat

import (
		"iter"
		"math"
)

type stat struct {
		val *int
		count *int
}

type statData []stat

func (data statData) N() int {
		n := 0
		for datum := range data {
				n += float64(datum.val * datum.count)
		}
}

func Std(data statData) float64 {
	  mean := Mean(data)
		varTot, n := 0.0, 0.0
		for datum := range data {
				varTot += datum.count * math.Pow(mean - datum.val, 2.0)
				n += datum.count
		}
		return math.sqrt(tokenVarTot / (n - 1))
}

func Mean(data statData) float64 {
		total, n := 0.0, data.N()
		for datum := range data {
				total += datum.val * datum.count
				n += datum.count
		}
		return (total / n)
}

func Median(data statData) float64 {
		vals = make([]*int, statInfo.N())
		n := 0
		for datum := range data {
				n += datum.count
				for i := 0; i < datum.count; i++ {
						vals = append(vals, datum.val)
				}
		}
		sort.Ints(vals)
		if !Even(n) {
				return float64(*vals[math.Floor(float64(n)/2)])
		}
		return (*vals[math.Floor(float64(n)/2)] + *vals[math.Ceil(float64(n)/2)]) / 2
}

func Even(number int) bool {
    return number%2 == 0
}
