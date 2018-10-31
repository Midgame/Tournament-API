package structs

import (
	"math"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type EntityType int

const (
	BOT EntityType = iota
	NODE
)

const (
	SCAN_RANGE      = 2
	GRID_WIDTH      = 19
	GRID_HEIGHT     = 19
	NUMBER_OF_NODES = 50
	MAX_NODE_VALUE  = 20
	MAX_CLAIMS      = 3
)

type GridEntity struct {
	Type     EntityType
	Id       string
	Location GridLocation
}

type Grid struct {
	Width  int
	Height int
	Bots   map[string]Bot
	Nodes  map[string]Node
}

type GridLocation struct {
	X int
	Y int
}

// Determines if a number is within a certain distance of a value.
func NumberWithinRange(value int, distance int, maxValue int, testValue int) bool {
	minDist := int(math.Max(0, float64(value-distance)))
	maxDist := int(math.Min(float64(value+distance), float64(maxValue)))

	withinRange := func(min int, max int, number int) bool {
		return min <= number && max >= number
	}

	return withinRange(minDist, maxDist, testValue)
}

func (grid Grid) RandomInitVals() (int, int, int) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	x := random.Intn(grid.Width + 1)
	y := random.Intn(grid.Height + 1)
	value := random.Intn(MAX_NODE_VALUE+1) + 1
	return x, y, value
}

// ScannableByBot returns true if the provided node is within scan range of
// the provided bot. False otherwise.
func (grid Grid) ScannableByBot(node Node, bot Bot) bool {

	return NumberWithinRange(bot.Location.X, SCAN_RANGE, grid.Width, node.Location.X) &&
		NumberWithinRange(bot.Location.Y, SCAN_RANGE, grid.Height, node.Location.Y)
}

// CheckMineValidity ensures that a bot that wants to mine a node is allowed to
func (grid Grid) CheckMineValidity(node Node, bot Bot) string {
	// If this node is owned by someone else or unowned, return an error except in debug mode
	if node.ClaimedBy != bot.Id {
		return "already_claimed"
	}

	// If this bot isn't within scan range of this node, return an error
	if !grid.ScannableByBot(node, bot) {
		return "too_far_away"
	}

	return ""
}

// CheckClaimValidity ensures that a claim on a node by a bot is valid.
func (grid Grid) CheckClaimValidity(node Node, bot Bot) string {

	// If this node is owned by someone else, return an error
	if node.ClaimedBy != "" && node.ClaimedBy != bot.Id {
		return "already_claimed"
	}

	// If this bot has too many claims already, return an error
	if len(bot.Claims) >= MAX_CLAIMS {
		return "too_many_claims"
	}

	// If this bot isn't within scan range of this node, return an error
	if !grid.ScannableByBot(node, bot) {
		return "too_far_away"
	}

	// Happy claim!
	return ""
}

// MoveBot returns a new location for the bot - this may be the same location if the requested
// coordinates are invalid!
func (grid Grid) MoveBot(bot Bot, x int, y int) GridLocation {
	// Is this a valid move for the bot?
	validMove := NumberWithinRange(bot.Location.X, 1, grid.Width, x) &&
		NumberWithinRange(bot.Location.Y, 1, grid.Height, y)

	if validMove {
		bot.Location = GridLocation{X: x, Y: y}
	} else {
		glog.Infof("[grid.MoveBot]: Invalid move for Bot: %s from (%d,%d) to (%d,%d).", bot.Id, bot.Location.X, bot.Location.Y, x, y)
	}
	return bot.Location
}

func (grid Grid) InitializeBot(callsign string) Bot {
	x, y, _ := grid.RandomInitVals()
	bot := Bot{
		GridEntity: GridEntity{
			Id:   callsign,
			Type: BOT,
			Location: GridLocation{
				X: x,
				Y: y,
			},
		},
		Claims: []string{},
	}
	return bot
}

func (grid Grid) InitializeNodes() map[string]Node {
	nodeMap := make(map[string]Node)
	for idx := 0; idx < NUMBER_OF_NODES; idx++ {
		x, y, value := grid.RandomInitVals()

		node := Node{
			GridEntity: GridEntity{
				Id:   uuid.New().String(),
				Type: NODE,
				Location: GridLocation{
					X: x,
					Y: y,
				},
			},
			ClaimedBy: "",
			Value:     value,
		}
		nodeMap[node.Id] = node
	}

	return nodeMap
}

func (grid *Grid) Initialize() {
	grid.Width = GRID_WIDTH
	grid.Height = GRID_HEIGHT
	grid.Bots = make(map[string]Bot)
	grid.Nodes = grid.InitializeNodes()
}
