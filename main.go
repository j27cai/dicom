package main

import (
    "fmt"
    "os"
    "image/png"

    "github.com/suyashkumar/dicom"
    "github.com/suyashkumar/dicom/pkg/tag"
)

func main() {
	// dataset, _ := dicom.ParseFile("test-mri/ST000001/SE000001/IM000001", nil) // See also: dicom.Parse which has a generic io.Reader API.

	// // Dataset will nicely print the DICOM dataset data out of the box.
	// fmt.Println(dataset)

	convertToPNG("test-mri/ST000001/SE000001/IM000001")
}


func convertToPNG(dicomFilePath string) error {
    dataset, _ := dicom.ParseFile(dicomFilePath, nil)
	pixelDataElement, _ := dataset.FindElementByTag(tag.PixelData)
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for i, fr := range pixelDataInfo.Frames {
		img, _ := fr.GetImage() // The Go image.Image for this frame
		f, _ := os.Create(fmt.Sprintf("image_%d.png", i))
		_ = png.Encode(f, img)
		_ = f.Close()
	}

	return nil
}