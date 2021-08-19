package storage

import (
	"crypto/sha1"
	"fmt"
	"imageUploader/uploading"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	userName = os.Getenv("USERN")
	dbName   = os.Getenv("DATABASE")
	hostIP   = os.Getenv("HOST")
	port, _  = strconv.Atoi(os.Getenv("PORT"))
	psWrd    = os.Getenv("PASSWORD")
)

var connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable ", hostIP, port, userName, psWrd, dbName)

var schema = `CREATE TABLE AD_Images (
	id text,
	region text,
	description text,
	location text,
	content bytea,
	size integer,
	name text,
	added text
)`

type Storage struct {
	db *sqlx.DB
}

func NewStorage() (*Storage, error) {
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		fmt.Printf("cannot connect to db: %+v", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Printf("cannot connect to db: %+v", err)
		return nil, err
	}
	fmt.Println("Database connection success!")
	strg := new(Storage)
	strg.db = db
	return strg, nil
}

func (strg *Storage) AddImage(pic uploading.Img) error {
	id, err := getUUID(pic.Name, pic.Location)
	if err != nil {
		return err
	}
	im := Img{
		ID:        id,
		ShortDesc: pic.ShortDesc,
		Region:    pic.Region,
		Location:  pic.Location,
		Content:   pic.Content,
		Size:      pic.Size,
		Name:      pic.Name,
		CreatedAt: strconv.Itoa(int(time.Now().Unix())),
	}
	err = strg.insertImg(im)
	if err != nil {
		fmt.Printf("can not insert image: %+v", err)
		return err
	}
	return nil
}

func getUUID(str1, str2 string) (string, error) {
	data := strings.Join([]string{str1, str2}, ":")
	hash := sha1.New()
	_, err := hash.Write([]byte(data))
	if err != nil {
		return "", err
	}
	hashed := fmt.Sprintf("%x", hash.Sum(nil))
	return hashed, nil
}

func (strg *Storage) GetAllImages() ([]Img, error) {
	images := []Img{}
	rows, err := strg.db.Queryx(`SELECT * FROM AD_Images`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Img
		err = rows.StructScan(&p)
		if err != nil {
			fmt.Println(err)
		}
		images = append(images, p)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return images, nil
}

func (strg *Storage) insertImg(p Img) error {
	insertStmt := `INSERT INTO Images(id, region ,description,location ,content ,size ,name ,added)
					VALUES (:id, :region ,:description,:location ,:content ,:size ,:name ,:added)`
	_, err := strg.db.NamedExec(insertStmt, p)
	if err != nil {
		return err
	}
	return nil
}
