package blob

import (
    "image"
    "image/png"
    "log"
    "os"
)

type Repository interface {
    WritePngToFile(image image.Image, path string) error
    ReadImageFromFile(path string) (image.Image, error)
}

type BlobStorage struct {
    logger *log.Logger
}

// NewBlobStorage creates a new instance of BlobStorage with the specified logger.
func NewBlobStorage(logger *log.Logger) (*BlobStorage, error) {
    return &BlobStorage{
        logger: logger,
    }, nil
}

func (b *BlobStorage) WritePngToFile(image image.Image, path string) error {
    // Create the output file
    outputFile, err := os.Create(path)
    if err != nil {
        b.logger.Printf("Error creating output file: %v", err)
        return err
    }
    defer outputFile.Close()

    err = png.Encode(outputFile, image)
    if err != nil {
        b.logger.Printf("Error encoding image to PNG: %v", err)
        return err
    }

    return nil
}

func (b *BlobStorage) ReadImageFromFile(path string) (image.Image, error) {
    // Open the image file
    file, err := os.Open(path)
    if err != nil {
        b.logger.Printf("Error opening image file: %v", err)
        return nil, err
    }
    defer file.Close()

    // Decode the image file
    img, _, err := image.Decode(file)
    if err != nil {
        b.logger.Printf("Error decoding image file: %v", err)
        return nil, err
    }

    return img, nil
}
