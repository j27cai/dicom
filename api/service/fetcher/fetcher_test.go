package fetcher

import (
    "errors"
    "image"
    "log"
    "testing"

    "dicom/api/model"
    "dicom/api/repository/blob"
    "dicom/api/repository/sql"
)

func TestDicomFetcher_GetImage_Success(t *testing.T) {
    mockImage := image.NewRGBA(image.Rect(0, 0, 100, 100))
    mockUUID := "mock_uuid"

    // Mock repositories
    mockBlobRepo := &blob.MockRepository{
        ReadImageFromFileFunc: func(path string) (image.Image, error) {
            return mockImage, nil
        },
    }
    mockSQLRepo := &sql.MockRepository{
        GetDicomByUUIDFunc: func(uuid string) (*model.Dicom, error) {
            if uuid == mockUUID {
                return &model.Dicom{ImageURL: "test_image_url"}, nil
            }
            return nil, errors.New("not found")
        },
    }

    // Create the DicomFetcher with mock repositories
    fetcher := NewDicomFetcher(mockBlobRepo, mockSQLRepo, log.Default())

    // Test GetImage method
    imageData, err := fetcher.GetImage(mockUUID)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if imageData == nil {
        t.Error("Expected image data, got nil")
    }
}

func TestDicomFetcher_GetTags_Success(t *testing.T) {
    mockUUID := "mock_uuid"

    // Mock repository
    mockSQLRepo := &sql.MockRepository{
        GetTagsByDicomUUIDFunc: func(uuid string) ([]model.Tag, error) {
            if uuid == mockUUID {
                return []model.Tag{{ID: 1, Tag: "Tag1"}, {ID: 2, Tag: "Tag2"}}, nil
            }
            return nil, errors.New("not found")
        },
    }

    // Create the DicomFetcher with mock repository
    fetcher := NewDicomFetcher(nil, mockSQLRepo, log.Default())

    // Test GetTags method
    tags, err := fetcher.GetTags(mockUUID)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if len(tags) != 2 {
        t.Errorf("Expected 2 tags, got %d", len(tags))
    }
    if tags[0].Tag != "Tag1" || tags[1].Tag != "Tag2" {
        t.Errorf("Unexpected tags: %+v", tags)
    }
}
