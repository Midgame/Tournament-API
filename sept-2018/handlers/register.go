package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/google/uuid"
)

type Response struct {
	Uuid string
}

// registerUser generates a new UUID for a user, adds that UUID to the list of known bots,
// and then returns the bot entity.
func RegisterUser() (structs.Bot, Response) {
	uuid := uuid.New().String()
	bot := structs.Bot{
		GridEntity: structs.GridEntity{
			Id:   uuid,
			Type: structs.BOT,
		},
		Claims: []string{},
	}

	response := Response{
		Uuid: uuid,
	}

	return bot, response
}
