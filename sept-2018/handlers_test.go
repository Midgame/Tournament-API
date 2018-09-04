package main

import (
	"testing"

	"github.com/HeadlightLabs/Tournament-API/sept-2018/handlers"
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

func TestRegisterUser(t *testing.T) {
	req := handlers.RegisterRequest{
		DebugMode: false,
		Callsign:  "",
	}
	bot, response := handlers.RegisterUser(req)
	if response.Callsign != bot.Id {
		t.Errorf("Response didn't return the correct UUID")
	}
	if bot.Id == "" {
		t.Errorf("Bot wasn't created with UUID properly")
	}
	if len(bot.Claims) != 0 {
		t.Errorf("Bot somehow started out with claims immediately after registration")
	}
	if bot.Location.X != 0 || bot.Location.Y != 0 {
		t.Errorf("Bot wasn't initialized with a location properly")
	}

	debugReq := handlers.RegisterRequest{
		DebugMode: true,
		Callsign:  "foobar",
	}
	debugBot, _ := handlers.RegisterUser(debugReq)
	if debugBot.Id != "foobar" {
		t.Errorf("Register function didn't accept callsign. Assigned: %s", bot.Id)
	}
	if !debugBot.DebugMode {
		t.Errorf("Register function didn't accept debug flag properly.")
	}
}

func TestStatus(t *testing.T) {
	knownBots := make(map[string]structs.Bot)
	knownBots["alpha"] = createBot("alpha", []string{})
	knownBots["beta"] = createBot("beta", []string{})
	gammaBot := createBot("gamma", []string{})
	gammaBot.DebugMode = true
	knownBots["gamma"] = gammaBot

	validReq := handlers.StatusRequest{Callsign: "beta"}
	validResult := handlers.Status(validReq, knownBots)
	if len(validResult.Bots) != 1 || validResult.Bots[0].Id != "beta" {
		t.Errorf("Non-debug result should find single bot with valid uuid. Bot found has ID: %s", validResult.Bots[0].Id)
	}

	invalidReq := handlers.StatusRequest{Callsign: "delta"}
	invalidResult := handlers.Status(invalidReq, knownBots)
	if len(invalidResult.Bots) > 0 {
		t.Errorf("Non-debug result should find no bots with invalid uuid, found: %d", len(invalidResult.Bots))
	}

	debugReq := handlers.StatusRequest{Callsign: "gamma"}
	debugResult := handlers.Status(debugReq, knownBots)
	if len(debugResult.Bots) != 3 {
		t.Errorf("Debug result should return all bots, found: %d", len(debugResult.Bots))
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
	nonExistentReq := handlers.ReleaseRequest{Callsign: "alpha", NodeId: "iota"}
	nonExistentResult := handlers.Release(nonExistentReq, knownNodes, knownBots)
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Non-existent node somehow mutated known bot claims: %d", len(knownBots["alpha"].Claims))
	}
	if nonExistentResult.Success {
		t.Errorf("Non-existent node somehow resulted in successful response")
	}

	// Trying to release someone else's node should result in error and not affect the other bot
	unownedReq := handlers.ReleaseRequest{Callsign: "alpha", NodeId: "delta"}
	unownedResult := handlers.Release(unownedReq, knownNodes, knownBots)
	if len(knownBots["beta"].Claims) != 1 || len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Node owned by other bot somehow mutated requesting bots claims")
	}
	if unownedResult.Success {
		t.Errorf("Unowned node somehow resulted in successful response")
	}

	// Trying to release your own node should result only in that node being released
	validReq := handlers.ReleaseRequest{Callsign: "alpha", NodeId: "epsilon"}
	validResult := handlers.Release(validReq, knownNodes, knownBots)
	if !validResult.Success {
		t.Errorf("Valid node somehow resulted in error response")
	}
	if len(knownBots["alpha"].Claims) > 1 || knownNodes["epsilon"].ClaimedBy != "" {
		t.Errorf("Valid node somehow didn't release claim. Bot claims: %d, Node claimed by: %s",
			len(knownBots["alpha"].Claims), knownNodes["epsilon"].ClaimedBy)
	}

}

func TestClaim(t *testing.T) {
	knownBots := make(map[string]structs.Bot)
	knownBots["alpha"] = createBot("alpha", []string{"gamma"})
	knownBots["beta"] = createBot("beta", []string{"delta"})

	knownNodes := make(map[string]structs.Node)
	knownNodes["gamma"] = createNode("gamma", "alpha")
	knownNodes["delta"] = createNode("delta", "beta")
	knownNodes["epsilon"] = createNode("epsilon", "")

	unclaimedReq := handlers.ClaimRequest{
		Callsign: "alpha",
		NodeId:   "epsilon",
	}
	unclaimedResult := handlers.Claim(unclaimedReq, knownNodes, knownBots)
	if !unclaimedResult.Success {
		t.Errorf("Trying to claim an unclaimed node should result in success")
	}
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Claiming unclaimed node didn't add claim to bot's list of claims: %d", len(knownBots["alpha"].Claims))
	}
	if knownNodes["epsilon"].ClaimedBy != "alpha" {
		t.Errorf("Claiming node didn't add claim to node's property")
	}

	claimedReq := handlers.ClaimRequest{
		Callsign: "alpha",
		NodeId:   "delta",
	}
	claimedResult := handlers.Claim(claimedReq, knownNodes, knownBots)
	if claimedResult.Success {
		t.Errorf("Trying to claim a node claimed by someone else should result in failure")
	}
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Claiming node owned by another bot should result in error")
	}
	if knownNodes["delta"].ClaimedBy != "beta" {
		t.Errorf("Claiming node owned by other bot shouldn't change node's claim")
	}

	alreadyClaimedReq := handlers.ClaimRequest{
		Callsign: "alpha",
		NodeId:   "epsilon",
	}
	alreadyClaimedResult := handlers.Claim(alreadyClaimedReq, knownNodes, knownBots)
	if !alreadyClaimedResult.Success {
		t.Errorf("Trying to claim a node already claimed should result in success")
	}
	if len(knownBots["alpha"].Claims) != 2 {
		t.Errorf("Claiming existing claimed node didn't preserve claim list")
	}
	if knownNodes["epsilon"].ClaimedBy != "alpha" {
		t.Errorf("Claiming existing node should keep node's claim")
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
