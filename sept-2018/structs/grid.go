package structs

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type EntityType int

const (
	BOT EntityType = iota
	NODE
)

const (
	SCAN_RANGE      = 5
	GRID_WIDTH      = 100
	GRID_HEIGHT     = 100
	NUMBER_OF_NODES = 20
	MAX_NODE_VALUE  = 20
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

// Determines if a number is within a certain distance of a value. Accounts for wrapping around
// (hence the maxValue param) in either direction.
func NumberWithinRange(value int, distance int, maxValue int, testValue int) bool {
	minDist := value - distance
	maxDist := value + distance

	withinRange := func(min int, max int, number int) bool {
		return min <= number && max >= number
	}

	// Easy case: No overlap
	if minDist >= 0 && maxDist <= maxValue {
		return withinRange(minDist, maxDist, testValue)
	}

	// If min is < 0:
	// Is this between 0, original value (no -> maybe return false if next is also false)
	// Or else is this between (max + min), max? (no -> return false)
	if minDist < 0 {
		beforeOverlap := withinRange(0, value, testValue)
		afterOverlap := withinRange(maxValue+minDist, maxValue, testValue)
		if !(beforeOverlap || afterOverlap) {
			return false
		}
	}

	// If maxDist is > maxValue:
	// Is this between original value, max? (no -> maybe return false if next is false)
	// Is this between 0, (max % gridMax)? (no -> return false)
	if maxDist > maxValue {
		beforeOverlap := withinRange(value, maxValue, testValue)
		afterOverlap := withinRange(0, maxDist%maxValue, testValue)
		if !(beforeOverlap || afterOverlap) {
			return false
		}
	}

	return true
}

// ScannableByBot returns true if the provided node is within scan range of
// the provided bot. False otherwise.
func (grid Grid) ScannableByBot(node Node, bot Bot) bool {

	return NumberWithinRange(bot.Location.X, SCAN_RANGE, grid.Width, node.Location.X) &&
		NumberWithinRange(bot.Location.Y, SCAN_RANGE, grid.Height, node.Location.Y)
}

func (grid Grid) MoveBot(bot Bot, x int, y int) GridLocation {
	// Is this a valid move for the bot?
	validMove := NumberWithinRange(bot.Location.X, 1, grid.Width, x) &&
		NumberWithinRange(int(bot.Location.Y), 1, grid.Height, y)

	if validMove {
		bot.Location = GridLocation{X: x, Y: y}
	}
	return bot.Location
}

func (grid Grid) initializeNodes() map[string]Node {
	nodeMap := make(map[string]Node)
	for idx := 0; idx < NUMBER_OF_NODES; idx++ {
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		x := random.Intn(grid.Width + 1)
		y := random.Intn(grid.Height + 1)
		value := random.Intn(MAX_NODE_VALUE + 1)

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
	grid.Nodes = grid.initializeNodes()
}
