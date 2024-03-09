package fetcher

import (
    "image"
    "log"

    "dicom/api/model"
    "dicom/api/repository/blob"
    "dicom/api/repository/sql"
)

type Fetcher interface {
    GetImage(uuid string) (image.Image, error)
}

type DicomFetcher struct {
    blobStorage blob.Repository
    sql         sql.Repository
    logger      *log.Logger
}

func NewDicomFetcher(blobStorage blob.Repository, sql sql.Repository, logger *log.Logger) *DicomFetcher {
    return &DicomFetcher{
        blobStorage: blobStorage,
        sql:          sql,
        logger:       logger,
    }
}

func (d *DicomFetcher) GetImage(uuid string) (image.Image, error) {
    dicom, err := d.sql.GetDicomByUUID(uuid)
    if err != nil {
        d.logger.Printf("Error retrieving DICOM by UUID: %v", err)
        return nil, err
    }

    img, err := d.blobStorage.ReadImageFromFile(dicom.ImageURL)
    if err != nil {
        d.logger.Printf("Error reading image from file: %v", err)
        return nil, err
    }

    return img, nil
}

func (d *DicomFetcher) GetTags(uuid string) ([]model.Tag, error) {
    tags, err := d.sql.GetTagsByDicomUUID(uuid)
    if err != nil {
        d.logger.Printf("Error retrieving tags by DICOM UUID: %v", err)
        return nil, err
    }

    return tags, nil
}