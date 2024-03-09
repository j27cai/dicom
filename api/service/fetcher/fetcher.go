package fetcher

import (
    "image"
    "dicom/api/repository/block"
    "dicom/api/repository/sql"
)

type Fetcher interface {
    GetImage(uuid string) (image.Image, error)
}

type DicomFetcher struct {
    blockStorage block.Repository
    sql          sql.Repository
}

func NewDicomFetcher(blockStorage block.Repository, sql sql.Repository) *DicomFetcher {
    return &DicomFetcher{
        blockStorage: blockStorage,
        sql:          sql,
    }
}

func (d *DicomFetcher) GetImage(uuid string) (*image.Image, error) {
    dicom, err := d.sql.GetDicomByUUID(uuid)
    if err != nil {
        return nil, err
    }

    imageUrl, err := d.blockStorage.ReadImageFromFile(dicom.ImageURL)
    if err != nil {
    	return nil, err
    }

    return &imageUrl, nil
}