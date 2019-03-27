package jsonSeed

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	_ "time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/cheggaaa/pb.v1"
)

// Struct to hold both Keys
type PubPriv struct {
	CodeName   string
	PublicKey  string
	PrivateKey string
}

// Generates and returns
func genKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	var private bytes.Buffer
	if err := pem.Encode(&private, privateKeyPEM); err != nil {
		return "", "", err
	}

	// generate public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	public := ssh.MarshalAuthorizedKey(pub)
	return string(public), private.String(), nil
}

func genWorker(id int, wg *sync.WaitGroup, mux *sync.Mutex, jobs <-chan string, list *[]PubPriv, bar *pb.ProgressBar) {
	for n := range jobs {
		codeName := n
		pub, priv, err := genKeyPair()
		if err != nil {
			log.Fatal(err)
		}

		keyPair := PubPriv{codeName, string(pub), string(priv)}
		//fmt.Printf("%s new struct\n", codeName)

		mux.Lock()
		*list = append(*list, keyPair)
		bar.Increment()
		mux.Unlock()

		wg.Done()
	}
}

func CreateJsonSeed() {
	jsonFile, err := os.Open("jsonSeed/names.json")
	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var nameList []string

	json.Unmarshal([]byte(byteValue), &nameList)

	fmt.Println("Starting key generation. This is going to take awhile...")

	jsonList := make([]PubPriv, 0)

	// Max is len(nameList)
	//nameNum := len(nameList)
	nameNum := 20
	bar := pb.StartNew(nameNum)

	// Channels for workers
	jobs := make(chan string, nameNum)

	var wg sync.WaitGroup
	var mux sync.Mutex

	// init worker
	workerNum := 2
	for i := 0; i < workerNum; i++ {
		go genWorker(1, &wg, &mux, jobs, &jsonList, bar)
	}

	//fill job queue
	for _, name := range nameList[:nameNum] {
		jobs <- name
		wg.Add(1)
	}

	wg.Wait()
	close(jobs)

	bar.FinishPrint("Key Generation Completed")
	fmt.Printf("%d Keys Generated\n", nameNum)

	jsonData, err := json.MarshalIndent(jsonList, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("jsonSeed/output.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done")
}
