package deadball

const (
	verboseWarning byte = iota
	verboseInfo
	verboseDebug
)

var (
	// verbosity determines how granular logs are
	verbosity byte = verboseInfo
)

// Verbosity returns the currently configured verbosity
func Verbosity() byte {
	return verbosity
}

// SetVerbosity changes the verbosity level
func SetVerbosity(newVerbosity byte) {
	verbosity = newVerbosity
}
