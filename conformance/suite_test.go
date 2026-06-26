package conformance

import "testing"

func TestSyntheticSuite(t *testing.T) {
	result := RunSuite(SyntheticSuite())
	if result.Failed != 0 {
		t.Fatalf("failed=%d cases=%+v", result.Failed, result.Cases)
	}
	if result.Passed != len(result.Cases) {
		t.Fatalf("passed=%d cases=%d", result.Passed, len(result.Cases))
	}
}
