package structs

type MoveRequest struct {
	SimpleRequest
	X int `json:"x"`
	Y int `json:"y"`
}

type SimpleRequest struct {
	Callsign  string `json:"callsign"`
	DebugMode bool   `json:"debug,string,omitempty"`
	NodeId    string `json:"node,omitempty"`
}
