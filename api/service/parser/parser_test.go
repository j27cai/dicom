package parser

import (
    "errors"
    "log"
    "testing"

    "dicom/api/repository/sql"
)

func TestDicomParser_GetDicomDatasetByPath_Success(t *testing.T) {
    mockSQLRepo := &sql.MockRepository{
        InsertDicomFunc: func(imageURL string, uuid string) (int64, error) {
            return 1, nil
        },
    }

    parser := NewDicomParser(mockSQLRepo, log.Default())

    dataset, uuid, err := parser.GetDicomDatasetByPath("test_file.dcm")
    if err != nil {
        t.Errorf("Unexpected error: %v", err)
    }
    if dataset == nil {
        t.Error("Expected dataset, got nil")
    }
    if uuid == "" {
        t.Error("Expected uuid, got ''")
    }
}

func TestDicomParser_GetDicomDatasetByPath_SQL_Error(t *testing.T) {
    mockSQLRepo := &sql.MockRepository{
        InsertDicomFunc: func(imageURL string, uuid string) (int64, error) {
            return 0, errors.New("SQL error")
        },
    }

    parser := NewDicomParser(mockSQLRepo, log.Default())

    _, _, err := parser.GetDicomDatasetByPath("test_file.dcm")
    if err == nil {
        t.Error("Expected error, got nil")
    }
}
