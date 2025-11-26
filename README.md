# Avagen

**Avagen** is a Go library to generate deterministic, symmetrical avatar images from strings. It's ideal for profile pictures or identicons.

## Features

- Symmetrical grid avatars
- Deterministic output based on input string
- Configurable size and grid dimensions
- Pure Go, minimal dependencies (uses [`xxh3`](https://github.com/zeebo/xxh3))

## Installation

```bash
go get github.com/simonov41coder/avagen
go get github.com/zeebo/xxh3
```

## Usage

```go
avatarBytes, err := GenerateAvatar("Aurel", 256, 6)
if err != nil {
    panic(err)
}
os.WriteFile("avatar.png", avatarBytes, 0644)
```

**Parameters:**

- `name` (string): Input string used to deterministically generate the avatar.
- `size` (int): Avatar size in pixels (default: 256, max: 1080).
- `grid` (int): Number of blocks in the avatar grid (default: 6).

---

Enjoy creating unique, deterministic avatars for your projects!
