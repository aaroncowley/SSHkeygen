package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// Struct to hold both Keys
type PubPriv struct {
	Public_Key, Private_Key string
}

// Generates and returns
func GenKeyPair() (string, string, error) {
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

func main() {
	jsonFile, err := os.Open("names.json")

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var nameList []string
	json.Unmarshal([]byte(byteValue), &nameList)

	var jsonMap = make(map[string]PubPriv)

	for i, _ := range nameList {
		pub, priv, err := GenKeyPair()

		/*
			fmt.Println(nameList[i])
			fmt.Println(string(pub))
			fmt.Println(string(priv))
		*/

		if err != nil {
			fmt.Println(err)
		}

		keyPair := PubPriv{string(pub), string(priv)}

		jsonMap[nameList[i]] = keyPair

	}

	jsonData, _ := json.MarshalIndent(jsonMap, "", " ")
	_ = ioutil.WriteFile("output.json", jsonData, 0644)

	fmt.Println("Done")
}
