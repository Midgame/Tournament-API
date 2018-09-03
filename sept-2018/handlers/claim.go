package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

type ClaimResponse struct {
	Success bool
}

// Claim establishes a claim on a node
// If this node is currently owned by another bot, returns an error.
// Returns a success otherwise.
func Claim(requestorId string, nodeId string, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) ClaimResponse {

	bot := knownBots[requestorId]
	node := knownNodes[nodeId]

	// If this node is owned by someone else, return an error
	if node.ClaimedBy != "" && node.ClaimedBy != requestorId {
		return ClaimResponse{
			Success: false,
		}
	}

	if node.ClaimedBy == requestorId {
		return ClaimResponse{
			Success: true,
		}
	}

	node.ClaimedBy = requestorId
	bot.Claims = append(bot.Claims, nodeId)
	knownNodes[nodeId] = node
	knownBots[requestorId] = bot

	return ClaimResponse{
		Success: true,
	}

}
