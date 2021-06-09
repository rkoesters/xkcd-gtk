// Package paths provides information on where to find and store files.
package paths

var b Builder

// Init initializes the paths package. Must be called before calling other
// package functions.
func Init(appID string) {
	b = Builder{
		appID: appID,
	}
}

// CacheDir returns the path to our app's cache directory.
func CacheDir() string {
	return b.CacheDir()
}

// EnsureCacheDir creates the CacheDir() directory, if it doesn't exist. It does
// not return an error if it already exists.
func EnsureCacheDir() error {
	return b.EnsureCacheDir()
}

// ConfigDir returns the path to our app's user configuration directory.
func ConfigDir() string {
	return b.ConfigDir()
}

// EnsureConfigDir creates the ConfigDir() directory, if it doesn't exist. It
// does not return an error if it already exists.
func EnsureConfigDir() error {
	return b.EnsureConfigDir()
}

// DataDir returns the path to our app's user data directory.
func DataDir() string {
	return b.DataDir()
}

// EnsureDataDir creates the DataDir() directory, if it doesn't exist. It does
// not return an error if it already exists.
func EnsureDataDir() error {
	return b.EnsureDataDir()
}

// LocaleDir returns the path to the system locale directory.
func LocaleDir() string {
	return b.LocaleDir()
}
