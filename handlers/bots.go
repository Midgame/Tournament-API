package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/structs"
)

func Bots(bots map[string]structs.Bot) map[string][]structs.BotStatus {

	response := []structs.BotStatus{}

	for _, bot := range bots {
		response = append(response, bot.GetStatus())
	}

	responseMap := make(map[string][]structs.BotStatus)
	responseMap["Bots"] = response
	return responseMap
}
