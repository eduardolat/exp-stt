// Package shortcut provides global hotkey management.
package shortcut

import (
	"errors"
	"strings"

	"github.com/varavelio/tribar/internal/config"
	"golang.design/x/hotkey"
)

// ErrInvalidKey is returned when an unknown key string is provided.
var ErrInvalidKey = errors.New("invalid key")

// ErrInvalidModifier is returned when an unknown modifier string is provided.
var ErrInvalidModifier = errors.New("invalid modifier")

// keyMap maps string key names to hotkey.Key constants.
var keyMap = map[string]hotkey.Key{
	"space":  hotkey.KeySpace,
	"a":      hotkey.KeyA,
	"b":      hotkey.KeyB,
	"c":      hotkey.KeyC,
	"d":      hotkey.KeyD,
	"e":      hotkey.KeyE,
	"f":      hotkey.KeyF,
	"g":      hotkey.KeyG,
	"h":      hotkey.KeyH,
	"i":      hotkey.KeyI,
	"j":      hotkey.KeyJ,
	"k":      hotkey.KeyK,
	"l":      hotkey.KeyL,
	"m":      hotkey.KeyM,
	"n":      hotkey.KeyN,
	"o":      hotkey.KeyO,
	"p":      hotkey.KeyP,
	"q":      hotkey.KeyQ,
	"r":      hotkey.KeyR,
	"s":      hotkey.KeyS,
	"t":      hotkey.KeyT,
	"u":      hotkey.KeyU,
	"v":      hotkey.KeyV,
	"w":      hotkey.KeyW,
	"x":      hotkey.KeyX,
	"y":      hotkey.KeyY,
	"z":      hotkey.KeyZ,
	"0":      hotkey.Key0,
	"1":      hotkey.Key1,
	"2":      hotkey.Key2,
	"3":      hotkey.Key3,
	"4":      hotkey.Key4,
	"5":      hotkey.Key5,
	"6":      hotkey.Key6,
	"7":      hotkey.Key7,
	"8":      hotkey.Key8,
	"9":      hotkey.Key9,
	"return": hotkey.KeyReturn,
	"escape": hotkey.KeyEscape,
	"delete": hotkey.KeyDelete,
	"tab":    hotkey.KeyTab,
	"left":   hotkey.KeyLeft,
	"right":  hotkey.KeyRight,
	"up":     hotkey.KeyUp,
	"down":   hotkey.KeyDown,
	"f1":     hotkey.KeyF1,
	"f2":     hotkey.KeyF2,
	"f3":     hotkey.KeyF3,
	"f4":     hotkey.KeyF4,
	"f5":     hotkey.KeyF5,
	"f6":     hotkey.KeyF6,
	"f7":     hotkey.KeyF7,
	"f8":     hotkey.KeyF8,
	"f9":     hotkey.KeyF9,
	"f10":    hotkey.KeyF10,
	"f11":    hotkey.KeyF11,
	"f12":    hotkey.KeyF12,
}

// modMap is defined in platform-specific files:
// - modifiers_linux.go
// - modifiers_windows.go
// - modifiers_darwin.go

// ParseShortcut converts a config.Shortcut to hotkey modifiers and key.
func ParseShortcut(s config.Shortcut) ([]hotkey.Modifier, hotkey.Key, error) {
	mods := make([]hotkey.Modifier, 0, len(s.Modifiers))
	for _, m := range s.Modifiers {
		mod, ok := modMap[strings.ToLower(m)]
		if !ok {
			return nil, 0, ErrInvalidModifier
		}
		mods = append(mods, mod)
	}

	key, ok := keyMap[strings.ToLower(s.Key)]
	if !ok {
		return nil, 0, ErrInvalidKey
	}

	return mods, key, nil
}
