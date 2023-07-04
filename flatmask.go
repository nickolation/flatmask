package flatmask

import (
	"sort"
	"strings"
)

// PathsProvider - representation of the mask instance with paths
// what can be reset, normalize, stringed, etc.
// fieldmaskpb.FieldMask implements PathsProvider interface and can be
// passed into Reduce.
type PathsProvider interface {
	// GetPaths - returns the string paths core of mask
	GetPaths() []string
}

// ReduceDegree - the degree of reduction of the mask.
// For example, there is a path {a_1, a_2,...,a_n},
// where n is the number of nodes (or fields in the field mask).
// If the degree of reduction is 3
// then we will leave only the first three nodes - {a_1, a_2, a_3}.
type ReduceDegree uint16

const (
	// The degree to which only the root node remains
	Total ReduceDegree = 1 + iota
	// The degree to which only the root node with his child remains
	RootChild
	// Other degrees can be passed as simple int number
)

// Reduce - reduces the mask to the specified degree and removes matches
func Reduce(mask PathsProvider, degree ReduceDegree) *Mask {
	if mask == nil || degree == 0 {
		return &Mask{}
	}

	maskPaths := mask.GetPaths()
	reducedMask := &Mask{
		p: make([]string, len(maskPaths)),
	}

	copy(reducedMask.p, maskPaths)

	reducedMask.reduce(int(degree))
	return reducedMask
}

type Mask struct {
	p []string
}

func (m *Mask) GetPaths() []string {
	return m.p
}

func (m *Mask) reduce(degree int) {
	maskPaths := m.GetPaths()
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

	m.removeDuplicates()
}

func (m *Mask) removeDuplicates() {
	if len(m.p) < 1 {
		return
	}

	sort.Strings(m.p)
	j := 1
	for i := 1; i < len(m.p); i++ {
		if m.p[i-1] != m.p[i] {
			m.p[j] = m.p[i]
			j++
		}
	}

	m.p = m.p[:j]
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
