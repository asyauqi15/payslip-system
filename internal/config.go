package internal

import (
	"encoding/base64"
	"time"
)

type Config struct {
	Database   DatabaseConfig   `mapstructure:"database"`
	HTTPServer HTTPServerConfig `mapstructure:"http_server"`
}

type DatabaseConfig struct {
	Source string `mapstructure:"source"`
}

type HTTPServerConfig struct {
	Port                      int           `mapstructure:"port"`
	AccessTokenSecretEncoded  string        `mapstructure:"access_token_secret_encoded"`
	RefreshTokenSecretEncoded string        `mapstructure:"refresh_token_secret_encoded"`
	AccessTokenDuration       time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration      time.Duration `mapstructure:"refresh_token_duration"`
	ReadHeaderTimeout         time.Duration `mapstructure:"read_header_timeout"`
	ReadTimeout               time.Duration `mapstructure:"read_timeout"`
	IdleTimeout               time.Duration `mapstructure:"idle_timeout"`
	WriteTimeout              time.Duration `mapstructure:"write_timeout"`
}

func (h *HTTPServerConfig) GetAccessTokenSecret() ([]byte, error) {
	return base64.StdEncoding.DecodeString(h.AccessTokenSecretEncoded)
}

func (h *HTTPServerConfig) GetRefreshTokenSecret() ([]byte, error) {
	return base64.StdEncoding.DecodeString(h.RefreshTokenSecretEncoded)
}
