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

func TestReleaseRequest(t *testing.T) {
	s = Server{}
	s.Initialize()
	route := "/release"

	s.KnownBots["alpha"] = createBot("alpha", []string{"delta"})
	s.KnownBots["beta"] = createBot("beta", []string{"gamma"})
	s.KnownNodes["gamma"] = createNode("gamma", "beta")
	s.KnownNodes["delta"] = createNode("delta", "alpha")
	s.KnownNodes["epsilon"] = createNode("epsilon", "alpha")

	getResponse := func(postBody string) (handlers.ReleaseResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		var response handlers.ReleaseResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	tt := []struct {
		payload  string
		response handlers.ReleaseResponse
	}{
		{`{}`, handlers.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"foobar", "node":"delta"}`, handlers.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"foobar"}`, handlers.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"gamma"}`, handlers.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"delta"}`, handlers.ReleaseResponse{Success: true, Error: false}},
		{`{"callsign":"alpha", "node":"epsilon"}`, handlers.ReleaseResponse{Success: false, Error: true}},
	}

	for _, tc := range tt {
		response, result := getResponse(tc.payload)
		if response != tc.response {
			t.Errorf("ReleaseHandler returned bad result. Expected: %v Actual: %v Raw: %v", tc.response, response, result)
		}
	}
}

func TestRegisterRequest(t *testing.T) {
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

func TestClaimRequest(t *testing.T) {
	s = Server{}
	s.Initialize()
	s.KnownBots["alpha"] = createBot("alpha", []string{})
	s.KnownBots["beta"] = createBot("beta", []string{"gamma"})
	s.KnownNodes["gamma"] = createNode("gamma", "beta")
	s.KnownNodes["delta"] = createNode("delta", "")

	route := "/claim"

	getResponse := func(postBody string) (handlers.ClaimResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		checkResponseCode(t, http.StatusOK, result.Code)

		var response handlers.ClaimResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	tt := []struct {
		payload  string
		response handlers.ClaimResponse
	}{
		{`{}`, handlers.ClaimResponse{Success: false, Error: true}},
		{`{"callsign":"alpha"}`, handlers.ClaimResponse{Success: false, Error: true}},
		{`{"node":"delta"}`, handlers.ClaimResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"delta"}`, handlers.ClaimResponse{Success: true, Error: true}},
		{`{"callsign":"alpha", "node":"gamma"}`, handlers.ClaimResponse{Success: false, Error: false}},
		{`{"callsign":"alpha", "node":"delta"}`, handlers.ClaimResponse{Success: false, Error: false}},
		{`{"callsign":"alpha", "node":"foobar"}`, handlers.ClaimResponse{Success: false, Error: true}},
	}

	for _, tc := range tt {
		response, result := getResponse(tc.payload)
		if response.Success != tc.response.Success && response.Error != tc.response.Error {
			t.Errorf("ClaimHandler returned bad result. Expected: %v Actual: %v Raw: %v", tc.response, response, result)
		}
	}

}

func TestStatusRequest(t *testing.T) {
	s = Server{}
	s.Initialize()
	s.KnownBots["alpha"] = createBot("alpha", []string{})
	s.KnownBots["beta"] = createBot("beta", []string{"gamma"})
	s.KnownNodes["gamma"] = createNode("gamma", "beta")
	betaBot := s.KnownBots["beta"]
	betaBot.DebugMode = true
	s.KnownBots["beta"] = betaBot

	route := "/status"

	getResponse := func(postBody string) (handlers.StatusResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		checkResponseCode(t, http.StatusOK, result.Code)

		var response handlers.StatusResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	tt := []struct {
		Payload string
		Length  int
		Error   bool
	}{
		{`{}`, 0, true},
		{`{"callsign":"alpha"}`, 1, false},
		{`{"callsign":"beta"}`, 2, false},
		{`{"callsign":"delta"}`, 0, true},
	}

	for _, tc := range tt {
		response, result := getResponse(tc.Payload)
		if len(response.Bots) != tc.Length && response.Error != tc.Error {
			t.Errorf("ClaimHandler returned bad result. Raw: %v", result)
		}
	}
}

func TestMineRequest(t *testing.T) {
	s = Server{}
	s.Initialize()
	s.KnownBots["alpha"] = createBot("alpha", []string{})
	s.KnownBots["beta"] = createBot("beta", []string{"gamma"})
	s.KnownNodes["gamma"] = createNode("gamma", "beta")

	route := "/mine"

	getResponse := func(postBody string) (handlers.MineResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		checkResponseCode(t, http.StatusOK, result.Code)

		var response handlers.MineResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	tt := []struct {
		Payload   string
		Mined     uint64
		Remaining uint64
		Error     bool
	}{
		{`{}`, 0, 0, true},
		{`{"callsign":"alpha", "node":"foobar"}`, 0, 0, true},
		{`{"callsign":"foobar", "node":"gamma"}`, 0, 0, true},
		{`{"callsign":"alpha", "node":"gamma"}`, 0, 0, true},
		{`{"callsign":"beta", "node":"gamma"}`, 1, 0, false},
		{`{"callsign":"beta", "node":"gamma"}`, 0, 0, true},
	}

	for _, tc := range tt {
		response, result := getResponse(tc.Payload)
		if tc.Error != response.Error || tc.Mined != response.AmountMined || tc.Remaining != response.AmountRemaining {
			t.Errorf("MineHandler returned bad result. Raw: %v", result)
		}
	}
}

func TestScanRequest(t *testing.T) {
	t.Errorf("Not implemented yet")
}

func TestMoveRequest(t *testing.T) {
	t.Errorf("Not implemented yet")
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

/** Helper functions */

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}
