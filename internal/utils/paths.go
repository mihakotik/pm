package utils

import (
	"os"
	"path/filepath"
)

const (
	DefaultLibDir     = "/var/lib/pm"
	DefaultCacheDir   = "/var/cache/pm"
	DefaultOptDir     = "/opt/pm"
	DefaultEtcDir     = "/etc/pm"
	DefaultBinDir     = "/usr/bin"
	DefaultDBName     = "installed.db"
	DefaultConfigName = "repos.conf"
)

// RootPrefix returns the optional sysroot prefix from PM_ROOT environment variable.
func RootPrefix() string {
	return os.Getenv("PM_ROOT")
}

// LibDir returns the state directory for pm (/var/lib/pm).
func LibDir() string {
	if root := RootPrefix(); root != "" {
		return filepath.Join(root, "var", "lib", "pm")
	}
	return DefaultLibDir
}

// CacheDir returns the download and package cache directory (/var/cache/pm).
func CacheDir() string {
	if root := RootPrefix(); root != "" {
		return filepath.Join(root, "var", "cache", "pm")
	}
	return DefaultCacheDir
}

// OptDir returns the installation directory for foreign packages (/opt/pm).
func OptDir() string {
	if root := RootPrefix(); root != "" {
		return filepath.Join(root, "opt", "pm")
	}
	return DefaultOptDir
}

// EtcDir returns configuration directory (/etc/pm).
func EtcDir() string {
	if root := RootPrefix(); root != "" {
		return filepath.Join(root, "etc", "pm")
	}
	return DefaultEtcDir
}

// BinDir returns the system binary directory for wrappers (/usr/bin).
func BinDir() string {
	if root := RootPrefix(); root != "" {
		return filepath.Join(root, "usr", "bin")
	}
	return DefaultBinDir
}

// DBPath returns the path to the installed packages SQLite database.
func DBPath() string {
	return filepath.Join(LibDir(), DefaultDBName)
}

// ConfigPath returns the path to the repository configuration file.
func ConfigPath() string {
	return filepath.Join(EtcDir(), DefaultConfigName)
}

// ProviderOptDir returns the path for a specific package in /opt/pm/<provider>/<pkg>.
func ProviderOptDir(provider, pkgName string) string {
	return filepath.Join(OptDir(), provider, pkgName)
}
