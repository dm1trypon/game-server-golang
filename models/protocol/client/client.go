package client

// Size of Init
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Screen of Init
type Screen struct {
	Size Size
}

// Init - main struct
type Init struct {
	Nickname string `json:"nickname"`
	Screen   Screen `json:"screen"`
}

// KeyBoard - main struct
type KeyBoard struct {
	Nickname string   `json:"nickname"`
	Method   string   `json:"method"`
	Keys     []string `json:"keys"`
}

// Position of Mouse
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Mouse - main struct
type Mouse struct {
	Nickname  string   `json:"x"`
	Method    string   `json:"method"`
	IsClicked bool     `json:"is_clicked"`
	Position  Position `json:"position"`
}
