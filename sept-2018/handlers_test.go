package main

import (
	"testing"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/handlers"
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
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
	if bot.Location.X != 0 || bot.Location.Y != 0 {
		t.Errorf("Bot wasn't initialized with a location properly")
	}
}

func createBot(uuid string) structs.Bot {
	bot := structs.Bot{
		GridEntity: structs.GridEntity{
			Id:   uuid,
			Type: structs.BOT,
			Location: structs.GridLocation{
				X: 0,
				Y: 0,
			},
		},
		Claims: []string{},
	}
	return bot
}

func TestStatus(t *testing.T) {
	knownBots := make(map[string]structs.Bot)
	knownBots["alpha"] = createBot("alpha")
	knownBots["beta"] = createBot("beta")
	knownBots["gamma"] = createBot("gamma")

	validResult := handlers.Status("beta", false, knownBots)
	if len(validResult.Bots) != 1 || validResult.Bots[0].Id != "beta" {
		t.Errorf("Non-debug result should find single bot with valid uuid. Bot found has ID: %s", validResult.Bots[0].Id)
	}

	invalidResult := handlers.Status("delta", false, knownBots)
	if len(invalidResult.Bots) > 0 {
		t.Errorf("Non-debug result should find no bots with invalid uuid, found: %d", len(invalidResult.Bots))
	}

	debugResult := handlers.Status("delta", true, knownBots)
	if len(debugResult.Bots) != 3 {
		t.Errorf("Debug result should return all bots, found: %d", len(debugResult.Bots))
	}

}
