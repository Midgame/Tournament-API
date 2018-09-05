package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Release releases a claim on a node
// If this is not a node owned by the requestor, returns an error.
// Returns an :ok otherwise.
func Release(req structs.SimpleRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) structs.StatusResponse {

	resp := structs.StatusResponse{
		Error: true,
	}

	// Return an error if this node does not exist or the bot does not exist
	node, ok := knownNodes[req.NodeId]
	if !ok {
		return resp
	}
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		return resp
	}
	resp.Bots = []structs.BotStatus{bot.GetStatus()}

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
		}
	}

	resp.Bots = []structs.BotStatus{bot.GetStatus()}
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}

	return resp

}
