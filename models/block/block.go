package block

// Position of Block
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size of Block
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Data of Block
type Data struct {
	Kind   string `json:"kind"`
	Health int    `json:"health"`
}

// Block - main struct
type Block struct {
	Position Position `json:"position"`
	Size     Size     `json:"size"`
	Data     Data     `json:"data"`
}
