package townmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Akatsuki-py/PokeTraveler/pkg/char"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	titleImage, _, _  = ebitenutil.NewImageFromFile("asset/townmap/asset/title.png", ebiten.FilterDefault)
	cursorImage, _, _ = ebitenutil.NewImageFromFile("asset/townmap/asset/cursor.png", ebiten.FilterDefault)
	pointImage, _, _  = ebitenutil.NewImageFromFile("asset/townmap/asset/point.png", ebiten.FilterDefault)
	townImage, _, _   = ebitenutil.NewImageFromFile("asset/townmap/asset/town.png", ebiten.FilterDefault)
)

var stageToRegion map[string]string

type TownMap struct {
	Regions map[string]*Region
	Cursor  Cursor
}

type Region struct {
	Points []Point `json:"points"`
	Image  *ebiten.Image
}

type Point struct {
	Name     string `json:"name"`
	X        uint   `json:"x"` // X は タイル(16*16)単位
	Y        uint   `json:"y"` // Y は タイル(16*16)単位
	Category string `json:"category"`
}

// New - コンストラクタ
func New() *TownMap {
	initStageToRegion()

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

func initStageToRegion() {
	file, err := ioutil.ReadFile("asset/townmap/regions.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file, &stageToRegion); err != nil {
		panic(err)
	}
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

// Open - Open townmap
func (tm *TownMap) Open(stagename string, avatar *ebiten.Image) *ebiten.Image {
	regionName, ok := stageToRegion[stagename]
	if !ok {
		panic(fmt.Errorf("region is not registerd: %s", stagename))
	}

	region := tm.Regions[regionName]

	background, err := ebiten.NewImageFromImage(region.Image, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	// アバターを配置
	point, ok := getPoint(region.Points, stagename)
	if ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(point.X*16+4), float64(point.Y*16+4))
		background.DrawImage(avatar, op)
	}

	// カーソルを配置
	var cursorX, cursorY int
	if tm.Cursor.Valid {
		cursorX, cursorY = tm.Cursor.GetXY()
	} else {
		cursorX, cursorY = int(point.X*16), int(point.Y*16) // カーソルの位置が未初期化
	}
	{
		tm.Cursor.SetXY(cursorX, cursorY)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cursorX+4), float64(cursorY+4))
		background.DrawImage(cursorImage, op)
	}

	// 最後にマップ名とマップイメージを合体
	cursorTileX, cursorTileY := tm.Cursor.GetTileXY()
	titleName := ""
	if point, ok := getPointByXY(region.Points, cursorTileX, cursorTileY); ok {
		titleName = point.Name
	}
	title := getTitle(titleName)
	result, _ := ebiten.NewImage(160, 144, ebiten.FilterDefault)
	result.DrawImage(title, nil)
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(0), float64(16))
		result.DrawImage(background, op)
	}

	return result
}

func getPoint(points []Point, stagename string) (p Point, ok bool) {
	for _, point := range points {
		if stagename == point.Name {
			return point, true
		}
	}

	return Point{}, false
}

func getPointByXY(points []Point, tileX, tileY int) (p Point, ok bool) {
	for _, point := range points {
		if tileX == int(point.X) && tileY == int(point.Y) {
			return point, true
		}
	}

	return Point{}, false
}

func getTitle(stagename string) *ebiten.Image {
	title, err := ebiten.NewImageFromImage(titleImage, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}

	char.RenderString(title, stagename, 2, 2)
	return title
}
