package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/structs"
)

func Nodes(nodes map[string]structs.Node) map[string][]structs.NodeStatus {

	response := []structs.NodeStatus{}

	for _, node := range nodes {
		response = append(response, node.GetStatus())
	}

	responseMap := make(map[string][]structs.NodeStatus)
	responseMap["Nodes"] = response
	return responseMap
}
