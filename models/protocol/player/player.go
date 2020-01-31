package player

// Position of Player
type Position struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	Rotation int `json:"rotation"`
}

// Size of Player
type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Speed of Player
type Speed struct {
	X   int `json:"x"`
	Y   int `json:"y"`
	Max int `json:"max"`
}

// Effect of Player
type Effect struct {
	Kind    string `json:"kind"`
	CurTime int    `json:"cur_time"`
	ExpTime int    `json:"exp_time"`
	Rate    int    `json:"rate"`
	Speed   int    `json:"speed"`
	Health  int    `json:"health"`
	Armor   int    `json:"armor"`
}

// Ammunition of Player
type Ammunition struct {
	Weapon     string     `json:"weapon"`
	Cartridges Cartridges `json:"cartridges"`
}

// Cartridges of Player
type Cartridges struct {
	Blaster int `json:"blaster"`
	Plazma  int `json:"plazma"`
	Sniper  int `json:"sniper"`
}

// GameStats of Player
type GameStats struct {
	Kills int `json:"kills"`
	Death int `json:"death"`
}

// LifeStats of Player
type LifeStats struct {
	Health int `json:"health"`
	Armor  int `json:"armor"`
}

// Effects of Player
type Effects []Effect

// Player - main struct
type Player struct {
	Nickname   string     `json:"nickname"`
	Position   Position   `json:"position"`
	Size       Size       `json:"size"`
	Speed      Speed      `json:"speed"`
	Effects    Effects    `json:"effects"`
	Ammunition Ammunition `json:"ammunition"`
	GameStats  GameStats  `json:"game_stats"`
	LifeStats  LifeStats  `json:"life_stats"`
}
