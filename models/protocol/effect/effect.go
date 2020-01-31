package effect

// Position of Effect
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size of Effect
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Data of Effect
type Data struct {
	Kind   string `json:"kind"`
	Rate   int    `json:"rate"`
	Health int    `json:"health"`
	Count  int    `json:"count"`
	Speed  int    `json:"speed"`
	Time   int    `json:"time"`
	Armor  int    `json:"armor"`
}

// Effect - main struct
type Effect struct {
	Position Position `json:"position"`
	Size     Size     `json:"size"`
	Data     Data     `json:"data"`
}
