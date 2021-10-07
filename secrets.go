package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type SecretAPIKey struct {
	BearerToken string
}

func writeSecretsFile(s string) {
	secret := SecretAPIKey{s}
	secretJson, err := json.Marshal(secret)
	check(err)
	err = ioutil.WriteFile("secret.json", []byte(secretJson), 0755)
	check(err)
}

func secretsFileExists() bool {
	_, err := ioutil.ReadFile("secret.json")
	if err != nil {
		return false
	}
	return true
}
