package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Status returns information about the requesting user's:
// Location, Claims, Total score
// If in debug mode, also returns this information for all other known bots
func Status(req structs.SimpleRequest, knownBots map[string]structs.Bot) structs.StatusResponse {
	botList := []structs.BotStatus{}
	resp := structs.StatusResponse{
		Bots:  botList,
		Error: false,
	}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		return resp
	}

	// In debug mode, return all bots
	if bot.DebugMode {
		for _, bot := range knownBots {
			resp.Bots = append(resp.Bots, bot.GetStatus())
		}
		return resp
	}

	// In regular mode, just return the requested id
	resp.Bots = append(resp.Bots, bot.GetStatus())
	return resp

}
