package stage

import (
	"encoding/json"
	"io/ioutil"
)

type Actions struct {
	List []*Action `json:"actions"`
}

// Action アクションの詳細
type Action struct {
	X     int      `json:"x"`
	Y     int      `json:"y"`
	Type  string   `json:"type"`
	Value []string `json:"value"`
}

// GetAction - Get Action
func (stage *Stage) GetAction(x, y int) (target *Action) {
	for _, action := range stage.Actions {
		if action.X == x/16 && action.Y == y/16 {
			target = action
			break
		}
	}
	return target
}

func (stage *Stage) loadActions(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	raw := new(Actions)
	if err := json.Unmarshal(file, raw); err != nil {
		panic(err)
	}
	stage.Actions = raw.List
}
