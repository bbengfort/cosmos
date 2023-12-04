package jcode_test

import (
	"encoding/json"
	"testing"

	"github.com/bbengfort/cosmos/pkg/jcode"
	"github.com/stretchr/testify/require"
)

func TestJoinCode(t *testing.T) {

	t.Run("Random", func(t *testing.T) {
		nRuns := 512
		if testing.Short() {
			nRuns = 10
		}

		for i := 0; i < nRuns; i++ {
			// Test creating a random join code
			code := jcode.New()
			require.Len(t, code, 16)

			// Test stringify the new code
			s := code.String()
			require.Len(t, s, 19)
			require.Contains(t, s, "-")

			// Test marshaling the join code
			data, err := json.Marshal(code)
			require.NoError(t, err, "could not marshal join code")

			// Test unmarshaling the join code
			compat := jcode.JoinCode("")
			err = json.Unmarshal(data, &compat)
			require.NoError(t, err, "could not unmarshal join code")
			require.Equal(t, code, compat)
		}
	})

	t.Run("JSON", func(t *testing.T) {
		testCases := []struct {
			Data     []byte
			Expected jcode.JoinCode
		}{
			{[]byte(`"ABCD-1234-EFGH-5678"`), jcode.JoinCode("ABCD1234EFGH5678")},
			{[]byte(`"ABCD1234EFGH5678"`), jcode.JoinCode("ABCD1234EFGH5678")},
			{[]byte(`"abcd-1234-efgh-5678"`), jcode.JoinCode("ABCD1234EFGH5678")},
			{[]byte(`"abcd1234efgh5678"`), jcode.JoinCode("ABCD1234EFGH5678")},
			{[]byte(`"   ABCD-1234-EFGH-5678  "`), jcode.JoinCode("ABCD1234EFGH5678")},
		}

		for i, tc := range testCases {
			var code jcode.JoinCode
			err := json.Unmarshal(tc.Data, &code)
			require.NoError(t, err, "expected no error on test case %d", i)
			require.Equal(t, tc.Expected, code, "mismatch on test case %d", i)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		testCases := [][]byte{
			[]byte(`""`),
			[]byte(`"ABC12"`),
			[]byte(`"ABCD-1234-EFGH-OZIK"`),
			[]byte(`100`),
			[]byte(`{"join_code": "ABCD-1234-EFGH-5678"}`),
		}

		for i, tc := range testCases {
			var code jcode.JoinCode
			err := json.Unmarshal(tc, &code)
			require.ErrorIs(t, err, jcode.ErrInvalidJoinCode, "expected invalid error on test case %d", i)
		}
	})
}
