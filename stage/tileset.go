package stage

type TileSet struct {
	Image string `json:"image"`
	List  []Tile `json:"tiles"`
}

type Tile struct {
	ID         int            `json:"id"`
	Properties []JSONProperty `json:"properties"`
}

type JSONProperty struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}
