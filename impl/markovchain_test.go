package impl

import (
	"fmt"
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
				Order:     2,
				Knowledge: map[string]frequencyGroup{},
			},
			ok: true,
		},
		{
			in: 47,
			out: &MarkovChain{
				Order:     47,
				Knowledge: map[string]frequencyGroup{},
			},
			ok: true,
		},
		{
			in: 1,
			out: &MarkovChain{
				Order:     1,
				Knowledge: map[string]frequencyGroup{},
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
				Order:     2,
				Knowledge: map[string]frequencyGroup{},
			},
			ok: true,
		},
		// Simple order = 1 case.
		{
			in: input{
				r:     "I see the tree built into the sidewalk.",
				order: 1,
			},
			out: &MarkovChain{
				Order: 1,
				Knowledge: map[string]frequencyGroup{
					"": frequencyGroup{
						"I": 1,
					},
					"I": frequencyGroup{
						"see": 1,
					},
					"see": frequencyGroup{
						"the": 1,
					},
					"the": frequencyGroup{
						"tree":     1,
						"sidewalk": 1,
					},
					"tree": frequencyGroup{
						"built": 1,
					},
					"built": frequencyGroup{
						"into": 1,
					},
					"into": frequencyGroup{
						"the": 1,
					},
					"sidewalk": frequencyGroup{
						".": 1,
					},
				},
			},
			ok: true,
		},
		// Simple order = 2 case.
		{
			in: input{
				r:     "The tree has two posts has two posts supporting it!",
				order: 2,
			},
			out: &MarkovChain{
				Order: 2,
				Knowledge: map[string]frequencyGroup{
					"$": frequencyGroup{
						"The": 1,
					},
					"$The": frequencyGroup{
						"tree": 1,
					},
					"The$tree": frequencyGroup{
						"has": 1,
					},
					"tree$has": frequencyGroup{
						"two": 1,
					},
					"has$two": frequencyGroup{
						"posts": 2,
					},
					"two$posts": frequencyGroup{
						"has":        1,
						"supporting": 1,
					},
					"posts$has": frequencyGroup{
						"two": 1,
					},
					"posts$supporting": frequencyGroup{
						"it": 1,
					},
					"supporting$it": frequencyGroup{
						"!": 1,
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
				Knowledge: map[string]frequencyGroup{
					"$$": frequencyGroup{
						"There's": 1,
					},
					"$$There's": frequencyGroup{
						"a": 1,
					},
					"$There's$a": frequencyGroup{
						"coat": 1,
					},
					"There's$a$coat": frequencyGroup{
						"on": 1,
					},
					"a$coat$on": frequencyGroup{
						"one": 1,
					},
					"coat$on$one": frequencyGroup{
						"post": 1,
					},
					"on$one$post": frequencyGroup{
						"?": 1,
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
				Knowledge: map[string]frequencyGroup{
					"": frequencyGroup{
						"floating": 1,
					},
					"floating": frequencyGroup{
						"punctuation": 1,
					},
					"punctuation": frequencyGroup{
						".": 1,
					},
					".": frequencyGroup{
						"works": 1,
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
					Order:     2,
					Knowledge: map[string]frequencyGroup{},
				},
			},
			out: &MarkovChain{
				Order: 2,
				Knowledge: map[string]frequencyGroup{
					"$": frequencyGroup{
						"The": 1,
					},
					"$The": frequencyGroup{
						"dog": 1,
					},
				},
			},
		},
		// Test case of no input data to premade Markov chain.
		{
			in: input{
				m: &MarkovChain{
					Order:     2,
					Knowledge: map[string]frequencyGroup{},
				},
			},
			out: &MarkovChain{
				Order:     2,
				Knowledge: map[string]frequencyGroup{},
			},
		},
		// Test case of training a premade Markov chain.
		{
			in: input{
				m: &MarkovChain{
					Order:     1,
					Knowledge: map[string]frequencyGroup{},
				},
				r: "translucent clear iridescent",
			},
			out: &MarkovChain{
				Order: 1,
				Knowledge: map[string]frequencyGroup{
					"": frequencyGroup{
						"translucent": 1,
					},
					"translucent": frequencyGroup{
						"clear": 1,
					},
					"clear": frequencyGroup{
						"iridescent": 1,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		in := tc.in

		// Save state for possibly printing later.
		mBeforeStr := fmt.Sprintf("%+v", in.m)

		err := in.m.Train(strings.NewReader(in.r))
		if err != nil {
			t.Errorf("(%s).Train(%s) = %s, want %+v", mBeforeStr, in.r, err, tc.out)
			continue
		}
		if !reflect.DeepEqual(in.m, tc.out) {
			t.Errorf("(%s).Train(%s) = %+v, want %+v", mBeforeStr, in.r, in.m, tc.out)
		}
	}

	// Test case with nil reader. Cannot be captured in above test loop.
	m := &MarkovChain{
		Order:     2,
		Knowledge: map[string]frequencyGroup{},
	}

	// Save state for possibly printing later.
	mBeforeStr := fmt.Sprintf("%+v", m)

	var r io.Reader
	if err := m.Train(r); err == nil {
		t.Errorf("%+v.Train(nil reader) = %+v, want err", mBeforeStr, m)
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

func testAddKnowledge(t *testing.T) {
	type input struct {
		m    *MarkovChain
		back []string
		next string
	}
	tests := []struct {
		in  input
		out *MarkovChain
	}{
		{
			in: input{
				m: &MarkovChain{
					Order:     2,
					Knowledge: map[string]frequencyGroup{},
				},
				back: []string{"alabama", "alaska"},
				next: "arizona",
			},
			out: &MarkovChain{
				Order: 2,
				Knowledge: map[string]frequencyGroup{
					"alabama$alaska": frequencyGroup{
						"arizona": 1,
					},
				},
			},
		},
		{
			in: input{
				m: &MarkovChain{
					Order: 2,
					Knowledge: map[string]frequencyGroup{
						"arkansas$california": frequencyGroup{
							"colorado": 2,
						},
					},
				},
				back: []string{"connecticut", "delaware"},
				next: "florida",
			},
			out: &MarkovChain{
				Order: 2,
				Knowledge: map[string]frequencyGroup{
					"arkansas$california": frequencyGroup{
						"colorado": 2,
					},
					"connecticut$delaware": frequencyGroup{
						"florida": 1,
					},
				},
			},
		},
		{
			in: input{
				m: &MarkovChain{
					Order: 2,
					Knowledge: map[string]frequencyGroup{
						"georgia$hawaii": frequencyGroup{
							"idaho":    4,
							"illinois": 6,
						},
					},
				},
				back: []string{"georgia", "hawaii"},
				next: "indiana",
			},
			out: &MarkovChain{
				Order: 2,
				Knowledge: map[string]frequencyGroup{
					"georgia$hawaii": frequencyGroup{
						"idaho":    4,
						"illinois": 6,
						"indiana":  1,
					},
				},
			},
		},
		{
			in: input{
				m: &MarkovChain{
					Order: 2,
					Knowledge: map[string]frequencyGroup{
						"iowa$kansas": frequencyGroup{
							"kentucky":  3,
							"louisiana": 1,
							"maine":     2,
						},
					},
				},
				back: []string{"iowa", "kansas"},
				next: "louisiana",
			},
			out: &MarkovChain{
				Order: 2,
				Knowledge: map[string]frequencyGroup{
					"iowa$kansas": frequencyGroup{
						"kentucky":  3,
						"louisiana": 2,
						"maine":     2,
					},
				},
			},
		},
	}
	for _, tc := range tests {
		in := tc.in
		// Save state for possibly printing later.
		mBeforeStr := fmt.Sprintf("%+v", in.m)

		in.m.addKnowledge(in.back, in.next)
		if !reflect.DeepEqual(in.m, tc.out) {
			t.Errorf("(%s).addKnowledge(%v, %s) yielded %+v, want %+v", mBeforeStr, in.back, in.next, in.m, tc.out)
		}
	}
}
