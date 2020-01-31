package protoworker

import "encoding/json"

// Init struct contains data for init player
type Init struct {
	Method     string
	Nickname   string
	Resolution Resolution
}

// Resolution struct contains resolution
type Resolution struct {
	Width  int
	Height int
}

// Position - position of cursor
type Position struct {
	X int
	Y int
}

// Mouse - position of cursor
type Mouse struct {
	Nickname  string
	Method    string
	IsClicked bool
	Position  Position
}

// Keyboard - events from player's keyboard
type Keyboard struct {
	Nickname string
	Method   string
	Keys     []string
}

// OnTCPMessage method parse data from client
func OnTCPMessage(data []byte) {
	init := &Init{}
	mouse := &Mouse{}
	keyboard := &Keyboard{}

	err := json.Unmarshal(data, &init)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &mouse)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &keyboard)
	if err != nil {
		return
	}

}
