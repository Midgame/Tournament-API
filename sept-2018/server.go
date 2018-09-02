package main

import (
	"encoding/json"
	"net/http"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	knownBots map[string]structs.Bot = make(map[string]structs.Bot)
)

type Status struct {
	Uuid string
}

// knownBotCount returns a count of the number of known bots
func knownBotCount() int {
	return len(knownBots)
}

// fetchBot is an internal helper method for fetching a bot out of the map structure
func fetchBot(uuid string) structs.Bot {
	return knownBots[uuid]
}

// registerUser generates a new UUID for a user, adds that UUID to the list of known bots,
// and then returns the UUID.
func registerUser() string {
	uuid := uuid.New().String()
	knownBots[uuid] = structs.Bot{
		GridEntity: structs.GridEntity{
			Id:   uuid,
			Type: structs.BOT,
		},
		Claims: []string{},
	}

	return uuid
}

// RegistrationHandler accepts registration from a new bot. It generates a UUID for the user, registers it,
// and returns the UUID to the user
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	uuid := registerUser()

	status := Status{Uuid: uuid}
	json.NewEncoder(w).Encode(status)
}

func main() {
	// Setup
	r := mux.NewRouter()
	r.HandleFunc("/register", RegistrationHandler)

	// Run
	http.Handle("/", r)
}
