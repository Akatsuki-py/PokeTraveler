package config

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten"
)

type Config struct {
	Avatar   int
	Controls Controls
}

type Controls struct {
	Up, Down, Left, Right, A, B, Start ebiten.Key
}

func (c Config) Validate() error {
	var errs []string

	if c.Avatar != 1 || c.Avatar != 2 {
		errs = append(errs, "avatar must be set to '1' or '2'")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}

	return nil
}

func NewDefaultConfig() Config {
	return Config{
		Avatar:   1,
		Controls: NewDefaultControls(),
	}
}

func NewDefaultControls() Controls {
	up, _ := keyTranslator("up")
	down, _ := keyTranslator("down")
	left, _ := keyTranslator("left")
	right, _ := keyTranslator("right")
	a, _ := keyTranslator("s")
	b, _ := keyTranslator("a")
	start, _ := keyTranslator("enter")

	return Controls{
		Up:    up,
		Down:  down,
		Left:  left,
		Right: right,
		A:     a,
		B:     b,
		Start: start,
	}
}

func NewConfigFromFile(f string) (Config, error) {
	var cnf struct {
		Avatar   int
		Controls struct {
			Up, Down, Left, Right, A, B, Start string
		}
	}

	if _, err := toml.DecodeFile(f, &cnf); err != nil {
		return Config{}, fmt.Errorf("unable to parse config file %s: %s", f, err)
	}

	up, _ := keyTranslator(cnf.Controls.Up)
	down, _ := keyTranslator(cnf.Controls.Down)
	left, _ := keyTranslator(cnf.Controls.Left)
	right, _ := keyTranslator(cnf.Controls.Right)
	a, _ := keyTranslator(cnf.Controls.A)
	b, _ := keyTranslator(cnf.Controls.B)
	start, _ := keyTranslator(cnf.Controls.Start)

	return Config{
		Avatar: cnf.Avatar,
		Controls: Controls{
			Up:    up,
			Down:  down,
			Left:  left,
			Right: right,
			A:     a,
			B:     b,
			Start: start,
		},
	}, nil
}

func keyTranslator(key string) (ebiten.Key, error) {
	keyMap := map[string]ebiten.Key{
		"0":            ebiten.Key0,
		"1":            ebiten.Key1,
		"2":            ebiten.Key2,
		"3":            ebiten.Key3,
		"4":            ebiten.Key4,
		"5":            ebiten.Key5,
		"6":            ebiten.Key6,
		"7":            ebiten.Key7,
		"8":            ebiten.Key8,
		"9":            ebiten.Key9,
		"a":            ebiten.KeyA,
		"b":            ebiten.KeyB,
		"c":            ebiten.KeyC,
		"d":            ebiten.KeyD,
		"e":            ebiten.KeyE,
		"f":            ebiten.KeyF,
		"g":            ebiten.KeyG,
		"h":            ebiten.KeyH,
		"i":            ebiten.KeyI,
		"j":            ebiten.KeyJ,
		"k":            ebiten.KeyK,
		"l":            ebiten.KeyL,
		"m":            ebiten.KeyM,
		"n":            ebiten.KeyN,
		"o":            ebiten.KeyO,
		"p":            ebiten.KeyP,
		"q":            ebiten.KeyQ,
		"r":            ebiten.KeyR,
		"s":            ebiten.KeyS,
		"t":            ebiten.KeyT,
		"u":            ebiten.KeyU,
		"v":            ebiten.KeyV,
		"w":            ebiten.KeyW,
		"x":            ebiten.KeyX,
		"y":            ebiten.KeyY,
		"z":            ebiten.KeyZ,
		"apostrophe":   ebiten.KeyApostrophe,
		"backslash":    ebiten.KeyBackslash,
		"backspace":    ebiten.KeyBackspace,
		"capslock":     ebiten.KeyCapsLock,
		"comma":        ebiten.KeyComma,
		"delete":       ebiten.KeyDelete,
		"down":         ebiten.KeyDown,
		"end":          ebiten.KeyEnd,
		"enter":        ebiten.KeyEnter,
		"equal":        ebiten.KeyEqual,
		"escape":       ebiten.KeyEscape,
		"f1":           ebiten.KeyF1,
		"f2":           ebiten.KeyF2,
		"f3":           ebiten.KeyF3,
		"f4":           ebiten.KeyF4,
		"f5":           ebiten.KeyF5,
		"f6":           ebiten.KeyF6,
		"f7":           ebiten.KeyF7,
		"f8":           ebiten.KeyF8,
		"f9":           ebiten.KeyF9,
		"f10":          ebiten.KeyF10,
		"f11":          ebiten.KeyF11,
		"f12":          ebiten.KeyF12,
		"graveaccent":  ebiten.KeyGraveAccent,
		"home":         ebiten.KeyHome,
		"insert":       ebiten.KeyInsert,
		"kp0":          ebiten.KeyKP0,
		"kp1":          ebiten.KeyKP1,
		"kp2":          ebiten.KeyKP2,
		"kp3":          ebiten.KeyKP3,
		"kp4":          ebiten.KeyKP4,
		"kp5":          ebiten.KeyKP5,
		"kp6":          ebiten.KeyKP6,
		"kp7":          ebiten.KeyKP7,
		"kp8":          ebiten.KeyKP8,
		"kp9":          ebiten.KeyKP9,
		"kpadd":        ebiten.KeyKPAdd,
		"kpdecimal":    ebiten.KeyKPDecimal,
		"kpdivide":     ebiten.KeyKPDivide,
		"kpenter":      ebiten.KeyKPEnter,
		"kpequal":      ebiten.KeyKPEqual,
		"kpmultiply":   ebiten.KeyKPMultiply,
		"kpsubtract":   ebiten.KeyKPSubtract,
		"left":         ebiten.KeyLeft,
		"leftbracket":  ebiten.KeyLeftBracket,
		"menu":         ebiten.KeyMenu,
		"minus":        ebiten.KeyMinus,
		"numlock":      ebiten.KeyNumLock,
		"pagedown":     ebiten.KeyPageDown,
		"pageup":       ebiten.KeyPageUp,
		"pause":        ebiten.KeyPause,
		"period":       ebiten.KeyPeriod,
		"printscreen":  ebiten.KeyPrintScreen,
		"right":        ebiten.KeyRight,
		"rightbracket": ebiten.KeyRightBracket,
		"scrolllock":   ebiten.KeyScrollLock,
		"semicolon":    ebiten.KeySemicolon,
		"slash":        ebiten.KeySlash,
		"space":        ebiten.KeySpace,
		"tab":          ebiten.KeyTab,
		"up":           ebiten.KeyUp,
		"alt":          ebiten.KeyAlt,
		"control":      ebiten.KeyControl,
		"shift":        ebiten.KeyShift,
	}

	lwr := strings.ToLower(key)
	res := keyMap[lwr]
	if _, ok := keyMap[lwr]; !ok {
		// TODO: Work out how to return an empty ebiten.Key
		return ebiten.KeyShift, fmt.Errorf("%s is not a valid key", key)
	}

	return res, nil
}
