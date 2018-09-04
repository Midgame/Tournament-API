package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/handlers"
)

var (
	s Server
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func TestRegister(t *testing.T) {
	s = Server{}
	s.Initialize()
	route := "/register"

	getResponse := func(postBody string) (handlers.RegisterResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		var response handlers.RegisterResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	// Simple test: No params
	simpleResponse, simpleResult := getResponse(`{}`)
	if simpleResponse.DebugMode || simpleResponse.Callsign == "" {
		t.Errorf("RegistrationHandler should return a valid result. Got: %v", simpleResult)
	}

	// One param: debug
	debugResponse, debugResult := getResponse(`{"debug": "true"}`)
	if !debugResponse.DebugMode {
		t.Errorf("RegistrationHandler should respect the debug flag. Got: %v", debugResult)
	}

	// One param: callsign
	callsignResponse, callsignResult := getResponse(`{"callsign": "foobar"}`)
	if callsignResponse.Callsign != "foobar" || callsignResponse.DebugMode {
		t.Errorf("RegistrationHandler should respect the callsign param. got %v", callsignResult)
	}

	// Both params: callsign and debug
	bothResponse, bothResult := getResponse(`{"callsign": "foobar", "debug": "true"}`)
	if !bothResponse.DebugMode || bothResponse.Callsign != "foobar" {
		t.Errorf("RegistrationHandler should respect both params. Debug: %v Callsign: %s Raw: %v", bothResponse.DebugMode, bothResponse.Callsign, bothResult)
	}
}

func TestInit(t *testing.T) {
	s = Server{}
	s.Initialize()

	if len(s.KnownBots) != 0 {
		t.Errorf("Known bots wasn't initialized to 0 properly")
	}
	if s.Grid.Height != GRID_HEIGHT {
		t.Errorf("Height not initialized properly for grid")
	}
	if s.Grid.Width != GRID_WIDTH {
		t.Errorf("Width not initialized properly for grid")
	}
	if len(s.Grid.Entities) != 0 {
		t.Errorf("Grid entities not initialized properly")
	}
}
