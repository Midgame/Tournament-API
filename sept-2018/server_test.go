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

type testCase struct {
	payload          string
	expectedCode     int
	expectedResponse structs.StatusResponse
}

func TestReleaseRequest(t *testing.T) {
	route := "/release"

	tt := []testCase{
		{`{}`, 200, createStatusResponse(true)}, // No params
	}

	testRoute(route, tt, t)
}

func TestRegisterRequest(t *testing.T) {
	route := "/register"

	tt := []testCase{
		{`{}`, 200, createStatusResponse(false)},                    // No params
		{`{"callsign":"foobar"}`, 200, createStatusResponse(false)}, // Valid params
	}

	testRoute(route, tt, t)
}

func TestClaimRequest(t *testing.T) {
	route := "/claim"

	tt := []testCase{
		{`{}`, 200, createStatusResponse(true)}, // No params
	}

	testRoute(route, tt, t)
}

func TestStatusRequest(t *testing.T) {
	route := "/status"

	tt := []testCase{
		{`{}`, 200, createStatusResponse(true)},                    // No params
		{`{"callsign":"alpha"}`, 200, createStatusResponse(false)}, // Valid params
		{`{"callsign":"delta"}`, 200, createStatusResponse(true)},  // Invalid bot
	}

	testRoute(route, tt, t)
}

func TestMineRequest(t *testing.T) {
	route := "/mine"

	tt := []testCase{
		{`{}`, 200, createStatusResponse(true)},                                    // No params
		{`{"callsign":"beta", "node": "gamma"}`, 200, createStatusResponse(false)}, // Valid params
	}

	testRoute(route, tt, t)
}

func TestScanRequest(t *testing.T) {
	route := "/scan"

	tt := []testCase{
		{`{}`, 200, createStatusResponse(true)},
	}
	testRoute(route, tt, t)
}

func TestMoveRequest(t *testing.T) {
	route := "/move"

	tt := []testCase{
		{`{}`, 200, createStatusResponse(true)},
	}
	testRoute(route, tt, t)
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

func testRoute(route string, testCases []testCase, t *testing.T) {
	s = Server{}
	s.Initialize()

	s.Grid.Bots["alpha"] = createBot("alpha", []string{"delta"})
	s.Grid.Bots["beta"] = createBot("beta", []string{"gamma"})
	s.Grid.Nodes["gamma"] = createNode("gamma", "beta")
	s.Grid.Nodes["delta"] = createNode("delta", "delta")
	s.Grid.Nodes["epsilon"] = createNode("epsilon", "")

	// Sending unparseable garbage should result in a 400 error
	testCases = append(testCases, testCase{`{"callsign: ""test"}`, 400, structs.StatusResponse{}})

	for _, tc := range testCases {
		response, result := getResponse(tc.payload, route, tc.expectedCode, t)
		if response.Error != tc.expectedResponse.Error {
			t.Errorf("\n(%s) returned bad result.\nOriginal postBody: %v\nExpected: %v\nActual: %v\nRaw result: %v\n\n", route, tc.payload, tc.expectedResponse, response, result)
		}
	}
}

/** Helper functions */

func createStatusResponse(errorExp bool) structs.StatusResponse {
	return structs.StatusResponse{
		Error: errorExp,
	}
}

func checkResponseCode(t *testing.T, expected, actual int, route string, postBody string) {
	if expected != actual {
		t.Errorf("\n(%s) Expected response code %d. Got %d. Post body: %v\n", route, expected, actual, postBody)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func getResponse(postBody string, route string, expectedCode int, t *testing.T) (structs.StatusResponse, string) {
	payload := []byte(postBody)
	request, _ := http.NewRequest("POST", route, bytes.NewBuffer(payload))
	result := executeRequest(request)
	checkResponseCode(t, expectedCode, result.Code, route, postBody)

	var response structs.StatusResponse
	json.Unmarshal([]byte(result.Body.String()), &response)
	return response, result.Body.String()
}

func createBot(uuid string, claims []string) structs.Bot {
	bot := structs.Bot{
		GridEntity: structs.GridEntity{
			Id:   uuid,
			Type: structs.BOT,
			Location: structs.GridLocation{
				X: 0,
				Y: 0,
			},
		},
		Claims: claims,
	}
	return bot
}

func createNode(uuid string, claimedBy string) structs.Node {
	node := structs.Node{
		GridEntity: structs.GridEntity{
			Id:   uuid,
			Type: structs.NODE,
			Location: structs.GridLocation{
				X: 0,
				Y: 0,
			},
		},
		ClaimedBy: claimedBy,
		Value:     1,
	}
	return node
}
