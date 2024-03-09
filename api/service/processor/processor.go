package processor

import (
    "log"

    "dicom/api/model"
    "dicom/api/repository/blob"
    "dicom/api/repository/sql"

    "github.com/suyashkumar/dicom"
    "github.com/suyashkumar/dicom/pkg/tag"
)

type Processor interface {
    ExtractDicomHeaders(id string, dicomDataset *dicom.Dataset) error
    ExtractDicomImage(id string, dicomDataset *dicom.Dataset) error
}

type DicomProcessor struct {
    sql     sql.Repository
    blob   blob.Repository
    logger  *log.Logger
}

func NewDicomProcessor(sqlRepo sql.Repository, blobRepo blob.Repository, logger *log.Logger) *DicomProcessor {
    p := &DicomProcessor{
        sql:    sqlRepo,
        blob:  blobRepo,
        logger: logger,
    }

    return p
}

// Get all the headers and read them and store them in SQL
func (p *DicomProcessor) ExtractDicomHeaders(id string, dicomDataset *dicom.Dataset) error {
    dicom, err := p.sql.GetDicomByUUID(id)
    if err != nil {
        p.logger.Printf("Error retrieving DICOM by UUID: %v", err)
        return err
    }

    for elem := dicomDataset.FlatStatefulIterator(); elem.HasNext(); {
        e := elem.Next()
        var tagName string
        if tagInfo, err := tag.Find(e.Tag); err == nil {
            tagName = tagInfo.Name
        }
        
        tag := model.Tag{
            Tag:          e.Tag.String(),
            Name:         tagName,
            VR:           e.ValueRepresentation.String(),
            Value:        e.Value.String(),
        }

        tagID, err := p.sql.InsertTag(tag)
        if err != nil {
            p.logger.Printf("Error inserting tag: %v", err)
            return err
        }

        _, err = p.sql.InsertDicomTag(dicom.ID, tagID)
        if err != nil {
            p.logger.Printf("Error inserting tag: %v", err)
            return err
        }
    }

    return nil
}

// Store image in "blob" storage and then store then generated the image URL and store the image URL in the database
func (p *DicomProcessor) ExtractDicomImage(id string, dicomDataset *dicom.Dataset) error {
    pixelDataElement, _ := dicomDataset.FindElementByTag(tag.PixelData)
    pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
    for _, fr := range pixelDataInfo.Frames {
        img, _ := fr.GetImage()
        dicom, err := p.sql.GetDicomByUUID(id)
        if err != nil {
            p.logger.Printf("Error retrieving DICOM by UUID: %v", err)
            return err
        }
        if err := p.blob.WritePngToFile(img, dicom.ImageURL); err != nil {
            p.logger.Printf("Error writing PNG to file: %v", err)
            return err
        }
    }

    return nil
}
