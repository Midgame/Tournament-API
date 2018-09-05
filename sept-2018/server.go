package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/handlers"
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Grid   structs.Grid
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

func createSimpleRequest(w http.ResponseWriter, r *http.Request) (structs.SimpleRequest, error) {
	var req structs.SimpleRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request payload: %v", r.Body)
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return req, err
	}
	defer r.Body.Close()
	return req, nil
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r)
	if err != nil {
		return
	}

	response := handlers.Status(req, s.Grid.Bots)
	json.NewEncoder(w).Encode(response)
}

// releaseHandler releases a claim on a node.
func (s *Server) releaseHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r)
	if err != nil {
		return
	}

	response := handlers.Release(req, s.Grid.Nodes, s.Grid.Bots)
	json.NewEncoder(w).Encode(response)
}

// claimHandler accepts a claim from an existing callsign.
func (s *Server) claimHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r)
	if err != nil {
		return
	}

	response := handlers.Claim(req, s.Grid.Nodes, s.Grid.Bots)
	json.NewEncoder(w).Encode(response)
}

// registrationHandler accepts registration from a new bot.
func (s *Server) registrationHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r)
	if err != nil {
		return
	}

	bot, response := handlers.RegisterUser(req)
	s.Grid.Bots[bot.Id] = bot

	json.NewEncoder(w).Encode(response)
}

// mineHandler accepts a mining request from a given callsign and node id.
func (s *Server) mineHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r)
	if err != nil {
		return
	}

	response := handlers.Mine(req, s.Grid.Nodes, s.Grid.Bots)

	json.NewEncoder(w).Encode(response)
}

// scanHandler accepts a scan request and returns information around a given callsign
func (s *Server) scanHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r)
	if err != nil {
		return
	}

	response := handlers.Scan(req, s.Grid.Nodes, s.Grid.Bots, s.Grid)

	json.NewEncoder(w).Encode(response)
}

// moveHandler accepts a move request from a given callsign and moves it to the requested location
func (s *Server) moveHandler(w http.ResponseWriter, r *http.Request) {
	var req structs.MoveRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		errorMsg := fmt.Sprintf("Invalid request payload: %v", r.Body)
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return
	}
	defer r.Body.Close()

	response := handlers.Move(req, s.Grid.Nodes, s.Grid.Bots, s.Grid)

	json.NewEncoder(w).Encode(response)
}

// Initializes the server with some defaults
func (s *Server) Initialize() {
	s.Grid = structs.Grid{}
	s.Grid.Initialize()
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

// Initializes all the routes
func (s *Server) initializeRoutes() {
	s.Router.HandleFunc("/register", s.registrationHandler).Methods("POST")
	s.Router.HandleFunc("/claim", s.claimHandler).Methods("POST")
	s.Router.HandleFunc("/release", s.releaseHandler).Methods("POST")
	s.Router.HandleFunc("/status", s.statusHandler).Methods("POST")
	s.Router.HandleFunc("/mine", s.mineHandler).Methods("POST")
	s.Router.HandleFunc("/scan", s.scanHandler).Methods("POST")
	s.Router.HandleFunc("/move", s.moveHandler).Methods("POST")
}

func (s *Server) Run() {
	http.Handle("/", s.Router)
}
