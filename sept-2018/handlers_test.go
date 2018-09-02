package main

import (
	"testing"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/handlers"
)

func TestRegisterUser(t *testing.T) {
	bot, response := handlers.RegisterUser()
	if response.Uuid != bot.Id {
		t.Errorf("Response didn't return the correct UUID")
	}
	if bot.Id == "" {
		t.Errorf("Bot wasn't created with UUID properly")
	}
	if len(bot.Claims) != 0 {
		t.Errorf("Bot somehow started out with claims immediately after registration")
	}
}
