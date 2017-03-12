package impl

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		in  int
		out *MarkovChain
		ok  bool
	}{
		{
			in: 2,
			out: &MarkovChain{
				order: 2,
			},
			ok: true,
		},
		{
			in: 47,
			out: &MarkovChain{
				order: 47,
			},
			ok: true,
		},
		{
			in: 1,
			out: &MarkovChain{
				order: 1,
			},
			ok: true,
		},
		{}, // The all-zero case is the test for order = 0.
		{
			in: -1,
		},
	}

	for _, tc := range tests {
		out, err := New(tc.in)
		if err == nil && !tc.ok {
			t.Errorf("New(%d) = %+v, want err", tc.in, out)
			continue
		}
		if err != nil && tc.ok {
			t.Errorf("New(%d) = %s, want %+v", tc.in, err, tc.out)
			continue
		}
		if !reflect.DeepEqual(out, tc.out) && tc.ok {
			t.Errorf("New(d) = %+v, want %+v", tc.in, out, tc.out)
		}
	}
}

func TestTrain(t *testing.T) {
	type input struct {
		r     string // For convenience of pretty test outputs.
		order int
	}
	tests := []struct {
		in  input
		out *MarkovChain
		ok  bool
	}{
		// Test vacuous case of no input data.
		{
			in: input{
				order: 2,
			},
			out: &MarkovChain{
				order: 2,
			},
			ok: true,
		},
		// Simple order = 1 case.
		{
			in: input{
				r:     "I see a tree built into the sidewalk.",
				order: 1,
			},
			out: &MarkovChain{
				order: 1,
				lessons: []lesson{
					lesson{
						back: []string{""},
						next: "I",
					},
					lesson{
						back: []string{"I"},
						next: "see",
					},
					lesson{
						back: []string{"see"},
						next: "a",
					},
					lesson{
						back: []string{"a"},
						next: "tree",
					},
					lesson{
						back: []string{"tree"},
						next: "built",
					},
					lesson{
						back: []string{"built"},
						next: "into",
					},
					lesson{
						back: []string{"into"},
						next: "the",
					},
					lesson{
						back: []string{"the"},
						next: "sidewalk.",
					},
				},
			},
			ok: true,
		},
		// Simple order = 2 case.
		{
			in: input{
				r:     "The tree has two posts supporting it.",
				order: 2,
			},
			out: &MarkovChain{
				order: 2,
				lessons: []lesson{
					lesson{
						back: []string{"", ""},
						next: "The",
					},
					lesson{
						back: []string{"", "The"},
						next: "tree",
					},
					lesson{
						back: []string{"The", "tree"},
						next: "has",
					},
					lesson{
						back: []string{"tree", "has"},
						next: "two",
					},
					lesson{
						back: []string{"has", "two"},
						next: "posts",
					},
					lesson{
						back: []string{"two", "posts"},
						next: "supporting",
					},
					lesson{
						back: []string{"posts", "supporting"},
						next: "it.",
					},
				},
			},
			ok: true,
		},
		// Simple order = 3 case.
		{
			in: input{
				r:     "There's a coat on one post.",
				order: 3,
			},
			out: &MarkovChain{
				order: 3,
				lessons: []lesson{
					{
						back: []string{"", "", ""},
						next: "There's",
					},
					{
						back: []string{"", "", "There's"},
						next: "a",
					},
					{
						back: []string{"", "There's", "a"},
						next: "coat",
					},
					{
						back: []string{"There's", "a", "coat"},
						next: "on",
					},
					{
						back: []string{"a", "coat", "on"},
						next: "one",
					},
					{
						back: []string{"coat", "on", "one"},
						next: "post.",
					},
				},
			},
			ok: true,
		},
	}

	for _, tc := range tests {
		// Set up.
		m, err := New(tc.in.order)
		if err != nil {
			t.Errorf("New(%d) got err, want nil (New is not SUT)", tc.in.order)
			continue
		}

		r := strings.NewReader(tc.in.r)

		err = m.Train(r)

		if err == nil && !tc.ok {
			t.Errorf("New(%d).Train(%s) = %+v, want err", tc.in.order, tc.in.r, m)
			continue
		}
		if err != nil && tc.ok {
			t.Errorf("New(%d).Train(%s) = %s, want %+v", tc.in.order, tc.in.r, err, tc.out)
			continue
		}
		if !reflect.DeepEqual(m, tc.out) && tc.ok {
			t.Errorf("New(%d).Train(%s) = %+v, want %+v", tc.in.order, tc.in.r, m, tc.out)
		}
	}

	// Test case with nil reader. Cannot be captured in above test loop.
	order := 2
	m, err := New(order)
	if err != nil {
		t.Errorf("New(%d) got err, want nil (New is not SUT)", 2)
		return
	}

	var r io.Reader
	if err := m.Train(r); err == nil {
		t.Errorf("New(%d).Train(nil reader) = %+v, want err", order, m)
	}
}
