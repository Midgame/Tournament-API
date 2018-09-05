package structs

type StatusResponse struct {
	Bots      []BotStatus
	Nodes     []NodeStatus
	DebugMode bool
	Error     bool
	ErrorMsg  string
}
