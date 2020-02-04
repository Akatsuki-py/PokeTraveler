package stage

type Actions struct {
	List []*Action `json:"actions"`
}

// Action アクションの詳細
type Action struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
