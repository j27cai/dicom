package sql

import (
	"dicom/api/common"
	"dicom/api/model"

    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

// Repository defines the interface for the SQL repository
type Repository interface {
    Close() error
    InsertDicom(imageURL string, uuid string) (int64, error)
    InsertTag(tag model.Tag) (int64, error)
    InsertDicomTag(dicomID, tagID int64) (int64, error)
    GetDicomByUUID(uuid string) (*model.Dicom, error)
}

type Database struct {
    db *sql.DB
}

// Setup initializes the SQLite3 database and creates necessary tables
func Setup() (*Database, error) {
    db, err := sql.Open("sqlite3", "./dicom.db")
    if err != nil {
        return nil, err
    }

    // Create the tables if they don't exist
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS dicom (
        id INTEGER PRIMARY KEY,
        uuid string TEXT UNIQUE
        image_url TEXT UNIQUE
    )`)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS dicomTags (
        dicomId INTEGER,
        tagId INTEGER
    )`)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY,
        uuid TEXT UNIQUE
        Tag TEXT,
        VR TEXT,
        V INTEGER,
        Value TEXT
    )`)

    if err != nil {
        return nil, err
    }

    return &Database{db: db}, nil
}

func (d *Database) Close() error {
    return d.db.Close()
}

func (d *Database) InsertDicom(imageURL string, uuid string) (int64, error) {
    result, err := d.db.Exec("INSERT INTO dicom (image_url, uuid) VALUES (?, ?)", imageURL, uuid)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

func (d *Database) InsertTag(tag model.Tag) (int64, error) {
	uuid := common.GenShortUUID()

    result, err := d.db.Exec("INSERT INTO tags (uuid, Tag, VR, V, Value) VALUES (?, ?, ?, ?, ?)", uuid, tag.Tag, tag.VR, tag.Value, tag.Name)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

func (d *Database) InsertDicomTag(dicomID, tagID int64) (int64, error) {
    result, err := d.db.Exec("INSERT INTO dicomTags (dicomId, tagId) VALUES (?, ?)", dicomID, tagID)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

func (d *Database) GetDicomByUUID(uuid string) (*model.Dicom, error) {
    var dicom model.Dicom
    row := d.db.QueryRow("SELECT id, image_url FROM dicom WHERE uuid = ?", uuid)
    err := row.Scan(&dicom.ID, &dicom.ImageURL)
    if err != nil {
        return nil, err
    }
    return &dicom, nil
}
