package handlers_test

import (
	"testing"

	"github.com/HeadlightLabs/Tournament-API/handlers"
	"github.com/HeadlightLabs/Tournament-API/structs"
)

func TestMine(t *testing.T) {

	knownBots := make(map[string]structs.Bot)
	knownBots["alpha"] = createBot("alpha", []string{"gamma"})
	knownBots["beta"] = createBot("beta", []string{"delta"})

	knownNodes := make(map[string]structs.Node)
	knownNodes["gamma"] = createNode("gamma", "alpha")
	knownNodes["delta"] = createNode("delta", "beta")
	knownNodes["epsilon"] = createNode("epsilon", "")

	// Case 4: Node unowned
	// Case 5: Node has no value left
	// Case 6: Callsign owns node
	makeReq := func(callsign string, node string) structs.SimpleRequest {
		return structs.SimpleRequest{
			Callsign: callsign,
			NodeId:   node,
		}
	}

	tt := []struct {
		callsign     string
		node         string
		expRemaining int
		expScore     int
		errorExp     bool
	}{
		{"alpha", "iota", 0, 0, true},    // Node does not exist
		{"omega", "gamma", 0, 0, true},   // Bot does not exist
		{"alpha", "delta", 0, 0, true},   // Node owned by someone else
		{"alpha", "epsilon", 0, 0, true}, // Node unowned
		{"alpha", "gamma", 0, 1, false},  // Successful mine
		{"alpha", "gamma", 0, 1, false},  // Successful mine but node tapped out
	}

	for _, tc := range tt {
		req := makeReq(tc.callsign, tc.node)
		actual := handlers.Mine(req, knownNodes, knownBots)
		if tc.errorExp {
			if !actual.Error {
				t.Errorf("[Mine] Error expected but not given. Actual result: %v", actual)
			}
			continue
		}

		remaining := actual.Nodes[0].Value
		score := actual.Status.Score
		if remaining != tc.expRemaining || score != tc.expScore {
			t.Errorf("[Mine] Actual and expected different. Actual (remaining, score): (%d,%d). Expected: (%d,%d)", remaining, score, tc.expRemaining, tc.expScore)
		}
	}
}

func TestMove(t *testing.T) {
	grid := structs.Grid{}
	grid.Initialize()

	knownBots := make(map[string]structs.Bot)
	knownBots["alpha"] = createBot("alpha", []string{})

	makeReq := func(callsign string, x int, y int) structs.MoveRequest {
		return structs.MoveRequest{
			Callsign: callsign,
			X:        x,
			Y:        y,
		}
	}
	makeExpResp := func(callsign string, x int, y int, errorExp bool) structs.GridLocation {
		return structs.GridLocation{
			X: x,
			Y: y,
		}
	}

	tt := []struct {
		x        int
		y        int
		noop     bool
		errorExp bool
	}{
		{1, 0, false, false},
		{1, 1, false, false},
		{2, 2, false, false},
		{1, 0, false, true},
		{100, 30, true, false},
	}

	for _, tc := range tt {
		callSign := "alpha"
		if tc.errorExp {
			callSign = "gamma"
		}

		originalX, originalY := knownBots[callSign].Location.X, knownBots[callSign].Location.Y

		req := makeReq(callSign, tc.x, tc.y)
		actual := handlers.Move(req, knownBots, grid)

		if tc.errorExp {
			if !actual.Error {
				t.Errorf("[Move] Error expected but not given for case: %v", tc)
			}
			continue
		}

		var expected structs.GridLocation
		if tc.noop {
			expected = makeExpResp(callSign, originalX, originalY, false)
		} else {
			expected = makeExpResp(callSign, tc.x, tc.y, false)
		}

		if actual.Status.Location != expected {
			t.Errorf("Move function didn't return expected result: %v. Actual: %v", expected, actual.Status.Location)
		}
	}
}

func TestRegisterUser(t *testing.T) {
	grid := structs.Grid{}
	grid.Initialize()

	req := structs.SimpleRequest{
		Callsign: "",
	}
	bot, response := handlers.RegisterUser(req, grid)
	if response.Status.Id != bot.Id {
		t.Errorf("Response didn't return the correct UUID")
	}
	if bot.Id == "" {
		t.Errorf("Bot wasn't created with UUID properly")
	}
	if len(bot.Claims) != 0 {
		t.Errorf("Bot somehow started out with claims immediately after registration")
	}
	if bot.Location.X == 0 || bot.Location.Y == 0 {
		t.Errorf("Bot wasn't initialized with a location properly")
	}

	debugReq := structs.SimpleRequest{
		Callsign: "foobar",
	}
	debugBot, _ := handlers.RegisterUser(debugReq, grid)
	if debugBot.Id != "foobar" {
		t.Errorf("Register function didn't accept callsign. Assigned: %s", bot.Id)
	}
}

func TestRelease(t *testing.T) {
	knownBots := make(map[string]structs.Bot)
	knownBots["alpha"] = createBot("alpha", []string{"gamma", "epsilon"})
	knownBots["beta"] = createBot("beta", []string{"delta"})

	knownNodes := make(map[string]structs.Node)
	knownNodes["gamma"] = createNode("gamma", "alpha")
	knownNodes["delta"] = createNode("delta", "beta")
	knownNodes["epsilon"] = createNode("epsilon", "alpha")

	// Trying to release a non-existent node should result in error
	nonExistentReq := structs.SimpleRequest{Callsign: "alpha", NodeId: "iota"}
	nonExistentResult := handlers.Release(nonExistentReq, knownNodes, knownBots)
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Non-existent node somehow mutated known bot claims: %d", len(knownBots["alpha"].Claims))
	}
	if !nonExistentResult.Error {
		t.Errorf("Non-existent node somehow resulted in successful response")
	}

	// Trying to release someone else's node should result in error and not affect the other bot
	unownedReq := structs.SimpleRequest{Callsign: "alpha", NodeId: "delta"}
	unownedResult := handlers.Release(unownedReq, knownNodes, knownBots)
	if len(knownBots["beta"].Claims) != 1 || len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Node owned by other bot somehow mutated requesting bots claims")
	}
	if !unownedResult.Error {
		t.Errorf("Unowned node somehow resulted in successful response")
	}

	// Trying to release your own node should result only in that node being released
	validReq := structs.SimpleRequest{Callsign: "alpha", NodeId: "epsilon"}
	validResult := handlers.Release(validReq, knownNodes, knownBots)
	if validResult.Error {
		t.Errorf("Valid node somehow resulted in error response")
	}
	if len(knownBots["alpha"].Claims) > 1 || knownNodes["epsilon"].ClaimedBy != "" {
		t.Errorf("Valid node somehow didn't release claim. Bot claims: %d, Node claimed by: %s",
			len(knownBots["alpha"].Claims), knownNodes["epsilon"].ClaimedBy)
	}

}

func TestClaim(t *testing.T) {
	grid := structs.Grid{}
	grid.Width = 100
	grid.Height = 100

	knownBots := make(map[string]structs.Bot)
	knownBots["alpha"] = createBot("alpha", []string{"gamma"})
	knownBots["beta"] = createBot("beta", []string{"delta"})

	knownNodes := make(map[string]structs.Node)
	knownNodes["gamma"] = createNode("gamma", "alpha")
	knownNodes["delta"] = createNode("delta", "beta")
	knownNodes["epsilon"] = createNode("epsilon", "")

	unclaimedReq := structs.SimpleRequest{
		Callsign: "alpha",
		NodeId:   "epsilon",
	}
	unclaimedResult := handlers.Claim(unclaimedReq, knownNodes, knownBots, grid)
	if unclaimedResult.Error {
		t.Errorf("Trying to claim an unclaimed node should result in success")
	}
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Claiming unclaimed node didn't add claim to bot's list of claims: %d", len(knownBots["alpha"].Claims))
	}
	if knownNodes["epsilon"].ClaimedBy != "alpha" {
		t.Errorf("Claiming node didn't add claim to node's property")
	}

	claimedReq := structs.SimpleRequest{
		Callsign: "alpha",
		NodeId:   "delta",
	}

	handlers.Claim(claimedReq, knownNodes, knownBots, grid)
	if knownNodes["delta"].ClaimedBy != "beta" {
		t.Errorf("Trying to claim a node claimed by someone else should result in failure")
	}
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Claiming node owned by another bot should result in error")
	}
	if knownNodes["delta"].ClaimedBy != "beta" {
		t.Errorf("Claiming node owned by other bot shouldn't change node's claim")
	}

	alreadyClaimedReq := structs.SimpleRequest{
		Callsign: "alpha",
		NodeId:   "epsilon",
	}
	alreadyClaimedResult := handlers.Claim(alreadyClaimedReq, knownNodes, knownBots, grid)
	if alreadyClaimedResult.Error {
		t.Errorf("Trying to claim a node already claimed should result in success")
	}
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Claiming existing claimed node didn't preserve claim list")
	}
	if knownNodes["epsilon"].ClaimedBy != "alpha" {
		t.Errorf("Claiming existing node should keep node's claim")
	}

}

func TestScan(t *testing.T) {
	grid := structs.Grid{}
	grid.Width = 100
	grid.Height = 100

	createNodeWithLocation := func(id string, x int, y int) structs.Node {
		node := createNode(id, "")
		node.Location = structs.GridLocation{
			X: x,
			Y: y,
		}
		return node
	}

	alpha := createBot("alpha", []string{})
	alpha.Location = structs.GridLocation{
		X: 10,
		Y: 10,
	}
	nodes := make(map[string]structs.Node)
	bots := make(map[string]structs.Bot)
	bots["alpha"] = alpha

	req := structs.SimpleRequest{
		Callsign: "alpha",
	}

	tt := []struct {
		x          int
		y          int
		expResults int
	}{
		{16, 10, 0},
		{4, 10, 0},
		{10, 16, 0},
		{10, 4, 0},
		{5, 10, 1},
		{15, 10, 1},
		{10, 15, 1},
		{10, 5, 1},
		{5, 5, 1},
		{100, 100, 0},
		{15, 15, 1},
		{16, 16, 0},
		{4, 4, 0},
	}

	for _, tc := range tt {
		node := createNodeWithLocation("gamma", tc.x, tc.y)
		nodes["gamma"] = node
		res := handlers.Scan(req, nodes, bots, grid)
		if len(res.Nodes) != tc.expResults {
			t.Errorf("Scan didn't return expected results. Expected: %d, Got: %d. TC: %v", tc.expResults, len(res.Nodes), tc)
		}
	}

}

/**
Helper functions
*/

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
