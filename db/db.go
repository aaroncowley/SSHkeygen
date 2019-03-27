package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"

	//uncheck first 2 to enable migrate
	_ "github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	"gopkg.in/cheggaaa/pb.v1"
)

var db *sql.DB

type KeyPair struct {
	CodeName   string
	PublicKey  string
	PrivateKey string
}

func connDB() (err error) {
	connStr := "postgresql://artifacts:2019capstone!!@localhost/artifacts"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return nil
}

/*
func InitDB() (err error) {
	if err = connDB(); err != nil {
		return err
	}
	log.Println("Success in contacting DB")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"artifacts", driver)
	if err != nil {
		return err
	}

	log.Println("Running Datbase migrations")

	//run all down.sql files
	if err := m.Down(); err != nil {
		return err
	}

	// run all up.sql files
	if err := m.Force(1); err != nil {
		return err
	}

	fmt.Println("Migrations Finished")
	return nil
}
*/

func InitDB() (err error) {
	if err = connDB(); err != nil {
		return err
	}

	if _, err = db.Query("DELETE FROM profiles"); err != nil {
		return err
	}

	if _, err = db.Query("DELETE FROM keys"); err != nil {
		return err
	}

	return nil
}

func InsertKeys() (err error) {
	fmt.Println("Opening output.json")

	jsonPath, _ := filepath.Abs("jsonSeed/output.json")
	keyJson, err := os.Open(jsonPath)
	if err != nil {
		return err
	}

	byteValue, _ := ioutil.ReadAll(keyJson)

	keyList := make([]KeyPair, 0)

	err = json.Unmarshal([]byte(byteValue), &keyList)
	if err != nil {
		return err
	}

	if err = connDB(); err != nil {
		return err
	}

	/* ---- Prepared statements ---- */
	// creates new profile row, returns its id
	sqlProfileRow := `
	INSERT INTO profiles (code_name) 
	VALUES ($1)
	RETURNING id`

	// creates new key row, returns its id
	sqlKeysRow := `
	INSERT INTO keys (profile_id, public_key, private_key, creation_date)
	VALUES($1, $2, $3, DEFAULT)
	RETURNING id`

	count := len(keyList)
	bar := pb.StartNew(count)
	for _, key := range keyList {
		codeName := key.CodeName
		pubKey := key.PublicKey
		privKey := key.PrivateKey

		profID := 0
		err = db.QueryRow(sqlProfileRow, codeName).Scan(&profID)
		if err != nil {
			return err
		}

		keyID := 0
		err = db.QueryRow(sqlKeysRow, profID, pubKey, privKey).Scan(&keyID)
		if err != nil {
			return err
		}
		bar.Increment()
	}
	bar.FinishPrint("Finished injesting keys.")

	return nil
}
