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
		out MarkovChain
		ok  bool
	}{}

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
