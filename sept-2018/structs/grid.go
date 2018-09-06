package structs

import (
	"math"
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
	NUMBER_OF_NODES = 80
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
	// TODO: This doesn't cover wraparound

	minDist := int(math.Max(0, float64(value-distance)))
	maxDist := int(math.Min(float64(value+distance), float64(maxValue)))

	withinRange := func(min int, max int, number int) bool {
		return min <= number && max >= number
	}

	return withinRange(minDist, maxDist, testValue)

}

func (grid Grid) randomInitVals() (int, int, int) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	x := random.Intn(grid.Width + 1)
	y := random.Intn(grid.Height + 1)
	value := random.Intn(MAX_NODE_VALUE + 1)
	return x, y, value
}

// ScannableByBot returns true if the provided node is within scan range of
// the provided bot. False otherwise.
func (grid Grid) ScannableByBot(node Node, bot Bot) bool {

	return NumberWithinRange(bot.Location.X, SCAN_RANGE, grid.Width, node.Location.X) &&
		NumberWithinRange(bot.Location.Y, SCAN_RANGE, grid.Height, node.Location.Y)
}

func (grid Grid) MoveBot(bot Bot, x int, y int, debugMode bool) GridLocation {
	// Is this a valid move for the bot?
	validMove := NumberWithinRange(bot.Location.X, 1, grid.Width, x) &&
		NumberWithinRange(bot.Location.Y, 1, grid.Height, y)

	if validMove || debugMode {
		bot.Location = GridLocation{X: x, Y: y}
	}
	return bot.Location
}

func (grid Grid) InitializeBot(callsign string, debug bool) Bot {
	x, y, _ := grid.randomInitVals()
	bot := Bot{
		GridEntity: GridEntity{
			Id:   callsign,
			Type: BOT,
			Location: GridLocation{
				X: x,
				Y: y,
			},
		},
		DebugMode: debug,
		Claims:    []string{},
	}
	return bot
}

func (grid Grid) initializeNodes() map[string]Node {
	nodeMap := make(map[string]Node)
	for idx := 0; idx < NUMBER_OF_NODES; idx++ {
		x, y, value := grid.randomInitVals()

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
