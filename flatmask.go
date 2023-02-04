package flatmask

import (
	"sort"
	"strings"
)

// MaskProvider - representation of the mask instance with paths
// what can be reset, normalize, stringed, etc.
// fieldmaskpb.FieldMask implements MaskProvider interface and can be
// passed into Reduce.
type MaskProvider interface {
	// GetPaths - returns the string paths core of mask
	GetPaths() []string

	// This logic is encapsulated in the mask implementation.
	//
	// Normalize - converts the mask to its canonical form where all paths are sorted
	// and redundant paths are removed.
	Normalize()

	// This logic is encapsulated in the mask implementation.
	//
	// Reset - resets the mask to the default state.
	Reset()

	// This logic is encapsulated in the mask implementation.
	//
	// String - returns ыtring representation of the mask, with all its paths
	String() string
}

type ReducedMask interface {
	// GetPaths returns the string paths of reduced mask.
	GetPaths() []string
}

type reducedPaths struct {
	p []string
}

func (rp *reducedPaths) GetPaths() []string {
	return rp.p
}

// ReduceDegree - the degree of reduction of the mask.
// For example, there is a path {a_1, a_2,...,a_n},
// where n is the number of nodes (or fields in the field mask).
// If the degree of reduction is 3
// then we will leave only the first three nodes - {a_1, a_2, a_3}.
type ReduceDegree int

const (
	// The degree to which only the root node remains
	Total ReduceDegree = 1 + iota
	// The degree to which only the root node with his child remains
	RootChild
	// Other degrees can be passed as simple int number
)

// Reduce - reduces the mask to the specified degree and removes matches
func Reduce(mask MaskProvider, degree ReduceDegree) ReducedMask {
	mask.Normalize()

	rd := reduceToDegree(mask, int(degree))
	rd.p = removeDuplicateStrings(rd.p)
	return rd
}

func reduceToDegree(mask MaskProvider, degree int) *reducedPaths {
	maskPaths := mask.GetPaths()
	for i, p := range maskPaths {
		if strings.IndexByte(p, '.') == -1 || len(p) < degree+1 {
			continue
		}

		ri := getMaxReducedIndex(p, degree)
		if ri == -1 {
			continue
		}

		maskPaths[i] = p[:ri]
	}

	return &reducedPaths{
		p: maskPaths,
	}
}

func getMaxReducedIndex(path string, degree int) int {
	ri, rCount := -1, 0
	for i := 0; i < len(path) && rCount < degree; i++ {
		if path[i] == '.' {
			ri = i
			rCount++
		}
	}

	if rCount < degree {
		return -1
	}

	return ri
}

func removeDuplicateStrings(s []string) []string {
	if len(s) < 1 {
		return s
	}

	sort.Strings(s)
	p := 1
	for i := 1; i < len(s); i++ {
		if s[i-1] != s[i] {
			s[p] = s[i]
			p++
		}
	}

	return s[:p]
}
