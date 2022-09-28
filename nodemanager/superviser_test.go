package nodemanager

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuperviser_lastBlockSeenLogPlugin(t *testing.T) {
	tests := []struct {
		name string
		line string
		want uint64
	}{
		{"block line zero", "FIRE BLOCK 0 5b02121274e67b59671b7e6c3711cc74", 0},
		{"block line some", "FIRE BLOCK 10 5b02121274e67b59671b7e6c3711cc74", 10},
		{"block line max", "FIRE BLOCK 18446744073709551615 5b02121274e67b59671b7e6c3711cc74", uint64(math.MaxUint64)},

		// Only logs an error
		{"block line missing data", "FIRE BLOCK 10", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Superviser{
				Logger: zlog,
			}

			s.lastBlockSeenLogPlugin(tt.line)
			require.Equal(t, tt.want, s.lastBlockSeen)
		})
	}
}
