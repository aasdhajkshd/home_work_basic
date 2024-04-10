package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntCount(t *testing.T) {
	cases := []struct {
		name     string
		testIter int
		expected int
	}{
		{
			name:     "Test 1",
			testIter: 5,
			expected: 777777,
		},
	}

	wg := &sync.WaitGroup{}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			for i := v.testIter; i > 0; i-- {
				c := 0
				intCount(&c, wg, true)
				wg.Wait()
				require.Equal(t, v.expected, c)
			}
		})
	}
}
