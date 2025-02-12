package config

import "embed"

type LogLevel int8

const (
	DebugLevel LogLevel = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

type Config struct {
	Address          string
	ConnectionString string
	SigningKeyPath   string
	JWTIssuer        string
	Log              LogLevel
	Migrations       embed.FS
}

func New(
	address string,
	connStr string,
	signingKeyPath string,
	jwtIssuer string,
	logLevel LogLevel,
	migrations embed.FS,
) (*Config, error) {
	return &Config{
		Address:          address,
		ConnectionString: connStr,
		SigningKeyPath:   signingKeyPath,
		JWTIssuer:        jwtIssuer,
		Log:              logLevel,
		Migrations:       migrations,
	}, nil
}
