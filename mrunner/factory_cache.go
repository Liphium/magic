package mrunner

import (
	"path/filepath"
)

// Get the cache directory
func (f Factory) CacheDirectory() string {
	return filepath.Join(f.mDir, "cache")
}

// Get the directory of a script in the cache directory.
func (f Factory) ScriptCacheDirectory(script string) string {
	return filepath.Join(f.CacheDirectory(), "script_"+script)
}

// Get the directory of a test in the cache directory.
func (f Factory) TestCacheDirectory(script string) string {
	return filepath.Join(f.CacheDirectory(), "test_"+script)
}

// Get the directory of a config in the cache directory.
func (f Factory) ConfigCacheDirectory(config string, profile string) string {
	return filepath.Join(f.CacheDirectory(), config+"_"+profile)
}
