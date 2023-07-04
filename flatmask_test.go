package flatmask

import (
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

func TestReduceFlatMasksChild(t *testing.T) {
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

func TestReduceFieldmasksChild(t *testing.T) {
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

func TestReduceFieldmasksZero(t *testing.T) {
	table := []fieldmaskTable{
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"a.b.c", "a.b.c.d", "g.awa"},
			},
			expected: []string{},
		},
	}

	for _, tt := range table {
		rm := Reduce(tt.mask, 0)
		if !testEq(rm.GetPaths(), tt.expected) {
			t.Errorf(
				"Reduce(%v): actual - %v, exptected - %v",
				tt.mask.GetPaths(), rm.GetPaths(), tt.expected,
			)
		}
	}
}

func TestReduceFieldmasksTotal(t *testing.T) {
	table := []fieldmaskTable{
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"a.b.c", "a.b.c.d", "g.awa"},
			},
			expected: []string{"a", "g"},
		},
	}

	for _, tt := range table {
		rm := Reduce(tt.mask, Total)
		if !testEq(rm.GetPaths(), tt.expected) {
			t.Errorf(
				"Reduce(%v): actual - %v, exptected - %v",
				tt.mask.GetPaths(), rm.GetPaths(), tt.expected,
			)
		}
	}
}

func TestReduceFieldmasksChildAfterRoot(t *testing.T) {
	table := []fieldmaskTable{
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"a.b.c", "a.b.c.d", "g.awa"},
			},
			expected: []string{"a.b", "g.awa"},
		},
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"1.12.512", "1", "2.52", "2.81", "2.52.3"},
			},
			expected: []string{"1", "1.12", "2.52", "2.81"},
		},
	}

	for _, tt := range table {
		_ = Reduce(tt.mask, Total)
		rm := Reduce(tt.mask, RootChild)

		if !testEq(rm.GetPaths(), tt.expected) {
			t.Errorf(
				"Reduce(%v): actual - %v, exptected - %v",
				tt.mask.GetPaths(), rm.GetPaths(), tt.expected,
			)
		}
	}
}

func TestReduceFieldmasksSecondChild(t *testing.T) {
	table := []fieldmaskTable{
		{
			mask: &fieldmaskpb.FieldMask{
				Paths: []string{"a.b.c.d.e", "a.b.c.d", "g.a", "g.a.c", "g.a.c.1"},
			},
			expected: []string{"a.b.c", "g.a", "g.a.c"},
		},
	}

	for _, tt := range table {
		rm := Reduce(tt.mask, 3)

		if !testEq(rm.GetPaths(), tt.expected) {
			t.Errorf(
				"Reduce(%v): actual - %v, exptected - %v",
				tt.mask.GetPaths(), rm.GetPaths(), tt.expected,
			)
		}
	}
}

func TestReduceFieldmasksNil(t *testing.T) {
	table := []fieldmaskTable{
		{
			mask:     nil,
			expected: []string{},
		},
	}

	for _, tt := range table {
		rm := Reduce(tt.mask, 3)

		if !testEq(rm.GetPaths(), tt.expected) {
			t.Errorf(
				"Reduce(%v): actual - %v, exptected - %v",
				tt.mask.GetPaths(), rm.GetPaths(), tt.expected,
			)
		}
	}
}


