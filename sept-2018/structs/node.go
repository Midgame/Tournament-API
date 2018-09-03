package structs

type Node struct {
	GridEntity
	ClaimedBy string
	Value     uint64
}

type NodeStatus struct {
	Id        string
	ClaimedBy string
	Location  GridLocation
}

// GetStatus returns some basic information about this node, including
// location, id, and any existing claims
func (node Node) GetStatus() NodeStatus {
	return NodeStatus{
		Id:        node.Id,
		Location:  node.Location,
		ClaimedBy: node.ClaimedBy,
	}
}
