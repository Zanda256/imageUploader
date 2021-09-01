package main

import (
	"fmt"
	storage "imageUploader/storage/RD"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	s, err := storage.NewStorage()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%T", s)
	err = InitDb(s)
	if err != nil {
		fmt.Printf("Major error:%s", err)
	}
}

func InitDb(s *storage.Repo) error {
	err := s.CreateDBIfNotExist()
	if err != nil {
		return err
	}
	createTableCmd := loadDBScript()
	s.Db.MustExec(createTableCmd)
	fmt.Println("create new table success!")
	return nil
}

func loadDBScript() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get working directory %+v", err)
	}
	fileDir := fmt.Sprintf("%s/%s", currentDir, "storage/RD/init_db.sql")
	initBytes, err := ioutil.ReadFile(fileDir)
	if err != nil {
		log.Fatalf("cannot read init_db file %+v", err)
	}
	return string(initBytes)
}
