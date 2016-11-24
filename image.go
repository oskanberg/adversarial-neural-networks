package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// AddLabel adds a label (string) to a gray image at x,y
func AddLabel(img *image.Gray, x, y int, label string) {
	col := color.Gray{0}
	point := fixed.P(x, y)

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(label)
}

// CreateGreyImage Creates a given size Grey image with given dimensions
func CreateGreyImage(dx, dy int, data []byte) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, dx, dy))
	img.Pix = data
	return img
}

// RandomGreyData returns a byte slice of greyscale
func RandomGreyData(dx, dy int) []byte {
	data := make([]byte, dx*dy)
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			data[(dy*y)+x] = byte(rand.Intn(255))
		}
	}
	return data
}

// UniformWhiteData returns byte slice of white pixels (255)
func UniformWhiteData(dx, dy int) []byte {
	data := make([]byte, dx*dy)
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			data[(dy*y)+x] = byte(255)
		}
	}
	return data
}

// UniformBlackData returns byte slice of white pixels (0)
func UniformBlackData(dx, dy int) []byte {
	data := make([]byte, dx*dy)
	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			data[(dy*y)+x] = byte(0)
		}
	}
	return data
}

// SaveImage saves an image.Image to disk
func SaveImage(m image.Image, destination string) {
	var buf bytes.Buffer
	err := png.Encode(&buf, m)
	if err != nil {
		panic(err)
	}
	// enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	// fmt.Println("IMAGE:" + enc)

	toimg, _ := os.Create(destination)
	defer toimg.Close()

	png.Encode(toimg, m)
}

// LoadImage loads an image.Image from disk
func LoadImage(location string) image.Image {
	r, _ := os.Open(location)
	defer r.Close()
	img, _ := png.Decode(r)
	return img
}
