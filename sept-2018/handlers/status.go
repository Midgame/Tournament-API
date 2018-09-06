package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

// Status returns information about the requesting user's:
// Location, Claims, Total score
// If in debug mode, also returns this information for all other known bots
func Status(req structs.SimpleRequest, nodes map[string]structs.Node, bots map[string]structs.Bot) structs.StatusResponse {
	resp := CheckParams(req, nodes, bots, false)
	if resp.Error {
		return resp
	}
	bot := bots[req.Callsign]

	// In debug mode, return all bots and nodes
	if req.DebugMode {
		for _, bot := range bots {
			resp.Bots = append(resp.Bots, bot.GetStatus())
		}
		for _, node := range nodes {
			resp.Nodes = append(resp.Nodes, node.GetStatus())
		}
		return resp
	}

	// In regular mode, just return the requested id
	resp.Bots = append(resp.Bots, bot.GetStatus())
	return resp

}
