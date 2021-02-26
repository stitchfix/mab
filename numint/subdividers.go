package numint

// EquallySpaced returns an SubDivider that divides intervals into a set number of equally-spaced sub-intervals.
func EquallySpaced(nSubIntervals int) SubDivider {
	return equallySpaced{
		nSubIntervals: nSubIntervals,
	}
}

type equallySpaced struct {
	nSubIntervals int
}

// SubDivide divides each interval into an equal number of sub-intervals and returns the slice of sub-intervals as a flat slice.
func (e equallySpaced) SubDivide(intervals []Interval) []Interval {
	result := make([]Interval, 0)

	for i := range intervals {
		subIntervals := e.subDivideSingle(intervals[i])
		for _, subInt := range subIntervals {
			result = append(result, subInt)
		}
	}

	return result
}

func (e equallySpaced) subDivideSingle(interval Interval) []Interval {
	h := (interval.B - interval.A) / float64(e.nSubIntervals)
	result := make([]Interval, e.nSubIntervals)

	x := interval.A
	for i := range result {
		result[i].A = x
		result[i].B = x + h
		x += h
	}
	return result
}
