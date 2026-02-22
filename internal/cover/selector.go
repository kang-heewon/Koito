package cover

import (
	"strconv"
	"strings"
)

type Image struct {
	URL        string
	Front      bool
	Back       bool
	Width      int
	Height     int
	Thumbnails map[string]string
}

type candidate struct {
	url      string
	width    int
	height   int
	priority int
}

func CoverImageExtract(images []Image) string {
	best := candidate{}
	hasBest := false

	for _, image := range images {
		c, ok := imageCandidate(image)
		if !ok {
			continue
		}
		if !hasBest || betterCandidate(c, best) {
			best = c
			hasBest = true
		}
	}

	if !hasBest {
		return ""
	}
	return best.url
}

func imageCandidate(image Image) (candidate, bool) {
	if strings.TrimSpace(image.URL) == "" && len(image.Thumbnails) == 0 {
		return candidate{}, false
	}

	url := strings.TrimSpace(image.URL)
	width := image.Width
	height := image.Height

	thumbURL, thumbWidth, thumbHeight, hasThumb := bestThumbnail(image.Thumbnails)
	thumbArea := thumbWidth * thumbHeight
	originalArea := width * height
	if hasThumb && thumbArea > originalArea {
		url = thumbURL
		width = thumbWidth
		height = thumbHeight
	}

	if url == "" {
		return candidate{}, false
	}

	return candidate{
		url:      url,
		width:    width,
		height:   height,
		priority: imagePriority(image),
	}, true
}

func imagePriority(image Image) int {
	if image.Front {
		return 2
	}
	if image.Back {
		return 0
	}
	return 1
}

func betterCandidate(next, current candidate) bool {
	if next.priority != current.priority {
		return next.priority > current.priority
	}

	nextArea := next.width * next.height
	currentArea := current.width * current.height
	if nextArea != currentArea {
		return nextArea > currentArea
	}

	if next.width != current.width {
		return next.width > current.width
	}

	if next.height != current.height {
		return next.height > current.height
	}

	return false
}

func bestThumbnail(thumbnails map[string]string) (string, int, int, bool) {
	bestURL := ""
	bestSize := 0

	for key, url := range thumbnails {
		if strings.TrimSpace(url) == "" {
			continue
		}
		size := thumbnailSize(key)
		if size == 0 {
			continue
		}
		if size > bestSize {
			bestSize = size
			bestURL = url
		}
	}

	if bestURL == "" {
		return "", 0, 0, false
	}

	return bestURL, bestSize, bestSize, true
}

func thumbnailSize(key string) int {
	normalized := strings.TrimSpace(strings.ToLower(key))
	switch normalized {
	case "small":
		return 250
	case "large":
		return 500
	}

	size, err := strconv.Atoi(normalized)
	if err != nil || size <= 0 {
		return 0
	}
	return size
}
