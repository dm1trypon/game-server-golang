package bullet

// Position of Bullet
type Position struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	Rotation int `json:"rotation"`
}

// Size of Bullet
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Speed of Bullet
type Speed struct {
	X   int `json:"x"`
	Y   int `json:"y"`
	Max int `json:"max"`
}

// Data of Bullet
type Data struct {
	Nickname string `json:"nickname"`
	Weapon   string `json:"weapon"`
	Health   int    `json:"health"`
}

// Bullet - main struct
type Bullet struct {
	Position Position `json:"position"`
	Size     Size     `json:"size"`
	Speed    Speed    `json:"speed"`
	Data     Data     `json:"data"`
}
