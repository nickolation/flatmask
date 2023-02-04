package flatmask

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type maskTable struct {
	p        []string
	expected []string
}

func (m *maskTable) GetPaths() []string {
	return m.p
}

func (m *maskTable) Reset() {
}

func (m *maskTable) Normalize() {
}

func (m *maskTable) String() string {
	return strings.Join(m.p, ".")
}

func TestReduceFlatMasks(t *testing.T) {
	table := []*maskTable{
		{
			p:        []string{"a.b.c"},
			expected: []string{"a.b"},
		},
		{
			p:        []string{"a.b.c", "a.b"},
			expected: []string{"a.b"},
		},
		{
			p:        []string{"a.b.c.d", "a.f.g"},
			expected: []string{"a.b", "a.f"},
		},
	}

	for _, tt := range table {
		rm := Reduce(tt, RootChild)
		if !testEq(rm.GetPaths(), tt.expected) {
			t.Errorf(
				"Reduce(%v): actual - %v, exptected - %v",
				tt.GetPaths(), rm.GetPaths(), tt.expected,
			)
		}
	}

}

type fieldmaskTable struct {
	mask     *fieldmaskpb.FieldMask
	expected []string
}

func TestReduceFieldmasks(t *testing.T) {
	table := []fieldmaskTable{
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"a.b.c"},
			},
			expected: []string{"a.b"},
		},
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"a.b.c", "a.b.c.d"},
			},
			expected: []string{"a.b"},
		},
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"a.b.c", "a.f.d.e"},
			},
			expected: []string{"a.b", "a.f"},
		},
	}

	for _, tt := range table {
		rm := Reduce(tt.mask, RootChild)
		if !testEq(rm.GetPaths(), tt.expected) {
			t.Errorf(
				"Reduce(%v): actual - %v, exptected - %v",
				tt.mask.GetPaths(), rm.GetPaths(), tt.expected,
			)
		}
	}
}
