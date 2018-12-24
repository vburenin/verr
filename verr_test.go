package verr

import "testing"

func TestError(t *testing.T) {
	err := Error("my error")
	if err.Code != 0 {
		t.Errorf("Incorrect error code, should be 0 but was %d", err.Code)
	}
}
