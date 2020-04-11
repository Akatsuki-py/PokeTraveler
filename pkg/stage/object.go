package stage

import "github.com/Akatsuki-py/PokeTraveler/pkg/object"

// GetObject - Get Object
func (stage *Stage) GetObject(x, y int) (target *object.Object) {
	for _, o := range stage.Objects {
		switch o.Direction {
		case object.Up:
			if o.X/16 == (x+15)/16 && ((o.Y-16)/16+1) == y/16 {
				target = o
			}
		case object.Down:
			if o.X/16 == (x+15)/16 && (o.Y+15)/16 == y/16 {
				target = o
			}
		case object.Right:
			if (o.X+15)/16 == x/16 && o.Y/16 == (y+15)/16 {
				target = o
			}
		case object.Left:
			if ((o.X-16)/16+1) == x/16 && o.Y/16 == (y+15)/16 {
				target = o
			}
		}

		if target != nil {
			break
		}
	}
	return target
}

func (stage *Stage) loadObjects(filename string) {
	stage.Objects = object.Load(filename)
}
