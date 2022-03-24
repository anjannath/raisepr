package helper

import (
	"testing"
	"time"
)

func TestRandomBranchName(t *testing.T) {
	t.Parallel()

	b1 := RandBranchName()
	time.Sleep(time.Second)
	b2 := RandBranchName()

	if b1 == b2 {
		t.Errorf("Generated same branch name twice b1=%s and b2=%s", b1, b2)
	}
}

func TestMissingOptsCausesError(t *testing.T) {
	t.Parallel()

	cfg, err := GetOptsFromEnv()
	if err == nil {
		t.Error("Missing options are not triggering error")
	}

	cnf := Config{}
	if cfg != cnf {
		t.Errorf("Wrong config struct returned, should be empty: %v", cfg)
	}
}
