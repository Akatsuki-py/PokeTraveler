package stage

type rawStage struct {
	Width    int            `json:"width"`
	Height   int            `json:"height"`
	Layers   []rawLayer     `json:"layers"`
	Tilesets []rawTileStart `json:"tilesets"`
}

type rawLayer struct {
	Data []int `json:"data"`
}

// TileStart 各タイルセットのベースIDとそのパス
type rawTileStart struct {
	FirstGID int    `json:"firstgid"`
	Source   string `json:"source"`
}
