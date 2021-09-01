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
	userName  = os.Getenv("USERN")
	dbName    = os.Getenv("DATABASE")
	hostIP    = os.Getenv("HOST")
	port, _   = strconv.Atoi(os.Getenv("PORT"))
	psWrd     = os.Getenv("PASSWORD")
	tableName = "AD_Images_Rel"
)

// dbname=%s
// , dbName
var connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable ", hostIP, port, userName, psWrd)

type Repo struct {
	Db *sqlx.DB
}

func NewStorage() (*Repo, error) {
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
	strg := new(Repo)
	strg.Db = db
	return strg, nil
}
func (strg *Repo) CreateDBIfNotExist() error {
	statement := `SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1);`

	row := strg.Db.QueryRowx(statement, dbName)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		fmt.Printf("Data base already exists")
		return nil
	}

	if !exists {
		statement = `CREATE DATABASE $1;`
		strg.Db.MustExec(statement, dbName)
		fmt.Println("create new database success!")
	}
	db1, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=disable ", hostIP, port, userName, psWrd, dbName))
	if err != nil {
		return err
	}
	strg.Db = db1
	return nil
}

func (strg *Repo) AddImage(pic uploading.Img) error {
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

func (strg *Repo) GetAllImages() ([]Img, error) {
	images := []Img{}
	rows, err := strg.Db.Queryx(`SELECT * FROM $1`, dbName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if images, err = iterateResult(rows); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return images, nil
}

func (strg *Repo) insertImg(p Img) error {
	insertStmt := `INSERT INTO AD_Images_Rel(id, region ,description,location ,content ,size ,name ,added)
					VALUES (:id, :region ,:description,:location ,:content ,:size ,:name ,:added)`
	_, err := strg.Db.NamedExec(insertStmt, p)
	if err != nil {
		return err
	}
	return nil
}

func (strg *Repo) GetImagesByCriteria(region, location string) ([]Img, error) {
	var (
		regArg, locArg string
		rows           *sqlx.Rows
		images         = make([]Img, 0)
		err            error
	)
	regArg, ok1 := validateCriteria(region)
	locArg, ok2 := validateCriteria(location)
	switch {
	case ok1 && ok2:
		if rows, err = strg.Db.Queryx(`SELECT * FROM AD_Images_Rel
								WHERE region = $1
								AND location = $2 `, regArg, locArg); err != nil {
			fmt.Println(err)
			return nil, err
		}
		if images, err = iterateResult(rows); err != nil {
			fmt.Println(err)
			return nil, err
		}

	case ok1 && !ok2:
		if rows, err = strg.Db.Queryx(`SELECT * FROM AD_Images_Rel
								WHERE region = $1 and `, regArg); err != nil {
			fmt.Println(err)
			return nil, err
		}
		if images, err = iterateResult(rows); err != nil {
			fmt.Println(err)
			return nil, err
		}

	case !ok1 && ok2:
		if rows, err = strg.Db.Queryx(`SELECT * FROM AD_Images_Rel
									WHERE location = $1 and `, locArg); err != nil {
			fmt.Println(err)
			return nil, err
		}
		if images, err = iterateResult(rows); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return images, nil
}
func iterateResult(rows *sqlx.Rows) ([]Img, error) {
	imgz := []Img{}
	defer rows.Close()
	for rows.Next() {
		var p Img
		err := rows.StructScan(&p)
		if err != nil {
			fmt.Println(err)
		}
		imgz = append(imgz, p)
	}
	err := rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	return imgz, nil
}
func validateCriteria(str string) (string, bool) {
	if str != "" {
		return str, true
	}
	return str, false
}

//
//
// 	if ; ok{
// 		 =location
// 	}

// SELECT EXISTS (
// 	SELECT FROM information_schema.tables
// 	WHERE  table_schema = 'schema_name'
// 	AND    table_name   = 'table_name'
// 	);

// SELECT 1 FROM information_schema.tables
// WHERE table_schema = 'schema_name'	AND table_name = 'table_name';

// SELECT 'CREATE DATABASE mydb'
// WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'AD_Images')
