package blob

import (
    "image"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
    WritePngToFileFunc  func(image image.Image, path string) error
    ReadImageFromFileFunc func(path string) (image.Image, error)
}

func (m *MockRepository) WritePngToFile(image image.Image, path string) error {
    if m.WritePngToFileFunc != nil {
        return m.WritePngToFileFunc(image, path)
    }
    return nil
}

func (m *MockRepository) ReadImageFromFile(path string) (image.Image, error) {
    if m.ReadImageFromFileFunc != nil {
        return m.ReadImageFromFileFunc(path)
    }
    return nil, nil
}
