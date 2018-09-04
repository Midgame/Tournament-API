package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	DebugMode bool   `json:"debug,string,omitempty"`
	Callsign  string `json:"callsign,omitempty"`
}

type RegisterResponse struct {
	Callsign  string
	DebugMode bool
}

// registerUser generates a new UUID for a user, adds that UUID to the list of known bots,
// and then returns the bot entity.
func RegisterUser(req RegisterRequest) (structs.Bot, RegisterResponse) {
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

	response := RegisterResponse{
		Callsign:  bot.Id,
		DebugMode: bot.DebugMode,
	}

	return bot, response
}
