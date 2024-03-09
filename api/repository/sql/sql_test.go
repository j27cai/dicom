package sql

import (
	"log"
	"os"
	"testing"

	"github.com/google/uuid" // Import the uuid package
	"dicom/api/model"
)

var (
	testDB *Database
)

func init() {
	// Set up a test logger
	logger := log.New(os.Stdout, "test ", log.LstdFlags)

	// Initialize a test database instance
	db, err := NewSqlDatabase(logger)
	if err != nil {
		log.Fatalf("Error initializing test database: %v", err)
	}
	testDB = db
}

func TestInsertDicomAndGetByUUID(t *testing.T) {
	imageURL := "test_image_url_" + uuid.New().String() // Generate a random image URL
	dicomUUID := uuid.New().String()                     // Generate a random DICOM UUID

	id, err := testDB.InsertDicom(imageURL, dicomUUID)
	if err != nil {
		t.Errorf("InsertDicom failed: %v", err)
	}

	dicom, err := testDB.GetDicomByUUID(dicomUUID)
	if err != nil {
		t.Errorf("GetDicomByUUID failed: %v", err)
	}

	if dicom.ImageURL != imageURL {
		t.Errorf("Retrieved DICOM image URL doesn't match: expected %s, got %s", imageURL, dicom.ImageURL)
	}
	if dicom.ID != id {
		t.Errorf("Retrieved DICOM ID doesn't match: expected %d, got %d", id, dicom.ID)
	}
}

func TestInsertAndRetrieveTags(t *testing.T) {
	imageURL := "test2_image_url_" + uuid.New().String() 
	dicomUUID := uuid.New().String()                     

	dicomId, err := testDB.InsertDicom(imageURL, dicomUUID)
	if err != nil {
		t.Errorf("InsertDicom failed: %v", err)
	}

	tag := model.Tag{
		Tag:   "TestTag",
		VR:    "TestVR",
		Value: "TestValue",
		Name:  "TestName",
	}

	tagID, err := testDB.InsertTag(tag)
	if err != nil {
		t.Errorf("InsertTag failed: %v", err)
	}

	_, err = testDB.InsertDicomTag(dicomId, tagID)
	if err != nil {
		t.Errorf("InsertDicomTag failed: %v", err)
	}

	tags, err := testDB.GetTagsByDicomUUID(dicomUUID)
	if err != nil {
		t.Errorf("GetTagsByDicomUUID failed: %v", err)
	}

	// Validate the retrieved tags data
	if len(tags) != 1 {
		t.Errorf("Unexpected number of tags retrieved: expected 1, got %d", len(tags))
	}

	// Validate the retrieved tag details
	retrievedTag := tags[0]
	if retrievedTag.ID != tagID {
		t.Errorf("Retrieved tag ID doesn't match: expected %d, got %d", tagID, retrievedTag.ID)
	}
	if retrievedTag.Tag != tag.Tag || retrievedTag.VR != tag.VR || retrievedTag.Value != tag.Value || retrievedTag.Name != tag.Name {
		t.Errorf("Retrieved tag details don't match: expected %+v, got %+v", tag, retrievedTag)
	}
}

func TestInsertDicom_Error(t *testing.T) {
	imageURL1 := "test3_image_url_" + uuid.New().String()
	imageURL2 := "test3_image_url_" + uuid.New().String()
	dicomUUID1 := uuid.New().String() // Generate a random DICOM UUID
	dicomUUID2 := uuid.New().String() // Generate another random DICOM UUID

	_, err := testDB.InsertDicom(imageURL1, dicomUUID1)
	if err != nil {
		t.Errorf("First InsertDicom failed: %v", err)
	}

	_, err = testDB.InsertDicom(imageURL2, dicomUUID2)
	if err != nil {
		t.Errorf("Second InsertDicom failed: %v", err)
	}

	// Attempt to insert the same DICOM UUID again to trigger a unique constraint error
	_, err = testDB.InsertDicom(imageURL1, dicomUUID1)
	if err == nil {
		t.Errorf("Expected an error for duplicate DICOM UUID insertion, but got nil")
	}
}

func TestGetDicomByUUID_Error(t *testing.T) {
	// Attempt to retrieve a DICOM with a non-existing UUID
	_, err := testDB.GetDicomByUUID("non_existing_uuid")
	if err == nil {
		t.Errorf("Expected an error for non-existing DICOM UUID, but got nil")
	}
}

func TestCloseDatabase(t *testing.T) {
	// Test Close method
	err := testDB.Close()
	if err != nil {
		t.Errorf("Close method failed: %v", err)
	}
}
