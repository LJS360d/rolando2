package model

import (
	"errors"
	"math/rand"
	"rolando/app/utils"
	"sync"
)

// MediaStorage stores media URLs
type MediaStorage struct {
	chainID string
	gifs    map[string]struct{}
	images  map[string]struct{}
	videos  map[string]struct{}
	mu      sync.RWMutex
}

func NewMediaStorage(chainID string, gifs, images, videos []string) *MediaStorage {
	storage := &MediaStorage{
		chainID: chainID,
		gifs:    make(map[string]struct{}),
		images:  make(map[string]struct{}),
		videos:  make(map[string]struct{}),
	}
	for _, gif := range gifs {
		storage.gifs[gif] = struct{}{}
	}
	for _, image := range images {
		storage.images[image] = struct{}{}
	}
	for _, video := range videos {
		storage.videos[video] = struct{}{}
	}
	return storage
}

func (ms *MediaStorage) AddMedia(url string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if utils.IsGif(url) {
		ms.gifs[url] = struct{}{}
	} else if utils.IsVideo(url) {
		ms.videos[url] = struct{}{}
	} else if utils.IsImage(url) {
		ms.images[url] = struct{}{}
	}
}

func (ms *MediaStorage) RemoveMedia(url string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.gifs, url)
	delete(ms.videos, url)
	delete(ms.images, url)
}

func (ms *MediaStorage) GetMedia(mediaType string) (string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var keys []string
	// Extract keys based on the media type
	switch mediaType {
	case "gif":
		for url := range ms.gifs {
			keys = append(keys, url)
		}
	case "image":
		for url := range ms.images {
			keys = append(keys, url)
		}
	case "video":
		for url := range ms.videos {
			keys = append(keys, url)
		}
	default:
		return "", errors.New("media type not found")
	}

	// If no media found, return an error
	if len(keys) == 0 {
		return "", errors.New("no media found for this type")
	}

	// Seed the random number generator and pick a random key
	randomIndex := rand.Intn(len(keys))
	return keys[randomIndex], nil
}
