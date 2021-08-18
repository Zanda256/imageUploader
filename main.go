package main

import (
	"fmt"
	storage "imageUploader/storage/RD"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T", s)
}
