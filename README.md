# dicom

This service reads from a local .dcm, process the dicom file, stores the various images/tags to persistence memory and then presents them

![alt text](https://github.com/j27cai/dicom/blob/main/architecture.png?raw=true)

## running the service

`docker-compose up`

on root level

runs on port 8001 by default per definition in docker-compose.yml

# REST API

The REST API to the dicom parser is described below.

## Process a new dicom file

Retrieves a file locally (for now) and then process the dicom file, storing the tags in a sql database and the image in "blob" storage (local) 

### Request

	`POST /dicom/`

    curl --location 'localhost:8000/dicom' \ --header 'Content-Type: application/json' \ --data '{"path": "test-xray/DICOM/PA000001/ST000001/SE000001/IM000002"}'

### Response

    {"id": "iEfcZk3Vn6H8iyqc3seHrm"}

## Get an image for a processed dicom file

Gets a image through a query parameter for a uniquely indentifiable dicom file provided as a response to the /dicom endpoint

### Request

	`GET /image`

    curl --location 'localhost:8001/image?id=iEfcZk3Vn6H8iyqc3seHrm'

### Response

    HTTP/1.1 200 OK
    Date: Sat, 09 Mar 2024 20:38:07 GMT
    Status: 200 OK
    Content-Type: image/png
    Transfer-Encoding: chunked

## Get all tags for a processed dicom file

Gets a image through a query parameter for a uniquely indentifiable dicom file provided as a response to the /dicom endpoint

### Request

	`GET /tags`

    curl --location 'localhost:8001/tags?id=iEfcZk3Vn6H8iyqc3seHrm'

### Response

    HTTP/1.1 200 OK
    Date: Sat, 09 Mar 2024 20:38:07 GMT
    Status: 200 OK
    Content-Type: application/json
    Transfer-Encoding: chunked

    {
	    "uuid": "iEfcZk3Vn6H8iyqc3seHrm",
	    "tags": [
	        {
	            "ID": 125,
	            "Tag": "(0002,0000)",
	            "VR": "VRUInt32List",
	            "Value": "[198]",
	            "Name": "FileMetaInformationGroupLength"
	        },
	        {
	            "ID": 126,
	            "Tag": "(0002,0001)",
	            "VR": "VRBytes",
	            "Value": "[0 1]",
	            "Name": "FileMetaInformationVersion"
	        }
	    ]
	}




Scratch notes

database for now:

sqlite (dicom - id, image_url
	    dicom-tags - dicom id, tag id
	    tag - id, Tag, VR, Name)

backup/blob storage for now:

local disk

backend:

golang

Later:

File encryption, security, monitoring/logging/observability (ELK/APM), error rate, accuracy of image generation, Upload in chunks (30mb) and if it breaks, duplicate case, damaged data, 
