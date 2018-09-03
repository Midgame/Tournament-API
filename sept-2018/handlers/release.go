package handlers

import (
	"github.com/HeadlightLabs/Tournament-API/sept-2018/structs"
)

type ReleaseResponse struct {
	Success bool
}

// Release releases a claim on a node
// If this is not a node owned by the requestor, returns an error.
// Returns an :ok otherwise.
func Release(requestorId string, nodeId string, knownNodes map[string]structs.Node, knownBots map[string]structs.Bot) ReleaseResponse {

	bot := knownBots[requestorId]
	node := knownNodes[nodeId]

	// Check if this node is owned by the requestor
	for idx, claim := range knownBots[requestorId].Claims {
		if claim == nodeId {

			// Update the node and the bot
			node.ClaimedBy = ""
			knownNodes[nodeId] = node

			// Slice removal
			bot.Claims[idx] = bot.Claims[len(bot.Claims)-1]
			bot.Claims[len(bot.Claims)-1] = ""
			bot.Claims = bot.Claims[:len(bot.Claims)-1]
			knownBots[requestorId] = bot

			return ReleaseResponse{
				Success: true,
			}
		}
	}

	// If not, error out
	return ReleaseResponse{
		Success: false,
	}

	// TODO: Should maybe calculate final score?
}
