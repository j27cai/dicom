package model

// Tag: (0400,0563)
//   Tag Name: ModifyingSystem
//   VR: VRStringList
//   VR Raw: LO
//   VL: 44
//   Value: [INTELEPACSPACS-4-6-1-P146LDSPACS-4-6-1-P146]


// Tag represents a DICOM tag
type Tag struct {
	ID    int64
    Tag   string
    VR    string
    Value string
    Name  string
}