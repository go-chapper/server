// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package avatar is reponsible for creating avatar images
// Thanks to https://github.com/taironas/tinygraphs
package avatar

import (
	"bytes"
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"chapper.dev/server/internal/modules/hash"
	"chapper.dev/server/internal/utils"
)

// Avatar represents an avatar image
type Avatar struct {
	Width       int
	Height      int
	Size        string
	Data        string
	RawData     string
	Image       *image.RGBA
	ImageBuffer *bytes.Buffer
}

var (
	// ErrNoSuchScheme gets returned when there is no such color scheme
	ErrNoSuchScheme = errors.New("No such scheme")
)

// New creates a new avatar struct and returns it
func New(size int, data string) *Avatar {
	return &Avatar{
		Width:       size,
		Height:      size,
		Size:        strconv.Itoa(size),
		Data:        hash.MD5(data),
		RawData:     data,
		Image:       image.NewRGBA(image.Rect(0, 0, size, size)),
		ImageBuffer: new(bytes.Buffer),
	}
}

// Generate will generate an avatar, encode and save it
func (a *Avatar) Generate(basePath string) error {
	pallete := GetRandomPalette()
	err := a.GenerateImage(pallete)
	if err != nil {
		return err
	}

	err = a.Encode()
	if err != nil {
		return err
	}

	return a.Save(basePath, a.Size, a.RawData)
}

// GenerateImage generates a new 6 x 6 quadrant avatar based on the data (hashed by MD5)
// and specified color palette
func (a *Avatar) GenerateImage(paletteName string) error {
	currentYQuadrant := 0
	quadrant := a.Width / 6
	colorMap := make(map[int]color.RGBA)
	colors, ok := Palette(paletteName)
	if !ok {
		return ErrNoSuchScheme
	}

	for y := 0; y < a.Height; y++ {
		yQuadrant := y / quadrant
		if yQuadrant != currentYQuadrant {
			colorMap = make(map[int]color.RGBA)
			currentYQuadrant = yQuadrant
		}

		for x := 0; x < a.Width; x++ {
			xQuadrant := x / quadrant
			if _, ok := colorMap[xQuadrant]; !ok {
				if xQuadrant < 3 {
					colorMap[xQuadrant] = pickColor(a.Data, colors, xQuadrant+3*yQuadrant)
				} else if xQuadrant < 6 {
					colorMap[xQuadrant] = colorMap[6-xQuadrant-1]
				} else {
					colorMap[xQuadrant] = colorMap[0]
				}
			}

			a.Image.Set(x, y, colorMap[xQuadrant])
		}
	}
	return nil
}

// Encode encodes the generated image as a JPEG image and writes it to ImageBuffer
func (a *Avatar) Encode() error {
	err := jpeg.Encode(a.ImageBuffer, a.Image, &jpeg.Options{
		Quality: 80,
	})
	if err != nil {
		return err
	}

	return nil
}

// Save saves the encoded image in ImageBuffer to disk at path base + size + name.jpg
func (a *Avatar) Save(base, size, name string) error {
	p := utils.Join(base, size)
	f := utils.Join(p, name+".jpg")
	err := os.MkdirAll(p, 0777)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f, a.ImageBuffer.Bytes(), 0777)
}

func pickColor(data string, colors []color.RGBA, index int) color.RGBA {
	l := len(colors)
	s := hex.EncodeToString([]byte{data[index]})

	if r, err := strconv.ParseInt(s, 16, 0); err == nil {
		for i := 0; i < l; i++ {
			if int(r)%l == i {
				return colors[i]
			}
		}
	} else {
		// TODO: Pass this error up the chain
		log.Printf("Error calling ParseInt(%v, 16, 0): %v\n", s, err)
	}
	return colors[0]
}
