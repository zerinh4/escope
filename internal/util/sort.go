package util

import (
	"sort"
	"strings"
)

type SortDirection int

const (
	Ascending SortDirection = iota
	Descending
)

func (d SortDirection) String() string {
	switch d {
	case Ascending:
		return "asc"
	case Descending:
		return "desc"
	default:
		return "asc"
	}
}

func SortStrings(strSlice []string, direction SortDirection) {
	sort.Slice(strSlice, func(i, j int) bool {
		result := strings.Compare(strSlice[i], strSlice[j])
		if direction == Descending {
			return result > 0
		}
		return result < 0
	})
}
