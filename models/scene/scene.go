package scene

// Position of Scene
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size of Scene
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Scene - main struct
type Scene struct {
	Position Position `json:"position"`
	Size     Size     `json:"size"`
}
