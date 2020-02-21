package config

// Net struct contains settings of server's ports
type Net struct {
	UDPPath string `json:"udp_path"`
	TCPPath string `json:"tcp_path"`
}

// Game struct contains game's setting
type Game struct {
	FPS        int `json:"fps"`
	MaxPlayers int `json:"max_players"`
}

// Player struct contains base player's setting
type Player struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Speed      int    `json:"speed"`
	Weapon     string `json:"weapon"`
	Cartridges int    `json:"cartridges"`
	Health     int    `json:"health"`
	Armor      int    `json:"armor"`
}

// Bullet struct contains base bullet's setting
type Bullet struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Speed  int    `json:"speed"`
	Weapon string `json:"weapon"`
	Health int    `json:"health"`
	Rate   int    `json:"rate"`
	TTL    int    `json:"ttl"`
}

// Block struct contains base block's setting
type Block struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Health int    `json:"health"`
	Kind   string `json:"kind"`
}

// Effect struct contains base effect's setting
type Effect struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Health []int  `json:"health"`
	Armor  []int  `json:"armor"`
	Rate   []int  `json:"rate"`
	Speed  []int  `json:"speed"`
	Сount  []int  `json:"count"`
	Kind   string `json:"kind"`
}

// Scene struct contains base scene's setting
type Scene struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// GameObjects struct contains base game's objects's setting
type GameObjects struct {
	Player  Player   `json:"player"`
	Bullets []Bullet `json:"bullets"`
	Blocks  []Block  `json:"blocks"`
	Effects []Effect `json:"effects"`
	Scene   Scene    `json:"scene"`
}

// Config struct - main struct of config
type Config struct {
	Net         Net         `json:"net"`
	Game        Game        `json:"game"`
	GameObjects GameObjects `json:"game_objects"`
}
