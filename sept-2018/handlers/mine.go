package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Mine extracts some resources from a given node.
// The amount extracted is deducted from the node, given to the bot, and returned in the response.
// The response "Error" flag will be set to true if the callsign doesn't own this node
func Mine(req structs.SimpleRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) structs.StatusResponse {

	resp := structs.StatusResponse{
		Bots:  []structs.BotStatus{},
		Nodes: []structs.NodeStatus{},
		Error: false,
	}

	// Return an error if this node does not exist
	node, ok := knownNodes[req.NodeId]
	if !ok {
		resp.Error = true
		resp.ErrorMsg = "node_not_found"
		return resp
	}
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}

	// Return an error if this bot does not exist
	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		resp.ErrorMsg = "bot_not_found"
		return resp
	}
	resp.Bots = []structs.BotStatus{bot.GetStatus()}

	// If this node is owned by someone else or unowned, return an error
	if node.ClaimedBy != req.Callsign {
		resp.Error = true
		resp.ErrorMsg = "node_already_claimed"
		return resp
	}

	// If this node has no value left, this is a no-op
	if node.Value > 0 {
		node.Value -= 1
		bot.Score += 1
	}

	// Update the game and response
	knownBots[req.Callsign] = bot
	knownNodes[req.NodeId] = node
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}
	resp.Bots = []structs.BotStatus{bot.GetStatus()}

	return resp
}
