package processor

import (
    "dicom/api/repository/block"
    "dicom/api/repository/sql"

    "github.com/suyashkumar/dicom"
    "github.com/suyashkumar/dicom/pkg/tag"
)

type DicomProcessor interface {
    ExtractDicomHeaders(dicomDataset *dicom.Dataset) error
    ExtractDicomImage(id string, dicomDataset *dicom.Dataset) error
}

type DicomExtractor struct {
    sql     sql.Repository
    block   block.Repository
}

func NewDicomExtractor(sqlRepo sql.Repository, blockRepo block.Repository) *DicomExtractor {
    d := &DicomExtractor{
        sql:   sqlRepo,
        block: blockRepo,
    }

    return d
}


// Get basically all the headers and read them and store them in sql
func (d DicomProcessor) ExtractDicomHeaders(dicomDataset *dicom.Dataset) error {
    // pixelDataElement, _ := dicomDataset.FindElementByTag(tag.PixelData)
    // pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
    // for i, fr := range pixelDataInfo.Frames {
    //     img, _ := fr.GetImage() // The Go image.Image for this frame
    //     f, _ := os.Create(fmt.Sprintf("image_%d.png", i))
    //     _ = png.Encode(f, img)
    //     _ = f.Close()
    // }

    return nil
}

// Store image in "block" storage and then store then generated the image url and store the image url in database
func (d DicomProcessor) ExtractDicomImage(id string, dicomDataset *dicom.Dataset) error {
	pixelDataElement, _ := dicomDataset.FindElementByTag(tag.PixelData)
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for i, fr := range pixelDataInfo.Frames {
		img, _ := fr.GetImage()
        dicom, err := sql.GetDicomByUUID(id)
        if err != nil {
            return err
        }
		block.WritePngToFile(img, dicom.ImageURL)
	}

	return nil
}