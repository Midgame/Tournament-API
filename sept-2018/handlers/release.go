package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

type ReleaseRequest struct {
	Callsign string `json:"callsign"`
	NodeId   string `json:"node"`
}

type ReleaseResponse struct {
	Success bool
	Error   bool
}

// Release releases a claim on a node
// If this is not a node owned by the requestor, returns an error.
// Returns an :ok otherwise.
func Release(req ReleaseRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) ReleaseResponse {

	resp := ReleaseResponse{
		Error:   true,
		Success: false,
	}

	// Return an error if this node does not exist or the bot does not exist
	node, ok := knownNodes[req.NodeId]
	if !ok {
		return resp
	}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		return resp
	}

	// Check if this node is owned by the requestor
	for idx, claim := range knownBots[req.Callsign].Claims {
		if claim == node.Id {

			// Update the node and the bot
			node.ClaimedBy = ""
			knownNodes[node.Id] = node

			// Slice removal
			bot.Claims[idx] = bot.Claims[len(bot.Claims)-1]
			bot.Claims[len(bot.Claims)-1] = ""
			bot.Claims = bot.Claims[:len(bot.Claims)-1]
			knownBots[req.Callsign] = bot

			resp.Error = false
			resp.Success = true
		}
	}

	// If not, error out
	return resp

	// TODO: Should maybe calculate final score?
}
