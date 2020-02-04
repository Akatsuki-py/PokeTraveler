package stage

// Warps ワープポイントのリスト
type Warps struct {
	List []*Warp `json:"warps"`
}

// Warp ワープポイント
type Warp struct {
	X   int    `json:"x"`
	Y   int    `json:"y"`
	Dst string `json:"dst"`
	Pos [2]int `json:"pos"`
}
