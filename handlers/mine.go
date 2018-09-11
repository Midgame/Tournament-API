package handlers

import (
	"math"

	"github.com/HeadlightLabs/Tournament-API/structs"
)

// Mine extracts some resources from a given node.
// The amount extracted is deducted from the node, given to the bot, and returned in the response.
// The response "Error" flag will be set to true if the callsign doesn't own this node
func Mine(req structs.SimpleRequest, nodes map[string]structs.Node, bots map[string]structs.Bot, grid structs.Grid) structs.StatusResponse {

	resp := CheckParams(req, nodes, bots, true)
	if resp.Error {
		return resp
	}
	bot := bots[req.Callsign]
	node := nodes[req.NodeId]

	err := grid.CheckMineValidity(node, bot)
	if err != "" {
		resp.Error = true
		resp.ErrorMsg = err
		return resp
	}

	// If this node has no value left, this is a no-op except in debug mode
	if node.Value > 0 {
		bot.Score += 1
	}
	node.Value = int(math.Max(float64(0), float64(node.Value-1)))

	// Update the game and response
	bots[req.Callsign] = bot
	nodes[req.NodeId] = node
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}
	resp.Status = bot.GetStatus()

	return resp
}
