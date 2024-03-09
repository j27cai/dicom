package blob

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"log"
	"os"
	"testing"
)

var (
	testImage      image.Image
	testImagePath  = "test_image.png"
	mockLogger     *MockLogger
	blockStorage   *BlobStorage
)

type MockLogger struct {
	Output *bytes.Buffer
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		Output: &bytes.Buffer{},
	}
}

// Printf mocks the Printf method of log.Logger
func (l *MockLogger) Printf(format string, v ...interface{}) {
	l.Output.WriteString(format)
}

func init() {
	// Create a test image
	testImage = image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Initialize the mock logger
	mockLogger = NewMockLogger()

	// Create a new BlobStorage instance
	blockStorage, _ = NewBlobStorage(log.New(mockLogger.Output, "", 0))
}

func TestWritePngToFile(t *testing.T) {
	err := blockStorage.WritePngToFile(testImage, testImagePath)
	if err != nil {
		t.Errorf("WritePngToFile returned an unexpected error: %v", err)
	}

	if _, err := os.Stat(testImagePath); os.IsNotExist(err) {
		t.Errorf("WritePngToFile failed to create the output file")
	}

	// Remove the temporary file
	err = os.Remove(testImagePath)
	if err != nil {
		t.Errorf("Failed to remove temporary file: %v", err)
	}
}

func TestReadImageFromFile(t *testing.T) {
	// Encode test image to PNG and write it to a file
	file, _ := os.Create(testImagePath)
	defer file.Close()
	err := png.Encode(file, testImage)
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	// Read the image file using BlobStorage
	imageFromStorage, err := blockStorage.ReadImageFromFile(testImagePath)
	if err != nil {
		t.Errorf("ReadImageFromFile returned an unexpected error: %v", err)
	}

	// Check if the dimensions of the retrieved image match the original image
	if imageFromStorage.Bounds().Size() != testImage.Bounds().Size() {
		t.Errorf("Dimensions of the retrieved image do not match the original image")
	}

	// Remove the temporary file
	err = os.Remove(testImagePath)
	if err != nil {
		t.Errorf("Failed to remove temporary file: %v", err)
	}
}

func TestWritePngToFile_Error(t *testing.T) {
	err := blockStorage.WritePngToFile(testImage, "/invalid_directory/test_image.png")
	if err == nil {
		t.Error("Expected an error when writing to a non-existent directory, but got nil")
	}

	expectedLogOutput := "Error creating output file:"
	if !bytes.Contains(mockLogger.Output.Bytes(), []byte(expectedLogOutput)) {
		t.Errorf("Expected logger output containing '%s', got '%s'", expectedLogOutput, mockLogger.Output.String())
	}
}

func TestReadImageFromFile_Error(t *testing.T) {
	_, err := blockStorage.ReadImageFromFile("/non_existent_file.png")
	if !errors.Is(err, os.ErrNotExist) {
		t.Error("Expected error due to non-existent file, but got a different error")
	}

	expectedLogOutput := "Error opening image file:"
	if !bytes.Contains(mockLogger.Output.Bytes(), []byte(expectedLogOutput)) {
		t.Errorf("Expected logger output containing '%s', got '%s'", expectedLogOutput, mockLogger.Output.String())
	}
}
