package impl

import (
	"reflect"
	"strings"
	"testing"
)

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
			out: &MarkovChain{},
			ok:  true,
		},
		// Simple order = 1 case.
		{
			in: input{
				r:     "I see a tree built into the sidewalk.",
				order: 1,
			},
			out: &MarkovChain{
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
		m := new(MarkovChain)
		r := strings.NewReader(tc.in.r)

		err := m.Train(r, tc.in.order)

		if err == nil && !tc.ok {
			t.Errorf("Train(%s, %d) = %+v, want err", tc.in.r, tc.in.order, m)
			continue
		}
		if err != nil && tc.ok {
			t.Errorf("Train(%s, %d) = %s, want %+v", tc.in.r, tc.in.order, err, tc.out)
			continue
		}
		if !reflect.DeepEqual(m, tc.out) {
			t.Errorf("Train(%s, %d) = %+v, want %+v", tc.in.r, tc.in.order, m, tc.out)
		}
	}
}
