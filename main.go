package main

import (
	"bytes"
	"crypto/md5"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strconv"
)

const (
	defaultSize       = 256
	defaultGrid       = 6
	maxAvatarSize     = 1080
	defaultAvatarName = "saitama"
)

// hashByte computes MD5 hash of a string and returns the byte slice
func hashByte(str string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hasher.Sum(nil)
}

// generateAvatar creates a deterministic avatar image based on a name
func generateAvatar(name string, size *int, grid *int) (*bytes.Buffer, error) {
	// Set defaults
	if size == nil {
		size = intPtr(defaultSize)
	}
	if grid == nil {
		grid = intPtr(defaultGrid)
	}

	hashed := hashByte(name)
	avatarColor := color.RGBA{
		R: hashed[0],
		G: hashed[1],
		B: hashed[2],
		A: 255,
	}

	img := image.NewRGBA(image.Rect(0, 0, *size, *size))
	blockSize := float64(*size) / float64(*grid)
	gridHalf := (*grid + 1) / 2

	// Generate the avatar pattern by iterating through half the grid
	// and mirroring blocks to create symmetry
	for y := 0; y < *grid; y++ {
		for x := 0; x < gridHalf; x++ {
			hashIndex := (x + y*gridHalf) % 16
			if hashed[hashIndex]%2 == 0 {
				drawSymmetricalBlocks(img, x, y, *size, *grid, blockSize, avatarColor)
			}
		}
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return nil, err
	}
	return &buffer, nil
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

// intPtr returns a pointer to an int
func intPtr(val int) *int {
	return &val
}

// parseIntParam parses a string to an int pointer, returning nil if parsing fails
func parseIntParam(value string) *int {
	if value == "" {
		return nil
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return intPtr(intVal)
}

// handleAvatar handles HTTP requests for avatar generation
func handleAvatar(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get name parameter with default
	name := query.Get("name")
	if name == "" {
		name = defaultAvatarName
	}

	// Parse size parameter
	size := parseIntParam(query.Get("resolution"))
	if size != nil && *size > maxAvatarSize {
		http.Error(w, "Requested payload exceeds maximum allowed size of 1080px", http.StatusRequestEntityTooLarge)
		return
	}

	// Parse grid parameter
	grid := parseIntParam(query.Get("grid"))

	// Generate avatar
	buf, err := generateAvatar(name, size, grid)
	if err != nil {
		http.Error(w, "Failed to generate avatar", http.StatusInternalServerError)
		log.Printf("Error generating avatar for '%s': %v", name, err)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("Failed to write image data: %v", err)
	}
}

func main() {
	http.HandleFunc("/avatar", handleAvatar)
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
