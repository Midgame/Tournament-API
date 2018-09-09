package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

const (
	CALLSIGN_INVALID      = "callsign_invalid"
	BOT_NOT_FOUND_ERROR   = "bot_not_found"
	NODE_NOT_FOUND_ERROR  = "node_not_found"
	ALREADY_CLAIMED_ERROR = "already_claimed"
	NOT_CLAIMED_ERROR     = "not_claimed"
)

func removeFromSlice(key string, arr []string) []string {
	for idx, val := range arr {
		if val == key {
			arr[idx] = arr[len(arr)-1]
			arr[len(arr)-1] = ""
			return arr[:len(arr)-1]
		}
	}
	return arr
}

func CheckParams(req structs.SimpleRequest, nodes map[string]structs.Node, bots map[string]structs.Bot, checkNodes bool) structs.StatusResponse {
	resp := structs.StatusResponse{
		Error: false,
		Nodes: []structs.NodeStatus{},
	}

	// All handlers care about callsigns, except the registration method
	bot, ok := bots[req.Callsign]
	if !ok {
		resp.Error = true
		resp.ErrorMsg = BOT_NOT_FOUND_ERROR
		return resp
	}
	resp.Status = bot.GetStatus()

	// Some handlers don't care about nodes
	if !checkNodes {
		return resp
	}

	// Return an error if this node does not exist
	node, ok := nodes[req.NodeId]
	if !ok {
		resp.Error = true
		resp.ErrorMsg = ALREADY_CLAIMED_ERROR
		return resp
	}
	resp.Nodes = []structs.NodeStatus{node.GetStatus()}

	return resp
}
