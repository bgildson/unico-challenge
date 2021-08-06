package server

import "errors"

var (
	// ErrEnvironmentConfigIsInvalid is used to represent an error in the Environment config
	ErrEnvironmentConfigIsInvalid = errors.New("the Environment config is invalid")
	// ErrPortConfigIsInvalid is used to represent an error in the Port config
	ErrPortConfigIsInvalid = errors.New("the Port config is invalid")
	// ErrDatabaseURLConfigIsInvalid is used to represent an error in the DatabaseURL config
	ErrDatabaseURLConfigIsInvalid = errors.New("the DatabaseURL config is invalid")
	// ErrLogsPathConfigIsInvalid is used to represent an error in the LogsPath config
	ErrLogsPathConfigIsInvalid = errors.New("the LogsPath config is invalid")
	// ErrPaginationDefaultLimitConfigIsInvalid is used to represent an error in the PaginationDefaultLimit config
	ErrPaginationDefaultLimitConfigIsInvalid = errors.New("the PaginationDefaultLimit config is invalid")
	// ErrPaginationMaxLimitConfigIsInvalid is used to represent an error in the PaginationMaxLimit config
	ErrPaginationMaxLimitConfigIsInvalid = errors.New("the PaginationMaxLimit config is invalid")
)

// Config represents the server config
type Config struct {
	Environment            string
	Port                   string
	DatabaseURL            string
	LogsPath               string
	PaginationDefaultLimit int
	PaginationMaxLimit     int
}

// NewConfig creates a Config for the server
func NewConfig(environment, port, databaseURL, logsPath string, paginationDefaultLimit, paginationMaxLimit int) Config {
	return Config{
		Environment:            environment,
		Port:                   port,
		DatabaseURL:            databaseURL,
		LogsPath:               logsPath,
		PaginationDefaultLimit: paginationDefaultLimit,
		PaginationMaxLimit:     paginationMaxLimit,
	}
}

// Validate applies the validation for every field in config
func (c Config) Validate() error {
	if c.Environment == "" {
		return ErrEnvironmentConfigIsInvalid
	}

	if c.Port == "" {
		return ErrPortConfigIsInvalid
	}

	if c.DatabaseURL == "" {
		return ErrDatabaseURLConfigIsInvalid
	}

	if c.LogsPath == "" {
		return ErrLogsPathConfigIsInvalid
	}

	if c.PaginationDefaultLimit == 0 {
		return ErrPaginationDefaultLimitConfigIsInvalid
	}

	if c.PaginationMaxLimit == 0 || c.PaginationMaxLimit < c.PaginationDefaultLimit {
		return ErrPaginationMaxLimitConfigIsInvalid
	}

	return nil
}
