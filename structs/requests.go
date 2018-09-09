package structs

type MoveRequest struct {
	Callsign string `json:"callsign"`
	X        int    `json:"x,string"`
	Y        int    `json:"y,string"`
}

type SimpleRequest struct {
	Callsign string `json:"callsign"`
	NodeId   string `json:"node,omitempty"`
}
