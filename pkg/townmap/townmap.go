package townmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	cursorImage, _, _ = ebitenutil.NewImageFromFile("asset/townmap/asset/cursor.png", ebiten.FilterDefault)
	pointImage, _, _  = ebitenutil.NewImageFromFile("asset/townmap/asset/point.png", ebiten.FilterDefault)
	townImage, _, _   = ebitenutil.NewImageFromFile("asset/townmap/asset/town.png", ebiten.FilterDefault)
)

type TownMap struct {
	Regions map[string]*Region
}

type Region struct {
	Points []Point `json:"points"`
	Image  *ebiten.Image
}

type Point struct {
	Name     string `json:"name"`
	X        uint   `json:"x"`
	Y        uint   `json:"y"`
	Category string `json:"category"`
}

// New - コンストラクタ
func New() *TownMap {
	file, err := ioutil.ReadFile("asset/townmap/townmap.json")
	if err != nil {
		panic(err)
	}

	tm := new(TownMap)
	if err := json.Unmarshal(file, &tm.Regions); err != nil {
		panic(err)
	}

	tm.initMap()

	return tm
}

func (tm *TownMap) initMap() {
	for name, region := range tm.Regions {
		filename := fmt.Sprintf("asset/townmap/%s.png", name)
		background, _, err := ebitenutil.NewImageFromFile(filename, ebiten.FilterDefault)
		if err != nil {
			panic(err)
		}

		for _, point := range region.Points {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(point.X*16), float64(point.Y*16))

			switch point.Category {
			case "town":
				background.DrawImage(townImage, op)
			case "point":
				background.DrawImage(pointImage, op)
			}
		}

		tmImage, _, err := ebitenutil.NewImageFromFile("asset/townmap/asset/frame.png", ebiten.FilterDefault)
		if err != nil {
			panic(err)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(8), float64(8))
		tmImage.DrawImage(background, op)
		region.Image = tmImage
	}
}
