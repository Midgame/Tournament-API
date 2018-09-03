package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

type StatusResponse struct {
	Bots []structs.BotStatus
}

// Status returns information about the requesting user's:
// Location, Claims, Last 5 actions, Total score
// If in debug mode, also returns this information for all other known bots
func Status(id string, debug bool, knownBots map[string]structs.Bot) StatusResponse {
	botList := []structs.BotStatus{}
	response := StatusResponse{
		Bots: botList,
	}

	// In debug mode, return all bots
	if debug {
		for _, bot := range knownBots {
			response.Bots = append(response.Bots, bot.GetStatus())
		}
		return response
	}

	// In regular mode, just return the requested id
	if bot, ok := knownBots[id]; ok {
		response.Bots = append(response.Bots, bot.GetStatus())
	}
	return response

}
