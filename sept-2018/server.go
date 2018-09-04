package main

import (
	"encoding/json"
	"net/http"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/handlers"
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/gorilla/mux"
)

const (
	GRID_WIDTH  = 100
	GRID_HEIGHT = 100
)

type Server struct {
	Router     *mux.Router
	KnownBots  map[string]structs.Bot
	KnownNodes map[string]structs.Node
	Grid       structs.Grid
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// RegistrationHandler accepts registration from a new bot. It generates a UUID for the user, registers it,
// and returns the UUID to the user
func (s *Server) registrationHandler(w http.ResponseWriter, r *http.Request) {

	var req handlers.RegisterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	bot, response := handlers.RegisterUser(req)
	s.KnownBots[bot.Id] = bot

	json.NewEncoder(w).Encode(response)
}

func (s *Server) Initialize() {
	s.KnownBots = make(map[string]structs.Bot)
	s.KnownNodes = make(map[string]structs.Node)
	s.Grid = structs.Grid{Width: GRID_WIDTH, Height: GRID_HEIGHT, Entities: [][]structs.GridEntity{}}
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/register", s.registrationHandler).Methods("POST")
}

func (s *Server) Run() {
	http.Handle("/", s.Router)
}
