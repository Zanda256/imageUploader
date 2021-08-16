package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ImageManager struct {
	l log.Logger
}

func (im *ImageManager) UploadFile(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	file, fheader, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	name := fheader.Filename
	size := fheader.Size
	desc := r.Form["description"]
	locale := r.Form["location"]
	region := r.Form["region"]

	// read all of the contents of our uploaded file into a
	// byte array
	imgBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

}
