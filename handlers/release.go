package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/structs"
)

// Release releases a claim on a node
// If this is not a node owned by the requestor, returns an error.
// Returns an :ok otherwise.
func Release(req structs.SimpleRequest, nodes map[string]structs.Node, bots map[string]structs.Bot) structs.StatusResponse {

	resp := CheckParams(req, nodes, bots, true)
	if resp.Error {
		return resp
	}
	bot := bots[req.Callsign]
	node := nodes[req.NodeId]
	resp.Error = true
	resp.ErrorMsg = NOT_CLAIMED_ERROR

	// Check if this node is owned by the requestor
	for idx, claim := range bots[req.Callsign].Claims {
		if claim == node.Id {

			// Update the node and the bot
			node.ClaimedBy = ""
			nodes[node.Id] = node

			// Slice removal
			bot.Claims[idx] = bot.Claims[len(bot.Claims)-1]
			bot.Claims[len(bot.Claims)-1] = ""
			bot.Claims = bot.Claims[:len(bot.Claims)-1]
			bots[req.Callsign] = bot

			// Update response
			resp.Error = false
			resp.ErrorMsg = ""
		}
	}

	resp.Status = bot.GetStatus()
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}

	return resp

}
