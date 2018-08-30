package main

import "testing"

func TestRegisterUser(t *testing.T) {
	uuid := registerUser()
	if uuid == "" {
		t.Errorf("registerUser function didn't return a uuid")
	}
	if knownBotCount() != 1 {
		t.Errorf("registerUser function didn't register the uuid in the knownBot map")
	}
}
