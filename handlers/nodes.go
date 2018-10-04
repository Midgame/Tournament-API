package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/structs"
)

func Nodes(nodes map[string]structs.Node) []structs.NodeStatus {

	response := []structs.NodeStatus{}

	for _, node := range nodes {
		response = append(response, node.GetStatus())
	}

	return response
}
