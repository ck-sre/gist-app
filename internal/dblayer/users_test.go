package dblayer

import (
	"gistapp.ck89.net/internal/assert"
	"testing"
)

func TestUserLayer_CheckExists(t *testing.T) {
	if testing.Short() {
		t.Skip("DBLayer: skipping integration test in short mode.")
	}
	tsts := []struct {
		name     string
		userId   int
		expected bool
	}{
		{
			name:     "user exists",
			userId:   1,
			expected: true,
		},
		{
			name:     "user does not exist",
			userId:   2,
			expected: false,
		},
		{
			name:     "Zero user id",
			userId:   0,
			expected: false,
		},
	}

	for _, tst := range tsts {
		t.Run(tst.name, func(t *testing.T) {
			db := newTstDB(t)
			n := UserLayer{db}

			exists, err := n.CheckExists(tst.userId)
			assert.Same(t, exists, tst.expected)
			assert.NilError(t, err)
		})
	}
}
