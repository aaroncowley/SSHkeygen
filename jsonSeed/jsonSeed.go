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
	mrand "math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
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

		mux.Lock() // lock resources
		*list = append(*list, keyPair)
		bar.Increment()
		mux.Unlock() // unlock resources

		wg.Done()
	}
}

func CreateJsonSeed(keyNum int) {
	jsonFile, err := os.Open("jsonSeed/names.json")
	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var nameList []string

	json.Unmarshal([]byte(byteValue), &nameList)

	jsonList := make([]PubPriv, 0)

	// Channels for workers
	jobs := make(chan string, keyNum)

	var wg sync.WaitGroup
	var mux sync.Mutex

	// init worker
	workerNum := runtime.NumCPU()
	nameNum := len(nameList)
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("Starting key generation with [%s] workers. This is going to take awhile ...\n",
		red(workerNum))

	color.Set(color.FgMagenta, color.Bold)
	bar := pb.StartNew(keyNum)
	for i := 0; i < workerNum; i++ {
		go genWorker(1, &wg, &mux, jobs, &jsonList, bar)
	}

	//fill job queue
	for i := 0; i < keyNum; i++ {
		now := strconv.FormatInt(time.Now().UnixNano(), 10)
		codeName := strings.Join([]string{nameList[mrand.Intn(nameNum)], now}, "")
		jobs <- codeName
		wg.Add(1)
	}

	wg.Wait()
	close(jobs)
	color.Unset()

	color.Set(color.FgCyan)
	bar.FinishPrint("Key Generation Completed")
	fmt.Printf("%s Keys Generated\n", red(keyNum))
	color.Unset()

	_ = os.Remove("jsonSeed/koutput.json")
	f, err := os.OpenFile("jsonSeed/koutput.json",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for _, jStr := range jsonList {
		jsonData, err := json.Marshal(jStr)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write(jsonData); err != nil {
			log.Println(err)
		}
		_, _ = f.WriteString("\n")
	}

	jsonData, err := json.MarshalIndent(jsonList, "", " ")
	err = ioutil.WriteFile("jsonSeed/output.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done")
	fmt.Println()
}
