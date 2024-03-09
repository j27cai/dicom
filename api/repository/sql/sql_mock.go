package sql

import (
    "dicom/api/model"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
    CloseFunc          func() error
    InsertDicomFunc    func(imageURL string, uuid string) (int64, error)
    InsertTagFunc      func(tag model.Tag) (int64, error)
    InsertDicomTagFunc func(dicomID, tagID int64) (int64, error)
    GetDicomByUUIDFunc func(uuid string) (*model.Dicom, error)
    GetTagsByDicomUUIDFunc func(uuid string) ([]model.Tag, error)
}

func (m *MockRepository) Close() error {
    if m.CloseFunc != nil {
        return m.CloseFunc()
    }
    return nil
}

func (m *MockRepository) InsertDicom(imageURL string, uuid string) (int64, error) {
    if m.InsertDicomFunc != nil {
        return m.InsertDicomFunc(imageURL, uuid)
    }
    return 0, nil
}

func (m *MockRepository) InsertTag(tag model.Tag) (int64, error) {
    if m.InsertTagFunc != nil {
        return m.InsertTagFunc(tag)
    }
    return 0, nil
}

func (m *MockRepository) InsertDicomTag(dicomID, tagID int64) (int64, error) {
    if m.InsertDicomTagFunc != nil {
        return m.InsertDicomTagFunc(dicomID, tagID)
    }
    return 0, nil
}

func (m *MockRepository) GetDicomByUUID(uuid string) (*model.Dicom, error) {
    if m.GetDicomByUUIDFunc != nil {
        return m.GetDicomByUUIDFunc(uuid)
    }
    return nil, nil
}

func (m *MockRepository) GetTagsByDicomUUID(uuid string) ([]model.Tag, error) {
    if m.GetTagsByDicomUUIDFunc != nil {
        return m.GetTagsByDicomUUIDFunc(uuid)
    }
    return nil, nil
}
