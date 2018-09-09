package structs

type StatusResponse struct {
	Status   BotStatus
	Nodes    []NodeStatus
	Error    bool
	ErrorMsg string
}
