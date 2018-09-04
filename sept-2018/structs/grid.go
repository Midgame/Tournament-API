package structs

type EntityType int

const (
	BOT EntityType = iota
	NODE
)

const SCAN_RANGE = 5

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

// Determines if a number is within a certain distance of a value. Accounts for wrapping around
// (hence the maxValue param) in either direction.
func NumberWithinRange(value int, distance int, maxValue int, testValue int) bool {
	minDist := value - distance
	maxDist := value + distance

	// Test cases:
	// 1) node too far to the right
	// 2) node too far to the left
	// 3) node too far up
	// 4) node too far down
	// 5) node too far up but within left/right range
	// 6) node too far left but within up/down range
	// 9) node on left edge, before/after overlap (but within range)
	// 10) node on right edge, before/after overlap (but within range)
	// 11) node on top edge, before/after overlap (but within range)
	// 12) node on bottom edge, before/after overlap (within range)
	// 13) node on left/right/top/bottom edge, after overlap (not within range)
	// 14) node on right edge, just within scan range (exactly 5 units away)

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

	return NumberWithinRange(int(bot.Location.X), SCAN_RANGE, int(grid.Width), int(node.Location.X)) &&
		NumberWithinRange(int(bot.Location.Y), SCAN_RANGE, int(grid.Height), int(node.Location.Y))
}
