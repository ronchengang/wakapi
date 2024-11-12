package utils

import "fmt"

var (
	// These variables are intended to be set during build time using -ldflags
	CommitHash string
)

// GetVersion returns a formatted string containing both version and commit hash
func GetCommitHash() string {
	if CommitHash != "" {
		return fmt.Sprintf("%s ", CommitHash)
	}
	return CommitHash
}
