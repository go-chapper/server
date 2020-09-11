// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package avatar is reponsible for creating avatar images
// Thanks to https://github.com/taironas/tinygraphs
package avatar

import (
	"image/color"
)

var pallete = map[string][]color.RGBA{
	"blue": {
		{3, 63, 99, 255},
		{40, 102, 110, 255},
		{80, 197, 183, 255},
		{255, 251, 255, 255},
	},
	"orange": {
		{250, 121, 33, 255},
		{254, 153, 32, 255},
		{252, 252, 98, 255},
		{255, 251, 255, 255},
	},
	"green": {
		{0, 127, 95, 255},
		{43, 147, 72, 255},
		{85, 166, 48, 255},
		{255, 251, 255, 255},
	},
}

// Palettes returns the names of available color palettes
func Palettes() []string {
	return []string{"blue", "orange", "green"}
}

// Palette returns a color palette with the provided 'paletteName'
func Palette(paletteName string) ([]color.RGBA, bool) {
	p, ok := pallete[paletteName]
	return p, ok
}
