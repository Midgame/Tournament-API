package structs

type Node struct {
	GridEntity
	ClaimedBy string
	Value     int
}

type NodeStatus struct {
	Id       string
	Location GridLocation
	Value    int
}

// GetStatus returns some basic information about this node, including
// location, id, and any existing claims
func (node Node) GetStatus() NodeStatus {
	return NodeStatus{
		Id:       node.Id,
		Location: node.Location,
		Value:    node.Value,
	}
}
