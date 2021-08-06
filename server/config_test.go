package server

import "testing"

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		name string
		in   Config
		out  error
	}{
		{
			name: "when Environment is invalid",
			in:   NewConfig("", "8080", "connectionstring", "/logs.txt", 10, 42),
			out:  ErrEnvironmentConfigIsInvalid,
		},
		{
			name: "when Port is invalid",
			in:   NewConfig("development", "", "connectionstring", "/logs.txt", 10, 42),
			out:  ErrPortConfigIsInvalid,
		},
		{
			name: "when DatabaseURL is invalid",
			in:   NewConfig("development", "8080", "", "/logs.txt", 10, 42),
			out:  ErrDatabaseURLConfigIsInvalid,
		},
		{
			name: "when LogsPath is invalid",
			in:   NewConfig("development", "8080", "connectionstring", "", 10, 42),
			out:  ErrLogsPathConfigIsInvalid,
		},
		{
			name: "when PaginationDefaultLimit is invalid",
			in:   NewConfig("development", "8080", "connectionstring", "/logs.txt", 0, 42),
			out:  ErrPaginationDefaultLimitConfigIsInvalid,
		},
		{
			name: "when PaginationMaxLimit is invalid",
			in:   NewConfig("development", "8080", "connectionstring", "/logs.txt", 10, 0),
			out:  ErrPaginationMaxLimitConfigIsInvalid,
		},
		{
			name: "when PaginationMaxLimit is invalid",
			in:   NewConfig("development", "8080", "connectionstring", "/logs.txt", 10, 9),
			out:  ErrPaginationMaxLimitConfigIsInvalid,
		},
		{
			name: "when success",
			in:   NewConfig("development", "8080", "connectionstring", "/logs.txt", 10, 42),
			out:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if r := tc.in.Validate(); r != tc.out {
				t.Errorf("was expecting %v, but returns %v", tc.out, r)
			}
		})
	}
}
