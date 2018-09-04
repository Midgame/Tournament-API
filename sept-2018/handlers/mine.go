package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

type MineRequest struct {
	Callsign string `json:"callsign"`
	NodeId   string `json:"node"`
}

type MineResponse struct {
	Callsign        string
	NodeId          string
	Error           bool
	AmountMined     uint64
	AmountRemaining uint64
}

// Mine extracts some resources from a given node.
// The amount extracted is deducted from the node, given to the bot, and returned in the response.
// The response "Error" flag will be set to true if the callsign doesn't own this node
func Mine(req MineRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) MineResponse {

	resp := MineResponse{
		Callsign:        req.Callsign,
		NodeId:          req.NodeId,
		Error:           false,
		AmountMined:     0,
		AmountRemaining: 0,
	}

	// Return an error if this node does not exist
	node, ok := knownNodes[req.NodeId]
	if !ok {
		resp.Error = true
		return resp
	}

	// Return an error if this bot does not exist
	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		return resp
	}

	// If this node is owned by someone else or unowned, return an error
	if node.ClaimedBy != req.Callsign {
		resp.Error = true
		return resp
	}

	// If this node has no value left, return an error
	if node.Value <= 0 {
		resp.Error = true
		return resp
	}

	// The happy path: This callsign owns this node
	node.Value -= 1
	bot.Score += 1
	resp.AmountMined = 1
	resp.AmountRemaining = node.Value

	// Update the game
	knownBots[req.Callsign] = bot
	knownNodes[req.NodeId] = node

	return resp
}
