package main

import (
	"encoding/json"
	"fmt"
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

// Writes out to the provided response writer with a JSON response. Only when response is successful.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Writes out to the provided response writer with an error code.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// releaseHandler releases a claim on a node.
func (s *Server) releaseHandler(w http.ResponseWriter, r *http.Request) {
	var req handlers.ReleaseRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request payload: %v", r.Body)
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return
	}
	defer r.Body.Close()

	response := handlers.Release(req, s.KnownNodes, s.KnownBots)
	json.NewEncoder(w).Encode(response)
}

// claimHandler accepts a claim from an existing callsign.
func (s *Server) claimHandler(w http.ResponseWriter, r *http.Request) {
	var req handlers.ClaimRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request payload: %v", r.Body)
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return
	}
	defer r.Body.Close()

	response := handlers.Claim(req, s.KnownNodes, s.KnownBots)
	json.NewEncoder(w).Encode(response)
}

// registrationHandler accepts registration from a new bot.
func (s *Server) registrationHandler(w http.ResponseWriter, r *http.Request) {
	var req handlers.RegisterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request payload: %v", r.Body)
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return
	}
	defer r.Body.Close()

	bot, response := handlers.RegisterUser(req)
	s.KnownBots[bot.Id] = bot

	json.NewEncoder(w).Encode(response)
}

// Initializes the server with some defaults
func (s *Server) Initialize() {
	s.KnownBots = make(map[string]structs.Bot)
	s.KnownNodes = make(map[string]structs.Node)
	s.Grid = structs.Grid{Width: GRID_WIDTH, Height: GRID_HEIGHT, Entities: [][]structs.GridEntity{}}
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

// Initializes all the routes
func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/register", s.registrationHandler).Methods("POST")
	s.Router.HandleFunc("/claim", s.claimHandler).Methods("POST")
	s.Router.HandleFunc("/release", s.releaseHandler).Methods("POST")
}

func (s *Server) Run() {
	http.Handle("/", s.Router)
}
