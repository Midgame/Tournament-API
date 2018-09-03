package structs

type Node struct {
	GridEntity
	ClaimedBy string
	Value     uint64
}
