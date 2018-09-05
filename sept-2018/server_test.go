package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/HeadlightLabs/Tournament-API/sept-2018"
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

var (
	s Server
)

func TestReleaseRequest(t *testing.T) {
	s = Server{}
	s.Initialize()
	route := "/release"

	s.Grid.Bots["alpha"] = createBot("alpha", []string{"delta"})
	s.Grid.Bots["beta"] = createBot("beta", []string{"gamma"})
	s.Grid.Nodes["gamma"] = createNode("gamma", "beta")
	s.Grid.Nodes["delta"] = createNode("delta", "alpha")
	s.Grid.Nodes["epsilon"] = createNode("epsilon", "alpha")

	getResponse := func(postBody string) (structs.ReleaseResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		var response structs.ReleaseResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	tt := []struct {
		payload  string
		response structs.ReleaseResponse
	}{
		{`{}`, structs.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"foobar", "node":"delta"}`, structs.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"foobar"}`, structs.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"gamma"}`, structs.ReleaseResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"delta"}`, structs.ReleaseResponse{Success: true, Error: false}},
		{`{"callsign":"alpha", "node":"epsilon"}`, structs.ReleaseResponse{Success: false, Error: true}},
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

	getResponse := func(postBody string) (structs.RegisterResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		var response structs.RegisterResponse
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
	s.Grid.Bots["alpha"] = createBot("alpha", []string{})
	s.Grid.Bots["beta"] = createBot("beta", []string{"gamma"})
	s.Grid.Nodes["gamma"] = createNode("gamma", "beta")
	s.Grid.Nodes["delta"] = createNode("delta", "")

	route := "/claim"

	getResponse := func(postBody string) (structs.ClaimResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		checkResponseCode(t, http.StatusOK, result.Code)

		var response structs.ClaimResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	tt := []struct {
		payload  string
		response structs.ClaimResponse
	}{
		{`{}`, structs.ClaimResponse{Success: false, Error: true}},
		{`{"callsign":"alpha"}`, structs.ClaimResponse{Success: false, Error: true}},
		{`{"node":"delta"}`, structs.ClaimResponse{Success: false, Error: true}},
		{`{"callsign":"alpha", "node":"delta"}`, structs.ClaimResponse{Success: true, Error: true}},
		{`{"callsign":"alpha", "node":"gamma"}`, structs.ClaimResponse{Success: false, Error: false}},
		{`{"callsign":"alpha", "node":"delta"}`, structs.ClaimResponse{Success: false, Error: false}},
		{`{"callsign":"alpha", "node":"foobar"}`, structs.ClaimResponse{Success: false, Error: true}},
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
	s.Grid.Bots["alpha"] = createBot("alpha", []string{})
	s.Grid.Bots["beta"] = createBot("beta", []string{"gamma"})
	s.Grid.Nodes["gamma"] = createNode("gamma", "beta")
	betaBot := s.Grid.Bots["beta"]
	betaBot.DebugMode = true
	s.Grid.Bots["beta"] = betaBot

	route := "/status"

	getResponse := func(postBody string) (structs.StatusResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		checkResponseCode(t, http.StatusOK, result.Code)

		var response structs.StatusResponse
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
	s.Grid.Bots["alpha"] = createBot("alpha", []string{})
	s.Grid.Bots["beta"] = createBot("beta", []string{"gamma"})
	s.Grid.Nodes["gamma"] = createNode("gamma", "beta")

	route := "/mine"

	getResponse := func(postBody string) (structs.MineResponse, string) {
		payload := []byte(postBody)
		request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
		result := executeRequest(request)
		checkResponseCode(t, http.StatusOK, result.Code)

		var response structs.MineResponse
		json.Unmarshal([]byte(result.Body.String()), &response)
		return response, result.Body.String()
	}

	tt := []struct {
		Payload   string
		Mined     int
		Remaining int
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

	// Test cases:
	// 1) node too far to the right
	// 2) node too far to the left
	// 3) node too far up
	// 4) node too far down
	// 5) node too far up but within left/right range
	// 6) node too far left but within up/down range
	// 9) node on left edge, before/after overlap (but within range)
	// 10) node on right edge, before/after overlap (but within range)
	// 11) node on top edge, before/after overlap (but within range)
	// 12) node on bottom edge, before/after overlap (within range)
	// 13) node on left/right/top/bottom edge, after overlap (not within range)
	// 14) node on right edge, just within scan range (exactly 5 units away)
}

func TestMoveRequest(t *testing.T) {
	t.Errorf("Not implemented yet")
}

func TestInit(t *testing.T) {
	s = Server{}
	s.Initialize()

	if len(s.Grid.Bots) != 0 {
		t.Errorf("Known bots wasn't initialized to 0 properly")
	}
	if s.Grid.Height == 0 {
		t.Errorf("Height not initialized properly for grid")
	}
	if s.Grid.Width == 0 {
		t.Errorf("Width not initialized properly for grid")
	}
	if len(s.Grid.Bots) != 0 {
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
