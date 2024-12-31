package model

import (
	"errors"
	"math/rand"
	"net/http"
	"rolando/cmd/log"
	"rolando/cmd/repositories"
	"rolando/cmd/utils"
	"slices"
	"sync"
)

// MediaStorage stores media URLs
type MediaStorage struct {
	chainID      string
	gifs         map[string]struct{}
	images       map[string]struct{}
	videos       map[string]struct{}
	messagesRepo repositories.MessagesRepository
	mu           sync.RWMutex
}

func NewMediaStorage(chainID string, gifs, images, videos []string, messagesRepo repositories.MessagesRepository) *MediaStorage {
	storage := &MediaStorage{
		chainID:      chainID,
		gifs:         make(map[string]struct{}),
		images:       make(map[string]struct{}),
		videos:       make(map[string]struct{}),
		messagesRepo: messagesRepo,
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

	var urls []string
	// Extract URLs based on the media type
	switch mediaType {
	case "gif":
		for url := range ms.gifs {
			urls = append(urls, url)
		}
	case "image":
		for url := range ms.images {
			urls = append(urls, url)
		}
	case "video":
		for url := range ms.videos {
			urls = append(urls, url)
		}
	default:
		return "", errors.New("media type not found")
	}
	ms.mu.RUnlock()
	// If no media found, return an error
	if len(urls) == 0 {
		return "", errors.New("no media found for this type")
	}

	// Try to get a valid URL
	validURL, err := ms.getValidUrlFromSet(urls)
	if err != nil {
		return "", err
	}

	return validURL, nil
}

// getValidUrlFromSet tries to return a valid URL from the set of URLs
func (ms *MediaStorage) getValidUrlFromSet(urls []string) (string, error) {
	for len(urls) > 0 {
		randomIndex := rand.Intn(len(urls))
		url := urls[randomIndex]

		// Validate the URL
		valid, err := ms.validateUrl(url)
		if err != nil {
			return ms.getValidUrlFromSet(slices.Concat(urls[:randomIndex], urls[randomIndex+1:]))
		}

		if valid {
			// Return the first valid URL
			return url, nil
		}

		// If invalid, remove the URL from the set
		urls = slices.Concat(urls[:randomIndex], urls[randomIndex+1:])
	}

	return "", errors.New("no valid media URLs found")
}

// validateUrl checks if the URL is valid (HTTP HEAD request)
func (ms *MediaStorage) validateUrl(url string) (bool, error) {
	// Make a HEAD request to check if the URL is valid
	resp, err := http.Head(url)
	if err != nil {
		// On error, remove the invalid URL from the appropriate file
		err = ms.removeInvalidUrl(url)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Check if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		err = ms.removeInvalidUrl(url)
		if err != nil {
			log.Log.Errorf("Error removing invalid URL: %v", err)
			return false, err
		}
		return false, nil
	}

	return true, nil
}

// removeInvalidUrl removes the invalid URL from DB messages
func (ms *MediaStorage) removeInvalidUrl(url string) error {
	ms.RemoveMedia(url)
	return ms.messagesRepo.DeleteGuildMessagesContaining(ms.chainID, url)
}
