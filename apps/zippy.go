package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
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
	}
	restApi.SetApp(router)
	log.Fatal(http.ListenAndServe(":12345", restApi.MakeHandler()))
}

type ZipFile struct {
	Filename string
	Data     string
}

func ScanFile(w rest.ResponseWriter, r *rest.Request) {
	//filename := r.PathParam("filename")

	zipFile := ZipFile{}
	err := r.DecodeJsonPayload(&zipFile)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if zipFile.Filename == "" {
		rest.Error(w, "filename required", 400)
		return
	}
	if zipFile.Data == "" {
		rest.Error(w, "data required", 400)
		return
	}
	w.WriteJson(&zipFile)
}
