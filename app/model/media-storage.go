package model

import (
	"errors"
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
	switch mediaType {
	case "gif":
		for url := range ms.gifs {
			return url, nil
		}
	case "image":
		for url := range ms.images {
			return url, nil
		}
	case "video":
		for url := range ms.videos {
			return url, nil
		}
	}
	return "", errors.New("media type not found")
}
