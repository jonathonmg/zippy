package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"flag"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {

	portPtr := flag.String("port", "12345", "valid tcp port")
	flag.Parse()

	log.Println("Starting server on port " + *portPtr)

	restApi := rest.NewApi()
	restApi.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/scanfile", ScanFile),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	restApi.SetApp(router)
	log.Fatal(http.ListenAndServe(":"+*portPtr, restApi.MakeHandler()))
}

type Payload struct {
	Filename   string
	Data       string
	SearchData string
}

type ReturnData struct {
	Filename string `json:"filename"`
	Result   bool   `json:"result"`
}

func ScanFile(responseWriter rest.ResponseWriter, r *rest.Request) {
	//filename := r.PathParam("filename")

	payload := Payload{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	if payload.Filename == "" {
		rest.Error(responseWriter, "filename required", 400)
		return
	}
	if payload.Data == "" {
		rest.Error(responseWriter, "data required", 400)
		return
	}
	if payload.SearchData == "" {
		rest.Error(responseWriter, "search data required", 400)
		return
	}

	fileName, result, err := ParsePayload(responseWriter, payload.Data, payload.SearchData)
	if err != nil {
		rest.Error(responseWriter, err.Error(), 400)
		return
	}
	returnData := &ReturnData{
		Filename: fileName,
		Result:   result}
	responseWriter.WriteJson(returnData)
}

func ParsePayload(responseWriter rest.ResponseWriter, fileData string, searchDataIn string) (fileNameOut string, result bool, err error) {

	decodedData, err := base64.StdEncoding.DecodeString(fileData)
	if err != nil {
		return "", false, err
	}
	decodedSearchData, err := base64.StdEncoding.DecodeString(searchDataIn)
	if err != nil {
		return "", false, err
	}
	searchData := string(decodedSearchData[:])
	readerAt := bytes.NewReader(decodedData)
	dataSize := int64(len(decodedData))
	zipReader, err := zip.NewReader(readerAt, dataSize)
	if err != nil {
		return "", false, err
	}
	for _, f := range zipReader.File {
		readCloser, err := f.Open()
		if err != nil {
			return "", false, err
		}
		fileBytes, err := ioutil.ReadAll(readCloser)
		fileContents := string(fileBytes[:])
		if strings.Contains(fileContents, searchData) {
			return f.Name, true, nil
		}
	}
	return "", false, nil
}
