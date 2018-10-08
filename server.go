package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/HeadlightLabs/Tournament-API/handlers"
	"github.com/HeadlightLabs/Tournament-API/structs"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	Grid   structs.Grid
	HTTP   *http.Server
}

// Logs the response and the original parameters just for thoroughness
func LogResponse(response structs.StatusResponse, route string) {

	resp, err := json.Marshal(response)
	if err != nil {
		glog.Infof("[%s](%s) Error marshalling JSON: %v", route, response.Status.Id, err)
	} else {
		glog.Infof("[%s](%s) Response: %s", route, response.Status.Id, string(resp))
	}

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

func createSimpleRequest(w http.ResponseWriter, r *http.Request, route string) (structs.SimpleRequest, error) {
	var req structs.SimpleRequest
	body, _ := ioutil.ReadAll(r.Body)
	glog.Infof("[%s][RAW] Request params: %v", route, string(body))
	decoder := json.NewDecoder(bytes.NewBuffer(body))

	if err := decoder.Decode(&req); err != nil {
		errorMsg := fmt.Sprintf("[%s][ERR] Invalid request. Params: %v. Error msg: %v", route, string(body), err)
		glog.Info(errorMsg)
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return req, err
	}
	defer r.Body.Close()
	return req, nil
}

// releaseHandler releases a claim on a node.
func (s *Server) releaseHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r, "RELEASE")
	if err != nil {
		return
	}

	response := handlers.Release(req, s.Grid.Nodes, s.Grid.Bots)
	LogResponse(response, "RELEASE")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

// claimHandler accepts a claim from an existing callsign.
func (s *Server) claimHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r, "CLAIM")
	if err != nil {
		return
	}

	response := handlers.Claim(req, s.Grid.Nodes, s.Grid.Bots, s.Grid)
	LogResponse(response, "CLAIM")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

// registrationHandler accepts registration from a new bot.
func (s *Server) registrationHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r, "REGISTER")
	if err != nil {
		return
	}

	bot, response := handlers.RegisterUser(req, s.Grid)
	s.Grid.Bots[bot.Id] = bot
	LogResponse(response, "REGISTER")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

// mineHandler accepts a mining request from a given callsign and node id.
func (s *Server) mineHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r, "MINE")
	if err != nil {
		return
	}

	response := handlers.Mine(req, s.Grid.Nodes, s.Grid.Bots, s.Grid)
	LogResponse(response, "MINE")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

// scanHandler accepts a scan request and returns information around a given callsign
func (s *Server) scanHandler(w http.ResponseWriter, r *http.Request) {
	req, err := createSimpleRequest(w, r, "SCAN")
	if err != nil {
		return
	}

	response := handlers.Scan(req, s.Grid.Nodes, s.Grid.Bots, s.Grid)
	LogResponse(response, "SCAN")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

// moveHandler accepts a move request from a given callsign and moves it to the requested location
func (s *Server) moveHandler(w http.ResponseWriter, r *http.Request) {
	var req structs.MoveRequest

	body, _ := ioutil.ReadAll(r.Body)
	glog.Infof("[MOVE][RAW] Request params: %v", string(body))

	decoder := json.NewDecoder(bytes.NewBuffer(body))
	if err := decoder.Decode(&req); err != nil {
		errorMsg := fmt.Sprintf("[MOVE][ERR] Invalid request. Params: %v. Error msg: %v", string(body), err)
		glog.Info(errorMsg)
		respondWithError(w, http.StatusBadRequest, errorMsg)
		return
	}
	defer r.Body.Close()

	response := handlers.Move(req, s.Grid.Bots, s.Grid)
	LogResponse(response, "MOVE")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

// Redirects to the documentation page, so anyone hitting this page through a browser gets more information
func (s *Server) redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://docs.headlightlabs.com", 302)
}

func (s *Server) botsHandler(w http.ResponseWriter, r *http.Request) {
	response := handlers.Bots(s.Grid.Bots)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

func (s *Server) nodesHandler(w http.ResponseWriter, r *http.Request) {
	response := handlers.Nodes(s.Grid.Nodes)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(response)
	glog.Flush()
}

// Initializes the server with some defaults
func (s *Server) Initialize() {
	s.Grid = structs.Grid{}
	s.Grid.Initialize()
	s.Router = mux.NewRouter()
	s.initializeRoutes()
	port := os.Getenv("PORT")
	s.HTTP = &http.Server{
		Addr:    ":" + port,
		Handler: s.Router,
	}
}

// Initializes all the routes
func (s *Server) initializeRoutes() {
	// Original tournament routes
	s.Router.HandleFunc("/register", s.registrationHandler).Methods("POST")
	s.Router.HandleFunc("/claim", s.claimHandler).Methods("POST")
	s.Router.HandleFunc("/release", s.releaseHandler).Methods("POST")
	s.Router.HandleFunc("/mine", s.mineHandler).Methods("POST")
	s.Router.HandleFunc("/scan", s.scanHandler).Methods("POST")
	s.Router.HandleFunc("/move", s.moveHandler).Methods("POST")
	s.Router.HandleFunc("/", s.redirectHandler).Methods("GET")

	// New tournament routes
	s.Router.HandleFunc("/bots", s.botsHandler).Methods("GET", "OPTIONS")
	s.Router.HandleFunc("/nodes", s.nodesHandler).Methods("GET", "OPTIONS")
}

func (s *Server) Run() {
	port := os.Getenv("PORT")
	glog.Info("Starting on port ", port)
	glog.Flush()
	go func() {
		if err := s.HTTP.ListenAndServe(); err != nil {
			glog.Infof("Server error: %s", err)
		}
	}()
}

func (s *Server) Shutdown() error {
	err := s.HTTP.Shutdown(nil)
	s.HTTP = nil
	s.Router = nil
	return err
}
