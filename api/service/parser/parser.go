package parser

import (
    "fmt"
    "log"

    "dicom/api/common"
    "dicom/api/repository/sql"

    "github.com/suyashkumar/dicom"
)

type Parser interface {
    GetDicomDatasetByPath(dicomFilePath string) (*dicom.Dataset, string, error)
    GetDicomDatasetByFile(file string) (error)
}

type DicomParser struct {
    sql    sql.Repository
    logger *log.Logger
}

func NewDicomParser(repo sql.Repository, logger *log.Logger) *DicomParser {
    p := &DicomParser{
        sql:    repo,
        logger: logger,
    }

    return p
}

func (p *DicomParser) GetDicomDatasetByPath(dicomFilePath string) (*dicom.Dataset, string, error) {
    dataset, err := dicom.ParseFile(dicomFilePath, nil)
    if err != nil {
        p.logger.Printf("Error parsing DICOM file: %v", err)
        return nil, "", err
    }

    uuid := common.GenShortUUID()
    imageURL := fmt.Sprintf("output/image_%s.png", uuid)

    _, err = p.sql.InsertDicom(imageURL, uuid)
    if err != nil {
        p.logger.Printf("Error inserting DICOM into database: %v", err)
        return nil, "", err
    }

    return &dataset, uuid, nil
}

// To be implemented. Preferably an uploader service would be implemented as well if transferring files over https
func (p *DicomParser) GetDicomDatasetByFile(file string) (*dicom.Dataset, error) {
    return nil, nil
}
