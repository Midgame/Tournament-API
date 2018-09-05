package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Claim establishes a claim on a node
// If this node is currently owned by another bot, returns Success: false.
// If this node does not exist, returns Error: true.
// Returns a success otherwise.
func Claim(req structs.SimpleRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) structs.ClaimResponse {

	resp := structs.ClaimResponse{
		Callsign: req.Callsign,
		NodeId:   req.NodeId,
		Error:    false,
	}

	// Return an error if this node does not exist
	node, ok := knownNodes[req.NodeId]
	if !ok {
		resp.Error = true
		return resp
	}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		return resp
	}

	// If this node is owned by someone else, return
	if node.ClaimedBy != "" && node.ClaimedBy != req.Callsign {
		resp.Success = false
		return resp
	}

	if node.ClaimedBy == req.Callsign {
		resp.Success = true
		return resp
	}

	node.ClaimedBy = req.Callsign
	bot.Claims = append(bot.Claims, req.NodeId)
	knownNodes[req.NodeId] = node
	knownBots[req.Callsign] = bot

	resp.Success = true
	return resp
}
