# dicom

database for now:

sqlite (dicom - id, image_url
	    dicom-tags - dicom id, tag id
	    tag - id, Tag, VR, Name)

backup/block storage for now:

local s3

backend:

golang



Later:

File encryption, security, monitoring/logging/observability (ELK/APM), error rate, accuracy of image generation, Upload in chunks (30mb) and if it breaks, duplicate case, damaged data, 
