package structs

type EntityType int

const (
	BOT EntityType = iota
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
