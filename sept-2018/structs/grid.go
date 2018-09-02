package structs

type EntityType int

const (
	BOT EntityType = iota
)

type GridEntity struct {
	Type EntityType
	Id   string
}
