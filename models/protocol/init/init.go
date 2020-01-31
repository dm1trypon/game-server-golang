package init

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
