package structs

type ClaimResponse struct {
	Callsign string
	NodeId   string
	Error    bool
	Success  bool
}

type MineResponse struct {
	Callsign        string
	NodeId          string
	Error           bool
	AmountMined     int
	AmountRemaining int
}

type MoveResponse struct {
	Location GridLocation
	Error    bool
}

type RegisterResponse struct {
	Callsign  string
	DebugMode bool
}

type ReleaseResponse struct {
	Success bool
	Error   bool
}

type ScanResponse struct {
	Nodes []NodeStatus
	Error bool
}

type StatusResponse struct {
	Bots  []BotStatus
	Error bool
}
