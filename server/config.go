package server

import "errors"

var (
	ErrEnvironmentConfigIsInvalid            = errors.New("the Environment config is invalid")
	ErrPortConfigIsInvalid                   = errors.New("the Port config is invalid")
	ErrDatabaseURLConfigIsInvalid            = errors.New("the DatabaseURL config is invalid")
	ErrLogsPathConfigIsInvalid               = errors.New("the LogsPath config is invalid")
	ErrPaginationDefaultLimitConfigIsInvalid = errors.New("the PaginationDefaultLimit config is invalid")
	ErrPaginationMaxLimitConfigIsInvalid     = errors.New("the PaginationMaxLimit config is invalid")
)

type Config struct {
	Environment            string
	Port                   string
	DatabaseURL            string
	LogsPath               string
	PaginationDefaultLimit int
	PaginationMaxLimit     int
}

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
