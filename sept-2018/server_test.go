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
	bot := fetchBot(uuid)
	if bot.Id != uuid {
		t.Errorf("Bot didn't register UUID properly")
	}
	if len(bot.Claims) != 0 {
		t.Errorf("Bot somehow started out with claims immediately after registration")
	}
}
