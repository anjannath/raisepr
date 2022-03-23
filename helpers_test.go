package main

import (
	"testing"
	"time"
)

func TestRandomBranchName(t *testing.T) {
	b1 := randBranchName()
	time.Sleep(time.Second)
	b2 := randBranchName()

	if b1 == b2 {
		t.Errorf("Generated same branch name twice b1=%s and b2=%s", b1, b2)
	}
}

func TestMissingOptsCausesError(t *testing.T) {
	cfg, err := getOptsFromEnv()
	if err == nil {
		t.Error("Missing options are not triggering error")
	}

	cnf := config{}
	if cfg != cnf {
		t.Errorf("Wrong config struct returned, should be empty: %v", cfg)
	}
}
