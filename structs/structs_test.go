package structs_test

import (
	"testing"

	"github.com/HeadlightLabs/Tournament-API/structs"
)

func TestNumberWithinRange(t *testing.T) {

	var max int = 100

	tt := []struct {
		botValue       int
		nodeValue      int
		distance       int
		expectedResult bool
	}{
		{50, 95, 5, false},
		{50, 5, 5, false},
		{97, 99, 5, true},
		{97, 2, 5, false},
		{4, 2, 5, true},
		{4, 99, 5, false},
		{97, 3, 5, false},
		{3, 97, 5, false},
		{3, 4, 1, true},
		{50, 49, 1, true},
		{50, 48, 1, false},
	}

	for _, tc := range tt {
		actual := structs.NumberWithinRange(tc.botValue, tc.distance, max, tc.nodeValue)
		if tc.expectedResult != actual {
			t.Errorf("Number within range didn't return correct answer. Bot: %d nodeValue: %d Actual result: %v", tc.botValue, tc.nodeValue, actual)
		}
	}

}

func TestBotStatus(t *testing.T) {
	tt := []struct {
		Claims []string
		Id     string
		X      int
		Y      int
		Score  int
	}{
		{[]string{""}, "foobar", 1, 2, 3},
		{[]string{"alpha", "tango"}, "baz", 3, 4, 5},
	}

	for _, tc := range tt {
		bot := structs.Bot{
			Claims: tc.Claims,
			GridEntity: structs.GridEntity{
				Id: tc.Id,
				Location: structs.GridLocation{
					X: tc.X,
					Y: tc.Y,
				},
			},
			Score: tc.Score,
		}
		actual := bot.GetStatus()
		if len(tc.Claims) != len(actual.Claims) || tc.Id != actual.Id || tc.Score != actual.Score || tc.X != actual.Location.X || tc.Y != actual.Location.Y {
			t.Errorf("Bot Status didn't return correct answer. Expected: %v, Actual: %v", tc, actual)
		}
	}
}

func TestNodeStatus(t *testing.T) {
	tt := []struct {
		Id    string
		X     int
		Y     int
		Value int
	}{
		{"foobar", 1, 2, 3},
		{"baz", 3, 4, 5},
	}

	for _, tc := range tt {
		node := structs.Node{
			GridEntity: structs.GridEntity{
				Id: tc.Id,
				Location: structs.GridLocation{
					X: tc.X,
					Y: tc.Y,
				},
			},
			Value: tc.Value,
		}
		actual := node.GetStatus()
		if tc.Id != actual.Id || tc.Value != actual.Value || tc.X != actual.Location.X || tc.Y != actual.Location.Y {
			t.Errorf("Node Status didn't return correct answer. Expected: %v, Actual: %v", tc, actual)
		}
	}
}

func TestInitializeNodes(t *testing.T) {
	grid := structs.Grid{}
	nodes := grid.InitializeNodes()
	if len(nodes) < 1 {
		t.Errorf("Nodes were not initialized properly!")
	}
}

func TestInitializeBot(t *testing.T) {
	grid := structs.Grid{}
	bot := grid.InitializeBot("foobar")
	if bot.Id != "foobar" {
		t.Errorf("Bot was not initialized properly!")
	}
	if len(bot.Claims) > 0 {
		t.Errorf("Bot was initialized with claims, which is incorrect.")
	}
}

func TestMoveBot(t *testing.T) {
	grid := structs.Grid{}
	grid.Width = 100
	grid.Height = 100
	bot := structs.Bot{
		Claims: []string{},
		GridEntity: structs.GridEntity{
			Id: "foobar",
			Location: structs.GridLocation{
				X: 50,
				Y: 50,
			},
		},
		Score: 0,
	}

	tt := []struct {
		x    int
		y    int
		expX int
		expY int
	}{
		{49, 49, 49, 49},
		{49, 50, 49, 50},
		{49, 51, 49, 51},
		{50, 50, 50, 50},
		{50, 49, 50, 49},
		{50, 51, 50, 51},
		{51, 49, 51, 49},
		{51, 50, 51, 50},
		{51, 51, 51, 51},
		{51, 52, 50, 50},
		{48, 49, 50, 50},
	}

	for _, tc := range tt {
		bot.Location.X = 50
		bot.Location.Y = 50
		actual := grid.MoveBot(bot, tc.x, tc.y)
		if actual.X != tc.expX || actual.Y != tc.expY {
			t.Errorf("Move didn't return correct location. Exp: (%d,%d). Actual: (%d,%d)", tc.expX, tc.expY, actual.X, actual.Y)
		}
	}
}

func TestRandomInitVals(t *testing.T) {
	grid := structs.Grid{}
	grid.Width = 100
	grid.Height = 100
	x, y, value := grid.RandomInitVals()

	if value == 0 {
		t.Errorf("Value should never be 0")
	}
	if x == y {
		t.Errorf("Extremely unlikely that x and y would be the same.")
	}
}

func TestCheckMineValidity(t *testing.T) {
	grid := structs.Grid{}
	grid.Width = 100
	grid.Height = 100

	alpha := createBot("alpha", 5, 5, []string{})
	gamma := createNode("gamma", 5, 6, "beta")
	epsilon := createNode("epsilon", 13, 13, "alpha")
	theta := createNode("epsilon", 13, 13, "delta")
	delta := createBot("delta", 18, 18, []string{"theta"})

	tt := []struct {
		bot      structs.Bot
		node     structs.Node
		expError string
	}{
		{alpha, gamma, "already_claimed"}, // Node claimed by someone else
		{alpha, epsilon, "too_far_away"},  // Not within scan range
		{delta, theta, ""},                // Happy path
	}

	for _, tc := range tt {
		actual := grid.CheckMineValidity(tc.node, tc.bot)
		if tc.expError != actual {
			t.Errorf("CheckMineValidity didn't return expected error: %s, Actual: %s. TC: %v", tc.expError, actual, tc)
		}
	}

}

func TestCheckClaimValidity(t *testing.T) {
	grid := structs.Grid{}
	grid.Width = 100
	grid.Height = 100

	alpha := createBot("alpha", 5, 5, []string{})
	gamma := createNode("gamma", 5, 6, "beta")
	beta := createBot("beta", 10, 10, []string{"gamma", "zeta", "iota"})
	epsilon := createNode("epsilon", 13, 13, "")
	delta := createBot("delta", 18, 18, []string{})

	tt := []struct {
		bot      structs.Bot
		node     structs.Node
		expError string
	}{
		{alpha, gamma, "already_claimed"},  // Node claimed by someone else
		{beta, epsilon, "too_many_claims"}, // Bot has too many claims
		{alpha, epsilon, "too_far_away"},   // Not within scan range
		{delta, epsilon, ""},               // Happy path
	}

	for _, tc := range tt {
		actual := grid.CheckClaimValidity(tc.node, tc.bot)
		if tc.expError != actual {
			t.Errorf("CheckClaimValidity didn't return expected error: %s, Actual: %s. TC: %v", tc.expError, actual, tc)
		}
	}

}

func TestScannableByBot(t *testing.T) {
	grid := structs.Grid{}
	grid.Width = 100
	grid.Height = 100

	alpha := createBot("alpha", 0, 0, []string{})
	delta := createBot("delta", 18, 18, []string{})
	omega := createBot("omega", 19, 19, []string{})

	gamma := createNode("gamma", 0, 1, "beta")
	epsilon := createNode("epsilon", 13, 13, "")
	theta := createNode("theta", 100, 100, "")
	psi := createNode("psi", 0, 0, "")

	tt := []struct {
		bot       structs.Bot
		node      structs.Node
		expResult bool
	}{
		{alpha, epsilon, false}, // Not anywhere close
		{alpha, gamma, true},    // Next to the node
		{delta, epsilon, true},  // Just within range
		{omega, epsilon, false}, // Just out of range
		{alpha, theta, false},   // Overlapped should return false
		{alpha, psi, true},      // On top of node should return true
	}

	for _, tc := range tt {
		actual := grid.ScannableByBot(tc.node, tc.bot)
		if tc.expResult != actual {
			t.Errorf("ScannableByBot didn't return expected reult: %v, Actual: %v. TC: %v", tc.expResult, actual, tc)
		}
	}
}

func createBot(uuid string, x int, y int, claims []string) structs.Bot {
	bot := structs.Bot{
		GridEntity: structs.GridEntity{
			Id:   uuid,
			Type: structs.BOT,
			Location: structs.GridLocation{
				X: x,
				Y: y,
			},
		},
		Claims: claims,
	}
	return bot
}

func createNode(uuid string, x int, y int, claimedBy string) structs.Node {
	node := structs.Node{
		GridEntity: structs.GridEntity{
			Id:   uuid,
			Type: structs.NODE,
			Location: structs.GridLocation{
				X: x,
				Y: y,
			},
		},
		ClaimedBy: claimedBy,
		Value:     1,
	}
	return node
}
