package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"imageUploader/uploading"

	"github.com/gorilla/mux"
)

type ImageManager struct {
	l     log.Logger
	adder uploading.Service
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
	image := uploading.Img{}
	image.Name = fheader.Filename
	image.Size = int(fheader.Size)
	image.ShortDesc = r.Form["description"][0]
	image.Location = r.Form["location"][0]
	image.Region = r.Form["region"][0]

	// read all of the contents of our uploaded file into a
	// byte array
	imgBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	image.Content = imgBytes
	err = im.adder.AddImages(image)
	if err != nil {
		fmt.Printf("failed to add image : %+v", err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintf(rw, "Image upload successful %s\n", image.Name)
	return
}

func SetUpRoutes(man ImageManager) {
	sm := mux.NewRouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/upload", man.UploadFile)
}
