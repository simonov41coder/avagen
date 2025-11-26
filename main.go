package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	
	"github.com/zeebo/xxh3"
)

const (
	defaultSize       = 256
	defaultGrid       = 6
	maxAvatarSize     = 1080
)

// hashByte computes 16bytes hash of a string and returns the byte slice
func hashByte(str string) []byte {
	hash := xxh3.Hash128([]byte(str))
	b := hash.Bytes()
	return b[:]
}

// GenerateAvatar creates a deterministic avatar image based on a name
// and returns the image data as bytes
func GenerateAvatar(name string, size, grid int) ([]byte, error) {
	// Validate and set defaults
	if size <= 0 {
		size = defaultSize
	}
	if size > maxAvatarSize {
		size = maxAvatarSize
	}
	if grid <= 0 {
		grid = defaultGrid
	}

	hashed := hashByte(name)

	// Generate color avoiding pure white and black
	r := hashed[0]
	g := hashed[1]
	b := hashed[2]

	// Ensure color is not too close to white or black
	// Clamp values to range [30, 225] to avoid extremes
	if r > 225 {
		r = 225
	} else if r < 30 {
		r = 30
	}
	if g > 225 {
		g = 225
	} else if g < 30 {
		g = 30
	}
	if b > 225 {
		b = 225
	} else if b < 30 {
		b = 30
	}

	avatarColor := color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	blockSize := float64(size) / float64(grid)
	gridHalf := (grid + 1) / 2

	// Generate the avatar pattern by iterating through half the grid
	// and mirroring blocks to create symmetry
	for y := 0; y < grid; y++ {
		for x := 0; x < gridHalf; x++ {
			hashIndex := (x + y*gridHalf) % 16
			if hashed[hashIndex]%2 == 0 {
				drawSymmetricalBlocks(img, x, y, size, grid, blockSize, avatarColor)
			}
		}
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// drawSymmetricalBlocks draws a colored block and its horizontal mirror
func drawSymmetricalBlocks(img *image.RGBA, x, y int, size int, grid int, blockSize float64, c color.Color) {
	x1 := int(float64(x) * blockSize)
	y1 := int(float64(y) * blockSize)
	x2 := int(float64(x+1) * blockSize)
	y2 := int(float64(y+1) * blockSize)

	// Draw left side
	leftRect := image.Rect(x1, y1, x2, y2)
	draw.Draw(img, leftRect, &image.Uniform{c}, image.Point{}, draw.Src)

	// Draw mirrored right side
	// Mirror position: grid - x - 1
	mirrorX := grid - x - 1
	rightX1 := int(float64(mirrorX) * blockSize)
	rightX2 := int(float64(mirrorX+1) * blockSize)
	rightRect := image.Rect(rightX1, y1, rightX2, y2)
	draw.Draw(img, rightRect, &image.Uniform{c}, image.Point{}, draw.Src)
}

// main demonstrates how to use the GenerateAvatar function
func main() {
	avatarBytes, err := GenerateAvatar("Aurel",256, 6)
	if err != nil {
		panic(err)
	}

	// Save the avatar to a file
	err = os.WriteFile("avatar.png", avatarBytes, 0644)
	if err != nil {
		panic(err)
	}

	println("Avatar generated successfully: avatar.png")
}
