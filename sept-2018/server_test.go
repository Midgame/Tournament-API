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

	// Simple test: No params. Should fail.
	simpleResponse, simpleResult := getResponse(`{}`)
	if simpleResponse.Success {
		t.Errorf("ClaimHandler shouldn't succeed with no params. Got: %v", simpleResult)
	}

	// One param: Callsign. Should fail.
	reqResponse, reqResult := getResponse(`{"callsign":"alpha"}`)
	if reqResponse.Success {
		t.Errorf("ClaimHandler shouldn't succeed with one param. Got: %v", reqResult)
	}

	// One param: Node. Should fail.
	nodeResponse, nodeResult := getResponse(`{"node":"delta"}`)
	if nodeResponse.Success {
		t.Errorf("ClaimHandler shouldn't succeed with one param. Got: %v", nodeResult)
	}

	// Valid claim. Should succeed.
	validResponse, validResult := getResponse(`{"callsign":"alpha", "node":"delta"}`)
	if !validResponse.Success {
		t.Errorf("ClaimHandler should succeed with valid claim. Got: %v", validResult)
	}

	// Invalid claim. Should fail.
	invalidResponse, invalidResult := getResponse(`{"callsign":"alpha", "node":"gamma"}`)
	if invalidResponse.Success {
		t.Errorf("ClaimHandler should not succeed with invalid claim. Got: %v", invalidResult)
	}

	// Already claimed by callsign. Should succeed.
	alreadyResponse, alreadyResult := getResponse(`{"callsign":"alpha", "node":"delta"}`)
	if !alreadyResponse.Success {
		t.Errorf("ClaimHandler should succeed with already claimed node. Got: %v", alreadyResult)
	}

	// Invalid node. Should return error.
	invalidNodeResponse, invalidNodeResult := getResponse(`{"callsign":"alpha", "node":"foobar"}`)
	if invalidNodeResponse.Success || !invalidNodeResponse.Error {
		t.Errorf("ClaimHandler should error with non-existant claimed node. Got: %v", invalidNodeResult)
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
