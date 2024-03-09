package parser

import (
    "fmt"

    "dicom/api/common"
    "dicom/api/repository/sql"

    "github.com/suyashkumar/dicom"
    // "github.com/suyashkumar/dicom/pkg/tag"
)


type Parser interface {
    GetDicomDatasetByPath(dicomFilePath string) (error)
    GetDicomDatasetByFile(file string) (error)
}

type DicomParser struct {
    sql     sql.Repository
}

func NewDicomConverter(repo sql.Repository) *DicomParser {
    d := &DicomParser{
        sql: repo,
    }

    return d
}

func (d DicomParser) GetDicomDatasetByPath(dicomFilePath string) (*dicom.Dataset, string, error) {
    dataset, err := dicom.ParseFile(dicomFilePath, nil)
    if err != nil {
    	return nil, "", err
    }

    uuid := common.GenShortUUID()
    image_url := fmt.Sprintf("output/image_%s.png", uuid)

    id, err := sql.InsertDicom(image_url, uuid)
    if err != nil {
    	return nil, "", err
    }

	return dataset, uuid, nil
}

// To be implemented. Preferably an uploader service would be implemented as well if transferring files over https
func (d DicomParser) GetDicomDatasetByFile(file string) (*dicom.Dataset, error) {
	return nil
}

// func (d DicomParser) GetDicomDataset() *dicom.Dataset {
// 	return d.dataset
// }