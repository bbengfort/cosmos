package enums_test

import (
	"testing"

	"github.com/bbengfort/cosmos/pkg/enums"
	"github.com/stretchr/testify/require"
)

func TestSize(t *testing.T) {

	t.Run("NumSystems", func(t *testing.T) {
		testCases := []struct {
			size enums.Size
			minn int
			maxn int
		}{
			{enums.Small, 20, 40},
			{enums.Medium, 100, 200},
			{enums.Large, 200, 400},
			{enums.Galactic, 500, 1000},
			{enums.Cosmic, 1000, 2000},
		}

		rounds := 100
		if testing.Short() {
			rounds = 5
		}

		for i, tc := range testCases {
			for j := 0; j < rounds; j++ {
				n := tc.size.NumSystems()
				require.GreaterOrEqual(t, n, tc.minn, "test case %d failed", i)
				require.LessOrEqual(t, n, tc.maxn, "test case %d failed", i)
			}
		}
	})

}
