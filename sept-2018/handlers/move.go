package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

type MoveRequest struct {
	Callsign string `json:"callsign"`
	X        int    `json:"x,string"`
	Y        int    `json:"y,string"`
}

type MoveResponse struct {
	Location structs.GridLocation
	Error    bool
}

// Move determines if a move is valid for a given callsign, and updates the bot's location if so (and returns the new location)
func Move(req MoveRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot, grid structs.Grid) MoveResponse {
	resp := MoveResponse{
		Error: false,
	}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		return resp
	}

	resp.Location = bot.Location

	newLocation := grid.MoveBot(bot, req.X, req.Y)
	resp.Location = newLocation
	bot.Location = newLocation
	return resp
}
