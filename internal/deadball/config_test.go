package deadball

import "testing"

func TestSetVerbosity(t *testing.T) {
	SetVerbosity(verboseDebug)
	if expected, actual := verboseDebug, Verbosity(); expected != actual {
		t.Fatalf("Fail: expected verbosity level %d, got %d", expected, actual)
	}

	t.Log("Pass")
}
