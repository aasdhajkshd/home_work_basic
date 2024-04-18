package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountLogLevels(t *testing.T) {
	cases := []struct {
		testName      string
		levelStr      string
		expectedCount int
		expectedError bool
	}{
		{
			testName:      "Level info",
			levelStr:      "info",
			expectedCount: 147,
			expectedError: false,
		},
		{
			testName:      "Level warn",
			levelStr:      "warn",
			expectedCount: 4,
			expectedError: false,
		},
		{
			testName:      "Level error",
			levelStr:      "trace",
			expectedCount: 119,
			expectedError: false,
		},
		{
			testName:      "Level unknown",
			levelStr:      "panic",
			expectedCount: 119,
			expectedError: false,
		},
	}
	inputLog, err := readLogFile("example.log")
	if err != nil {
		require.FailNow(t, "whoops, check the file name: %s", err)
	}
	level := &LogLevel{}
	for _, j := range cases {
		t.Run(j.testName, func(t *testing.T) {
			level.setLevel(j.levelStr)
			resultCount := countLogLevels(inputLog, *level)
			require.Equal(t, j.expectedCount, resultCount)
			if j.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
