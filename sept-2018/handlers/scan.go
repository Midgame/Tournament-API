package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

type ScanResponse struct {
	Nodes []structs.NodeStatus
}

// Scan releases information about the nodes surrounding the requestor.
// Returns all nodes within a 5x5 grid around the requestor
func Scan(requestorId string, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) ScanResponse {
	return ScanResponse{
		Nodes: []structs.NodeStatus{},
	}
}
