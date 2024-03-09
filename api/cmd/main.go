package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	
	"github.com/gorilla/mux"

	"dicom/api/client"
	"dicom/api/repository/blob"
	"dicom/api/repository/sql"
	"dicom/api/service/fetcher"
	"dicom/api/service/parser"
	"dicom/api/service/processor"
)

func main() {
	logger := log.New(os.Stdout, "dicom-service: ", log.LstdFlags)
	
	blobStorage, err := blob.NewBlobStorage(logger)
	if err != nil {
		panic(err)
	}

	
	sqlRepo, err := sql.NewSqlDatabase(logger)
	if err != nil {
		panic(err)
	}

	// Instantiate fetcher service
	dicomFetcher := fetcher.NewDicomFetcher(blobStorage, sqlRepo, logger)

	// Instantiate parser service
	dicomParser := parser.NewDicomParser(sqlRepo, logger)

	// Instantiate processor service
	dicomProcessor := processor.NewDicomProcessor(sqlRepo, blobStorage, logger)

	// Set up HTTP server
	router := mux.NewRouter()
	client.Setup(router, dicomParser, dicomProcessor, dicomFetcher, logger)

	// Define server settings
	serverAddr := ":8080"
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	fmt.Println("server started")

	// Start server
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
