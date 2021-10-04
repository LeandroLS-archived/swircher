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
	Data []Followers
}

type Followers struct {
	Id       string
	Name     string
	UserName string
}

func isNameSimilar(userName string, userToBeFound string) bool {
	return strings.Contains(userName, userToBeFound)
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
	baseUrl := fmt.Sprintf("https://api.twitter.com/2/users/%v/followers", userId)
	var bearer = "Bearer " + bearerToken
	req, err := http.NewRequest("GET", baseUrl, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Cant make request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Cant make request")
	}
	var twitterData TwitterData

	err = json.Unmarshal(body, &twitterData)

	if err != nil {
		fmt.Println(err)
	}

	for _, user := range twitterData.Data {
		if isNameSimilar(user.Name, userToBeFound) {
			fmt.Printf("Name: %s Username: %s\n", user.Name, user.UserName)
		}
	}
}
