package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type TwitterData struct {
	Data []Followers
}

type Followers struct {
	Id       string
	Name     string
	UserName string
}

func main() {
	args := os.Args[1:]
	lenArgs := len(args)
	if lenArgs != 2 {
		log.Fatalf("Two command lines are needed %v informed", lenArgs)
	}
	baseUrl := fmt.Sprintf("https://api.twitter.com/2/users/%v/followers", args[1])
	var bearer = "Bearer " + args[0]
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
		fmt.Printf("Name: %s Username: %s\n", user.Name, user.UserName)
	}
}
