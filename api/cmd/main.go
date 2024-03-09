package main

import (
	"net/http"
	"github.com/gorilla/mux"

	"dicom/api/client"
	"dicom/api/repository/block"
	"dicom/api/repository/sql"
	"dicom/api/service/fetcher"
	"dicom/api/service/parser"
	"dicom/api/service/processor"
)

func main() {
	// Instantiate block storage repository
	blockStorage, err := block.Setup()
	if err != nil {
		panic(err)
	}

	// Instantiate SQL repository
	sqlRepo, err := repository.Setup()
	if err != nil {
		panic(err)
	}

	// Instantiate fetcher service
	dicomFetcher := fetcher.NewDicomProvider(blockStorage, sqlRepo)

	// Instantiate parser service
	dicomParser := parser.NewDicomConverter(sqlRepo)

	// Instantiate processor service
	dicomProcessor := processor.NewDicomExtractor(sqlRepo, blockStorage)

	// Set up HTTP server
	router := mux.NewRouter()
	client.Setup(router, dicomParser, dicomProcessor)

	// Define server settings
	serverAddr := ":8080"
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
