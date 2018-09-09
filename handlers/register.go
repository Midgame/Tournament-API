package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/google/uuid"
)

// registerUser generates a new UUID for a user, adds that UUID to the list of known bots,
// and then returns the bot entity.
func RegisterUser(req structs.SimpleRequest, grid structs.Grid) (structs.Bot, structs.StatusResponse) {
	// TODO: Error out if callsign already chosen

	if req.Callsign == "" {
		req.Callsign = uuid.New().String()
	}
	bot := grid.InitializeBot(req.Callsign)

	response := structs.StatusResponse{
		Status: bot.GetStatus(),
	}

	return bot, response
}
