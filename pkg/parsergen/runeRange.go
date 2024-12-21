package parsergen

import (
	"slices"
	"math"
)

var RUNE_MAX = rune(math.Pow(2, 31) - 1)

// A rune range represents a range of runes (inclusive, exclusive)
type RuneRange struct {
	low  rune
	high rune
}

// minimizes the ranges
func minimizeRanges(ranges []RuneRange) []RuneRange {
	if len(ranges) <= 1 {
		return ranges
	}

	slices.SortFunc(ranges, func(a RuneRange, b RuneRange) int {
		return int(a.low) - int(b.low)
	})

	outputRanges := make([]RuneRange, 0)
	currRange := ranges[0]

	for _, r := range ranges {
		if r.low > currRange.high {
			outputRanges = append(outputRanges, currRange)
			currRange = r
		} else if r.high > currRange.high {
			currRange.high = r.high
		}
	}

	outputRanges = append(outputRanges, currRange)

	return outputRanges
}

func makeRangesThatIgnore(min rune, max rune, ignore ...rune) []RuneRange {
	slices.SortFunc(ignore, func(a,b rune) int {
		return int(a-b)
	})

	outputRanges := make([]RuneRange, 0)
	currLow := min

	for _, i := range ignore {
		if i < min || i > max {
			continue
		}

		if i == min {
			// weird edge case
			currLow = min + 1
			continue
		}

		outputRanges = append(outputRanges, RuneRange{
			low: currLow,
			high: i,
		})

		currLow = i + 1
	}

	if currLow - 1 < max {
		outputRanges = append(outputRanges, RuneRange{
			low: currLow,
			high: max,
		})
	}

	return outputRanges
}