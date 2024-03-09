package sql

import (
    "dicom/api/common"
    "dicom/api/model"

    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

// Repository defines the interface for the SQL repository
type Repository interface {
    Close() error
    InsertDicom(imageURL string, uuid string) (int64, error)
    InsertTag(tag model.Tag) (int64, error)
    InsertDicomTag(dicomID, tagID int64) (int64, error)
    GetDicomByUUID(uuid string) (*model.Dicom, error)
    GetTagsByDicomUUID(uuid string) ([]model.Tag, error)
}

type Database struct {
    db     *sql.DB
    logger *log.Logger
}

func NewSqlDatabase(logger *log.Logger) (*Database, error) {
    db, err := sql.Open("sqlite3", "./dicom.db")
    if err != nil {
        logger.Printf("Error opening database: %v", err)
        return nil, err
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS dicom (
        id INTEGER PRIMARY KEY,
        uuid string TEXT UNIQUE,
        image_url TEXT UNIQUE
    )`)
    if err != nil {
        logger.Printf("Error creating dicom table: %v", err)
        return nil, err
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS dicomTags (
        dicomId INTEGER,
        tagId INTEGER
    )`)
    if err != nil {
        logger.Printf("Error creating dicomTags table: %v", err)
        return nil, err
    }

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY,
        uuid TEXT UNIQUE,
        Tag TEXT,
        VR TEXT,
        Value TEXT,
        Name TEXT
    )`)

    if err != nil {
        logger.Printf("Error creating tags table: %v", err)
        return nil, err
    }

    return &Database{db: db, logger: logger}, nil
}


func (d *Database) Close() error {
    return d.db.Close()
}

func (d *Database) InsertDicom(imageURL string, uuid string) (int64, error) {
    result, err := d.db.Exec("INSERT INTO dicom (image_url, uuid) VALUES (?, ?)", imageURL, uuid)
    if err != nil {
        d.logger.Printf("Error inserting DICOM: %v", err)
        return 0, err
    }

    return result.LastInsertId()
}

func (d *Database) InsertTag(tag model.Tag) (int64, error) {
    uuid := common.GenShortUUID()

    result, err := d.db.Exec("INSERT INTO tags (uuid, Tag, VR, Value, Name) VALUES (?, ?, ?, ?, ?)", uuid, tag.Tag, tag.VR, tag.Value, tag.Name)
    if err != nil {
        d.logger.Printf("Error inserting tag: %v", err)
        return 0, err
    }

    return result.LastInsertId()
}

func (d *Database) InsertDicomTag(dicomID, tagID int64) (int64, error) {
    result, err := d.db.Exec("INSERT INTO dicomTags (dicomId, tagId) VALUES (?, ?)", dicomID, tagID)
    if err != nil {
        d.logger.Printf("Error inserting DICOM tag: %v", err)
        return 0, err
    }

    return result.LastInsertId()
}

func (d *Database) GetDicomByUUID(uuid string) (*model.Dicom, error) {
    var dicom model.Dicom
    row := d.db.QueryRow("SELECT id, image_url FROM dicom WHERE uuid = ?", uuid)
    err := row.Scan(&dicom.ID, &dicom.ImageURL)
    if err != nil {
        d.logger.Printf("Error getting DICOM by UUID: %v", err)
        return nil, err
    }
    return &dicom, nil
}

func (d *Database) GetTagsByDicomUUID(uuid string) ([]model.Tag, error) {
    var tags []model.Tag

    // Select all tags associated with the DICOM UUID
    query := `
        SELECT tags.id, tags.Tag, tags.VR, tags.Value, tags.Name
        FROM dicomTags
        JOIN tags ON dicomTags.tagId = tags.id
        JOIN dicom ON dicomTags.dicomId = dicom.id
        WHERE dicom.uuid = ?
    `

    rows, err := d.db.Query(query, uuid)
    if err != nil {
        d.logger.Printf("Error getting tags by DICOM UUID: %v", err)
        return nil, err
    }
    defer rows.Close()

    // Iterate through the rows and populate the tags slice
    for rows.Next() {
        var tag model.Tag
        if err := rows.Scan(&tag.ID, &tag.Tag, &tag.VR, &tag.Value, &tag.Name); err != nil {
            d.logger.Printf("Error scanning tag row: %v", err)
            return nil, err
        }
        tags = append(tags, tag)
    }
    if err := rows.Err(); err != nil {
        d.logger.Printf("Error iterating over tag rows: %v", err)
        return nil, err
    }

    return tags, nil
}

