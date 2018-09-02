package main

import (
	"encoding/json"
	"net/http"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/handlers"
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/gorilla/mux"
)

var (
	knownBots map[string]structs.Bot = make(map[string]structs.Bot)
)

// RegistrationHandler accepts registration from a new bot. It generates a UUID for the user, registers it,
// and returns the UUID to the user
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	bot, response := handlers.RegisterUser()
	knownBots[bot.Id] = bot

	json.NewEncoder(w).Encode(response)
}

func main() {
	// Setup
	r := mux.NewRouter()
	r.HandleFunc("/register", RegistrationHandler)

	// Run
	http.Handle("/", r)
}
