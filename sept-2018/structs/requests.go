package structs

type MoveRequest struct {
	Callsign string `json:"callsign"`
	X        int    `json:"x,string"`
	Y        int    `json:"y,string"`
}

type SimpleRequest struct {
	Callsign  string `json:"callsign"`
	DebugMode bool   `json:"debug,string,omitempty"`
	NodeId    string `json:"node,omitempty"`
}
