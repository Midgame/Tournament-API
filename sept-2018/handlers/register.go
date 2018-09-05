package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/google/uuid"
)

// registerUser generates a new UUID for a user, adds that UUID to the list of known bots,
// and then returns the bot entity.
func RegisterUser(req structs.SimpleRequest) (structs.Bot, structs.RegisterResponse) {
	if req.Callsign == "" {
		req.Callsign = uuid.New().String()
	}
	bot := structs.Bot{
		GridEntity: structs.GridEntity{
			Id:   req.Callsign,
			Type: structs.BOT,
			Location: structs.GridLocation{
				X: 0,
				Y: 0,
			},
		},
		DebugMode: req.DebugMode,
		Claims:    []string{},
	}

	response := structs.RegisterResponse{
		Callsign:  bot.Id,
		DebugMode: bot.DebugMode,
	}

	return bot, response
}
