package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Scan releases information about the nodes surrounding the requestor.
// Returns all nodes within a 5x5 grid around the requestor
func Scan(req structs.SimpleRequest, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot, grid structs.Grid) structs.StatusResponse {
	nodeList := []structs.NodeStatus{}
	resp := structs.StatusResponse{
		Error: false,
	}

	bot, ok := knownBots[req.Callsign]
	if !ok {
		resp.Error = true
		return resp
	}

	for _, node := range knownNodes {
		if grid.ScannableByBot(node, bot) {
			nodeList = append(nodeList, node.GetStatus())
		}
	}

	resp.Nodes = nodeList
	return resp
}
