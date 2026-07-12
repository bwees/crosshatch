package config

import "os"

// DataDir returns the directory where persistent data (the SQLite database) is
// stored. Defaults to the current directory for local development; production
// sets DATA_DIR to a mounted volume.
func DataDir() string {
	dir := os.Getenv("DATA_DIR")
	if dir == "" {
		return "."
	}
	return dir
}
