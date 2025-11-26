# Avagen

**Avagen** is a Go library to generate deterministic, symmetrical avatar images from strings. Perfect for profile pics or identicons.

## Features
- Symmetrical grid avatars
- Deterministic based on input string
- Configurable size and grid
- Pure Go, minimal dependencies (`xxh3`)

## Installation
```bash
go get github.com/zeebo/xxh3

Usage

avatarBytes, err := GenerateAvatar("Aurel", 256, 6)
if err != nil { panic(err) }
os.WriteFile("avatar.png", avatarBytes, 0644)

name: input string

size: avatar size in pixels (default 256, max 1080)

grid: number of blocks in the grid (default 6)
