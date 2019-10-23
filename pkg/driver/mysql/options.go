package mysql

import (
	"time"
)

type DatabaseOption func(config *Config)

// MaxConnLifetime
func WithMaxConnLifetime(lifetime time.Duration) DatabaseOption {
	return func(config *Config) {
		config.MaxLifetime = lifetime
	}
}

// MaxOpenConns
func WithMaxOpenConns(limit int) DatabaseOption {
	return func(config *Config) {
		config.MaxOpenConns = limit
	}
}

// MaxIdleConns
func WithMaxIdleConns(limit int) DatabaseOption {
	return func(config *Config) {
		config.MaxIdleConns = limit
	}
}

// MaxAllowedPacket
func WithMaxAllowedPacket(limit int) DatabaseOption {
	return func(config *Config) {
		config.MaxAllowedPacket = limit
	}
}

// AllowOldPasswords
func WithAllowOldPasswords(b bool) DatabaseOption {
	return func(config *Config) {
		config.AllowOldPasswords = b
	}
}

// AllowCleartextPasswords
func WithAllowCleartextPasswords(b bool) DatabaseOption {
	return func(config *Config) {
		config.AllowCleartextPasswords = b
	}
}

// AllowNativePasswords
func WithAllowNativePasswords(b bool) DatabaseOption {
	return func(config *Config) {
		config.AllowNativePasswords = b
	}
}
