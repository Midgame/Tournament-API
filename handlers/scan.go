package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/structs"
)

// Scan releases information about the nodes surrounding the requestor.
// Returns all nodes within a 5x5 grid around the requestor
func Scan(req structs.SimpleRequest, nodes map[string]structs.Node, bots map[string]structs.Bot, grid structs.Grid) structs.StatusResponse {
	nodeList := []structs.NodeStatus{}

	resp := CheckParams(req, nodes, bots, false)
	if resp.Error {
		return resp
	}
	bot := bots[req.Callsign]

	for _, node := range nodes {
		if grid.ScannableByBot(node, bot) {
			nodeList = append(nodeList, node.GetStatus())
		}
	}

	resp.Nodes = nodeList
	return resp
}
