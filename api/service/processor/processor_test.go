package processor

import (
    "errors"
    "image"
    "log"
    "testing"

    "dicom/api/model"
    "dicom/api/repository/blob"
    "dicom/api/repository/sql"
    "dicom/api/service/parser"

    "github.com/suyashkumar/dicom"
)

func TestDicomProcessor_ExtractDicomHeaders_Success(t *testing.T) {
    // Mock dataset and UUID
    mockDataset := &dicom.Dataset{}
    mockUUID := "mock_uuid"

    // Mock SQL repository
    mockSQLRepo := &sql.MockRepository{
        GetDicomByUUIDFunc: func(uuid string) (*model.Dicom, error) {
            return &model.Dicom{ID: 1}, nil
        },
        InsertTagFunc: func(tag model.Tag) (int64, error) {
            return 1, nil
        },
        InsertDicomTagFunc: func(dicomID, tagID int64) (int64, error) {
            return 1, nil
        },
    }

    // Create the DicomProcessor with mock repository
    processor := NewDicomProcessor(mockSQLRepo, nil, log.Default())

    // Test ExtractDicomHeaders method
    err := processor.ExtractDicomHeaders(mockUUID, mockDataset)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
}

func TestDicomProcessor_ExtractDicomHeaders_SQL_Error(t *testing.T) {
    // Mock SQL repository with error
    mockSQLRepo := &sql.MockRepository{
        GetDicomByUUIDFunc: func(uuid string) (*model.Dicom, error) {
            return nil, errors.New("SQL error")
        },
    }

    // Create the DicomProcessor with mock repository
    processor := NewDicomProcessor(mockSQLRepo, nil, log.Default())

    // Test ExtractDicomHeaders method with SQL error
    err := processor.ExtractDicomHeaders("mock_uuid", &dicom.Dataset{})
    if err == nil {
        t.Error("Expected error, got nil")
    }
}

func TestDicomProcessor_ExtractDicomImage_Success(t *testing.T) {
    mockSQLRepo := &sql.MockRepository{
        GetDicomByUUIDFunc: func(uuid string) (*model.Dicom, error) {
            return &model.Dicom{ID: 1, ImageURL: "output/image_mock_uuid.png"}, nil
        },
    }

    parser :=  parser.NewDicomParser(mockSQLRepo, log.Default()) 

    mockBlobRepo := &blob.MockRepository{
        WritePngToFileFunc: func(image.Image, string) error {
            return nil
        },
    }

    processor := NewDicomProcessor(mockSQLRepo, mockBlobRepo, log.Default())

    dataset, _, _ := parser.GetDicomDatasetByPath("test_file.dcm")

    // Test ExtractDicomImage method
    err := processor.ExtractDicomImage("mock_uuid", dataset)
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
}

func TestDicomProcessor_ExtractDicomImage_SQL_Error(t *testing.T) {
    mockSQLRepo := &sql.MockRepository{
        GetDicomByUUIDFunc: func(uuid string) (*model.Dicom, error) {
            return nil, errors.New("SQL error")
        },
    }

    parser :=  parser.NewDicomParser(mockSQLRepo, log.Default()) 
    processor := NewDicomProcessor(mockSQLRepo, nil, log.Default())

    dataset, _, _ := parser.GetDicomDatasetByPath("test_file.dcm")

    // Test ExtractDicomImage method with SQL error
    err := processor.ExtractDicomImage("mock_uuid", dataset)
    if err == nil {
        t.Error("Expected error, got nil")
    }
}
