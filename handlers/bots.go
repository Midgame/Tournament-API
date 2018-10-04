package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/structs"
)

func Bots(bots map[string]structs.Bot) []structs.BotStatus {

	response := []structs.BotStatus{}

	for _, bot := range bots {
		response = append(response, bot.GetStatus())
	}

	return response
}
