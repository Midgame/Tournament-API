package structs

type EntityType int

const (
	BOT EntityType = iota
	NODE
)

type GridEntity struct {
	Type     EntityType
	Id       string
	Location GridLocation
}

type Grid struct {
	Width    uint64
	Height   uint64
	Entities [][]GridEntity
}

type GridLocation struct {
	X uint64
	Y uint64
}

// ScannableByBot returns true if the provided node is within scan range of
// the provided bot. False otherwise.
func (grid Grid) ScannableByBot(node Node, bot Bot) bool {
	return false
}
