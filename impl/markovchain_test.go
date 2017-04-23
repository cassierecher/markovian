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
				Order: 2,
			},
			ok: true,
		},
		{
			in: 47,
			out: &MarkovChain{
				Order: 47,
			},
			ok: true,
		},
		{
			in: 1,
			out: &MarkovChain{
				Order: 1,
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

func TestTrain_New(t *testing.T) {
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
				Order: 2,
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
				Order: 1,
				Lessons: []lesson{
					lesson{
						Back: []string{""},
						Next: "I",
					},
					lesson{
						Back: []string{"I"},
						Next: "see",
					},
					lesson{
						Back: []string{"see"},
						Next: "a",
					},
					lesson{
						Back: []string{"a"},
						Next: "tree",
					},
					lesson{
						Back: []string{"tree"},
						Next: "built",
					},
					lesson{
						Back: []string{"built"},
						Next: "into",
					},
					lesson{
						Back: []string{"into"},
						Next: "the",
					},
					lesson{
						Back: []string{"the"},
						Next: "sidewalk",
					},
					lesson{
						Back: []string{"sidewalk"},
						Next: ".",
					},
				},
			},
			ok: true,
		},
		// Simple order = 2 case.
		{
			in: input{
				r:     "The tree has two posts supporting it!",
				order: 2,
			},
			out: &MarkovChain{
				Order: 2,
				Lessons: []lesson{
					lesson{
						Back: []string{"", ""},
						Next: "The",
					},
					lesson{
						Back: []string{"", "The"},
						Next: "tree",
					},
					lesson{
						Back: []string{"The", "tree"},
						Next: "has",
					},
					lesson{
						Back: []string{"tree", "has"},
						Next: "two",
					},
					lesson{
						Back: []string{"has", "two"},
						Next: "posts",
					},
					lesson{
						Back: []string{"two", "posts"},
						Next: "supporting",
					},
					lesson{
						Back: []string{"posts", "supporting"},
						Next: "it",
					},
					lesson{
						Back: []string{"supporting", "it"},
						Next: "!",
					},
				},
			},
			ok: true,
		},
		// Simple order = 3 case.
		{
			in: input{
				r:     "There's a coat on one post?",
				order: 3,
			},
			out: &MarkovChain{
				Order: 3,
				Lessons: []lesson{
					{
						Back: []string{"", "", ""},
						Next: "There's",
					},
					{
						Back: []string{"", "", "There's"},
						Next: "a",
					},
					{
						Back: []string{"", "There's", "a"},
						Next: "coat",
					},
					{
						Back: []string{"There's", "a", "coat"},
						Next: "on",
					},
					{
						Back: []string{"a", "coat", "on"},
						Next: "one",
					},
					{
						Back: []string{"coat", "on", "one"},
						Next: "post",
					},
					{
						Back: []string{"on", "one", "post"},
						Next: "?",
					},
				},
			},
			ok: true,
		},
		{
			in: input{
				r:     "floating punctuation . works",
				order: 1,
			},
			out: &MarkovChain{
				Order: 1,
				Lessons: []lesson{
					{
						Back: []string{""},
						Next: "floating",
					},
					{
						Back: []string{"floating"},
						Next: "punctuation",
					},
					{
						Back: []string{"punctuation"},
						Next: ".",
					},
					{
						Back: []string{"."},
						Next: "works",
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

		err = m.Train(strings.NewReader(tc.in.r))

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

func TestTrain_Existing(t *testing.T) {
	type input struct {
		r string // For convenience of pretty test outputs.
		m *MarkovChain
	}
	tests := []struct {
		in  input
		out *MarkovChain
	}{
		// Test case of training a premade but empty chain.
		{
			in: input{
				r: "The dog",
				m: &MarkovChain{
					Order:   2,
					Lessons: []lesson{},
				},
			},
			out: &MarkovChain{
				Order: 2,
				Lessons: []lesson{
					lesson{
						Back: []string{"", ""},
						Next: "The",
					},
					lesson{
						Back: []string{"", "The"},
						Next: "dog",
					},
				},
			},
		},
		// Test case of no input data to premade Markov chain.
		{
			in: input{
				m: &MarkovChain{
					Order: 2,
					Lessons: []lesson{
						lesson{
							Back: []string{"red", "orange"},
							Next: "yellow",
						},
					},
				},
			},
			out: &MarkovChain{
				Order: 2,
				Lessons: []lesson{
					lesson{
						Back: []string{"red", "orange"},
						Next: "yellow",
					},
				},
			},
		},
		// Test case of training a premade Markov chain.
		{
			in: input{
				m: &MarkovChain{
					Order: 1,
					Lessons: []lesson{
						lesson{
							Back: []string{"green"},
							Next: "blue",
						},
						lesson{
							Back: []string{"blue"},
							Next: "indigo",
						},
						lesson{
							Back: []string{"indigo"},
							Next: "violet",
						},
					},
				},
				r: "translucent clear iridescent",
			},
			out: &MarkovChain{
				Order: 1,
				Lessons: []lesson{
					lesson{
						Back: []string{"green"},
						Next: "blue",
					},
					lesson{
						Back: []string{"blue"},
						Next: "indigo",
					},
					lesson{
						Back: []string{"indigo"},
						Next: "violet",
					},
					lesson{
						Back: []string{""},
						Next: "translucent",
					},
					lesson{
						Back: []string{"translucent"},
						Next: "clear",
					},
					lesson{
						Back: []string{"clear"},
						Next: "iridescent",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		// Copy.
		m := &MarkovChain{
			Order: tc.in.m.Order,
		}
		// Copy demands the lengths sync up.
		m.Lessons = make([]lesson, len(tc.in.m.Lessons))
		copy(m.Lessons, tc.in.m.Lessons)

		err := m.Train(strings.NewReader(tc.in.r))

		if err != nil {
			t.Errorf("%+v.Train(%s) = %s, want %+v", tc.in.m, tc.in.r, err, tc.out)
			continue
		}
		if !reflect.DeepEqual(m, tc.out) {
			t.Errorf("%+v.Train(%s) = %+v, want %+v", tc.in.m, tc.in.r, m, tc.out)
		}
	}

	// Test case with nil reader. Cannot be captured in above test loop.
	m := &MarkovChain{
		Order: 2,
		Lessons: []lesson{
			lesson{
				Back: []string{"one", "two"},
				Next: "three",
			},
		},
	}
	m2 := &MarkovChain{
		Order: m.Order,
	}
	// Copy demands the lengths sync up.
	m2.Lessons = make([]lesson, len(m.Lessons))
	copy(m2.Lessons, m.Lessons)

	var r io.Reader
	if err := m2.Train(r); err == nil {
		t.Errorf("%+v.Train(nil reader) = %+v, want err", m, m2)
	}
}

func TestBuildKey(t *testing.T) {
	tests := []struct {
		in  []string
		out string
	}{
		// Empty case.
		{},
		// One simple string.
		{
			in:  []string{"antiquing"},
			out: "antiquing",
		},
		// Multiple simple strings.
		{
			in: []string{
				"backyard",
				"maple",
				"tree",
			},
			out: "backyard$maple$tree",
		},
		// Add delimiters.
		{
			in:  []string{"ca$h"},
			out: `ca\$h`,
		},
		{
			in: []string{
				"sea$hells",
				"$eem",
				"super",
			},
			out: `sea\$hells$\$eem$super`,
		},
		// Add escapes.
		{
			in:  []string{`back\slash`},
			out: `back\\slash`,
		},
		{
			in:  []string{`\\\`},
			out: `\\\\\\`,
		},
		// Delimiters and escapes.
		{
			in: []string{
				`$$\$\\`,
				`fan\ta$tic`,
				`stuff`,
			},
			out: `\$\$\\\$\\\\$fan\\ta\$tic$stuff`,
		},
	}
	for _, tc := range tests {
		out := buildKey(tc.in)
		if out != tc.out {
			t.Errorf("buildKey(%v) = %s, want %s", tc.in, out, tc.out)
		}
	}
}
