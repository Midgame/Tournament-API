package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/structs"
)

// Claim establishes a claim on a node
// If this node is currently owned by another bot, returns Success: false.
// If this node does not exist, returns Error: true.
// Returns a success otherwise.
func Claim(req structs.SimpleRequest, nodes map[string]structs.Node, bots map[string]structs.Bot, grid structs.Grid) structs.StatusResponse {

	resp := CheckParams(req, nodes, bots, true)
	if resp.Error {
		return resp
	}
	bot := bots[req.Callsign]
	node := nodes[req.NodeId]

	// If this is a noop, return
	if node.ClaimedBy == req.Callsign {
		return resp
	}

	err := grid.CheckClaimValidity(node, bot)
	if err != "" {
		resp.Error = true
		resp.ErrorMsg = err
		return resp
	}

	bot.Claims = append(bot.Claims, node.Id)
	node.ClaimedBy = bot.Id
	nodes[req.NodeId] = node
	bots[req.Callsign] = bot
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}
	resp.Status = bot.GetStatus()

	return resp
}
