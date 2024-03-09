package block

import (
    "image"
    "image/png"
    "os"
)

type Repository interface {
    WritePngToFile(image image.Image, path string) error
    ReadImageFromFile(path string) (image.Image, error)
}

type BlockStorage struct {
	path string
}

func Setup() (*BlockStorage, error) {
	return &BlockStorage{}, nil
}

func (b *BlockStorage) WritePngToFile(image image.Image, path string) error {
    // Create the output file
    outputFile, err := os.Create(path)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    err = png.Encode(outputFile, image)
    if err != nil {
        return err
    }

    return nil
}

func (b *BlockStorage) ReadImageFromFile(path string) (image.Image, error) {
    // Open the image file
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Decode the image file
    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }

    return img, nil
}
