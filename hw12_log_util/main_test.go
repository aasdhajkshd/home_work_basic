package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLog(t *testing.T) {
	//nolint: lll
	inputLog := []string{
		`79.105.72.219:18446 [14/Apr/2024:06:33:50.930] fe_https fe_https/<NOSRV> 0/-1/-1/-1/0 301 100 - - LR-- 1/1/0/0/0 0/0 "GET / HTTP/1.1"`,
		`78.153.140.179:54722 [14/Apr/2024:19:12:45.653] fe_https~ be_vkworkmail_https/vkworkmail 0/0/13/1/14 200 177 - - ---- 2/2/0/0/0 0/0 "GET /.aws/credentials HTTP/1.1"`,
		`162.243.136.67:45650 [14/Apr/2024:19:16:37.221] fe_https/HTTPS: SSL handshake failure`,
		`51.195.216.60:34370 [14/Apr/2024:19:18:44.190] fe_https fe_https/<NOSRV> 0/-1/-1/-1/0 403 197 - - PR-- 1/1/0/0/0 0/0 "GET /wp-login.php HTTP/1.1"`,
	}
	cases := []struct {
		testName      string
		level         logLevel
		expectedLog   []string
		expectedError bool
	}{
		{
			testName:      "Level info",
			level:         Info,
			expectedLog:   []string{"Status code 200: OK: 1\n", "Status code 301: Moved Permanently: 1\n", "Status code 403: Forbidden: 1\n"}, //nolint: lll
			expectedError: false,
		},
		{
			testName:      "Level warn",
			level:         Warn,
			expectedLog:   []string{"2024-04-14 06:33:50.93 +0000 UTC 79.105.72.219 fe_https /\n"},
			expectedError: false,
		},
		{
			testName:      "Level error",
			level:         Error,
			expectedLog:   []string{"2024-04-14 19:18:44.19 +0000 UTC 51.195.216.60 fe_https /wp-login.php\n"},
			expectedError: false,
		},
	}
	for _, j := range cases {
		t.Run(j.testName, func(t *testing.T) {
			parsedLog, err := parseLog(inputLog, j.level)
			require.Equal(t, j.expectedLog, parsedLog)
			if j.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
