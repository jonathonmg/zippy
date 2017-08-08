package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"github.com/ant0ine/go-json-rest/rest"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func scanfile(res http.ResponseWriter, req *http.Request) {

}

func main() {

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
	log.Fatal(http.ListenAndServe(":12345", restApi.MakeHandler()))
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

func ScanFile(w rest.ResponseWriter, r *rest.Request) {
	//filename := r.PathParam("filename")

	payload := Payload{}
	err := r.DecodeJsonPayload(&payload)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if payload.Filename == "" {
		rest.Error(w, "filename required", 400)
		return
	}
	if payload.Data == "" {
		rest.Error(w, "data required", 400)
		return
	}
	if payload.SearchData == "" {
		rest.Error(w, "search data required", 400)
		return
	}

	//WriteFile(payload.Filename, payload.Data)
	fileName, result := ParsePayload(payload.Data, payload.SearchData)
	returnData := &ReturnData{
		Filename: fileName,
		Result:   result}
	w.WriteJson(returnData)
}

func ParsePayload(fileData string, searchDataIn string) (fileNameOut string, result bool) {

	decodedData, err := base64.StdEncoding.DecodeString(fileData)
	CheckError(err)
	decodedSearchData, err := base64.StdEncoding.DecodeString(searchDataIn)
	CheckError(err)
	searchData := string(decodedSearchData[:])
	readerAt := bytes.NewReader(decodedData)
	dataSize := int64(len(decodedData))
	zipReader, err := zip.NewReader(readerAt, dataSize)
	CheckError(err)
	for _, f := range zipReader.File {
		readCloser, err := f.Open()
		CheckError(err)
		fileBytes, err := ioutil.ReadAll(readCloser)
		fileContents := string(fileBytes[:])
		if strings.Contains(fileContents, searchData) {
			return f.Name, true
		}
	}
	return "", false
}

func WriteFile(fileName string, fileData string) {
	decodedData, err := base64.StdEncoding.DecodeString(fileData)
	CheckError(err)
	fileHandle, err := os.Create(fileName)
	CheckError(err)
	defer fileHandle.Close()
	if _, err := fileHandle.Write(decodedData); err != nil {
		log.Fatal(err)
	}
	if err := fileHandle.Sync(); err != nil {
		log.Fatal(err)
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
