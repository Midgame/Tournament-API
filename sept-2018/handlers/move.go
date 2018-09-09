package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Move determines if a move is valid for a given callsign, and updates the bot's location if so (and returns the new location)
func Move(req structs.MoveRequest, knownBots map[string]structs.Bot, grid structs.Grid) structs.StatusResponse {
	resp := structs.StatusResponse{
		Nodes: []structs.NodeStatus{},
		Error: false,
	}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		resp.ErrorMsg = BOT_NOT_FOUND_ERROR
		return resp
	}

	newLocation := grid.MoveBot(bot, req.X, req.Y)
	bot.Location = newLocation
	knownBots[bot.Id] = bot

	// Error checking
	if bot.Location.X != req.X && bot.Location.Y != req.Y {
		resp.Error = true
		resp.ErrorMsg = INVALID_MOVE
	}

	resp.Status = bot.GetStatus()
	return resp
}
