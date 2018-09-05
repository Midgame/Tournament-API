package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Claim establishes a claim on a node
// If this node is currently owned by another bot, returns Success: false.
// If this node does not exist, returns Error: true.
// Returns a success otherwise.
func Claim(req structs.SimpleRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) structs.StatusResponse {

	resp := structs.StatusResponse{
		Error: false,
	}

	// Return an error if this node does not exist
	node, ok := knownNodes[req.NodeId]
	if !ok {
		resp.Error = true
		return resp
	}
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		return resp
	}
	resp.Bots = []structs.BotStatus{bot.GetStatus()}

	// If this node is owned by someone else, return
	if node.ClaimedBy != "" && node.ClaimedBy != req.Callsign {
		return resp
	}

	// If this is a noop, return
	if node.ClaimedBy == req.Callsign {
		return resp
	}

	node.ClaimedBy = req.Callsign
	bot.Claims = append(bot.Claims, req.NodeId)
	knownNodes[req.NodeId] = node
	knownBots[req.Callsign] = bot
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}
	resp.Bots = []structs.BotStatus{bot.GetStatus()}

	return resp
}
