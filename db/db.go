package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
	"time"

	//uncheck first 2 to enable migrate
	_ "github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	"github.com/fatih/color"
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

	red := color.New(color.FgRed).SprintFunc()
	color.Set(color.FgCyan)
	fmt.Printf("Resetting Database in ")
	for i := 5; i > 0; i-- {
		fmt.Printf("%s...\n", red(i))
		time.Sleep(time.Second)
	}
	fmt.Println(red("0"))

	if _, err = db.Query("DELETE FROM keys"); err != nil {
		return err
	}
	if _, err = db.Query("DELETE FROM profiles"); err != nil {
		return err
	}

	fmt.Println("Database", red("Reset"))
	fmt.Println()
	return nil
}

func InsertKeys() (err error) {
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

	color.Set(color.FgMagenta)
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
	color.Unset()

	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s Rows Inserted\n", red(count))

	fmt.Println()
	return nil
}

func PullNextKey() (keys KeyPair, err error) {
	err = connDB()
	if err != nil {
		return keys, err
	}

	pullKeyRow := `
	SELECT profiles.code_name, keys.public_key, keys.private_key, profiles.id
	FROM keys, profiles
	WHERE profiles.id = keys.profile_id AND profiles.used != TRUE
	ORDER BY keys.creation_date DESC
	LIMIT 1`

	var profID int
	var retKeys KeyPair

	err = db.QueryRow(pullKeyRow).Scan(
		&retKeys.CodeName,
		&retKeys.PublicKey,
		&retKeys.PrivateKey,
		&profID)
	if err != nil {
		return keys, err
	}

	//fmt.Printf("%+v\n", retKeys)

	updateUsedRow := `
	UPDATE profiles SET used = TRUE WHERE id = $1
	RETURNING code_name`

	var usedCodeName string
	err = db.QueryRow(updateUsedRow, profID).Scan(&usedCodeName)
	if err != nil {
		return keys, err
	}

	//fmt.Printf("Updated %s as used\n", usedCodeName)
	keys = retKeys
	return keys, nil
}
