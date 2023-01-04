package service

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(LaptopID, imageType string, imageData bytes.Buffer) (string, error)
}

type ImageInfo struct {
	LaptopID string
	Type     string
	Path     string
}

type DiskImageStore struct {
	mutex       sync.Mutex
	imageFolder string
	images      map[string]*ImageInfo
}

func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

func (store *DiskImageStore) Save(LaptopID, imageType string, imageData bytes.Buffer) (string, error) {

	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image ID: %w", err)
	}

	image_path := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)
	log.Print("image path: ", image_path)
	file, err := os.Create(image_path)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}
	defer file.Close()

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image data to file: %w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageID.String()] = &ImageInfo{
		LaptopID: LaptopID,
		Type:     imageType,
		Path:     image_path,
	}

	return imageID.String(), nil
}
