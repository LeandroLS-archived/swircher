package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type TwitterData struct {
	Data   []Followers
	Status int
}

type Followers struct {
	Id       string
	Name     string
	UserName string
}

func isNameSimilar(userName string, userToBeFound string) bool {
	return strings.Contains(userName, userToBeFound)
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	args := os.Args[1:]
	lenArgs := len(args)
	if lenArgs != 3 {
		log.Fatalf("Two command lines are needed %v informed", lenArgs)
	}
	bearerToken := args[0]
	userId := args[1]
	userToBeFound := args[2]
	if !secretsFileExists() {
		writeSecretsFile(bearerToken)
	}
	baseUrl := fmt.Sprintf("https://api.twitter.com/2/users/%v/followers", userId)
	var bearer = "Bearer " + bearerToken
	req, err := http.NewRequest("GET", baseUrl, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	handleErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	var twitterData TwitterData
	err = json.Unmarshal([]byte(body), &twitterData)
	handleErr(err)
	//When twitter API returns some status, is because the API fails. Weird :(
	if twitterData.Status == 0 {
		log.Fatalln("Cant make request", err)
	}
	for _, user := range twitterData.Data {
		fmt.Println(user)
		if isNameSimilar(user.Name, userToBeFound) {
			fmt.Printf("Name: %s Username: %s\n", user.Name, user.UserName)
		}
	}

}
