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

type TwitterAPIResponseUserInfo struct {
	Data UserInfo
}

type TwitterAPIResponseFollowers struct {
	Data   []Followers
	Status int
}

type UserInfo struct {
	Id       string
	Name     string
	UserName string
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

func makeRequest(method string, url string, bearerToken string) []byte {
	req, err := http.NewRequest(method, url, nil)
	handleErr(err)
	bearer := bearerToken
	req.Header.Add("Authorization", "Bearer "+bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	handleErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	return body
}

func getUserByUserName(userName string) TwitterAPIResponseUserInfo {
	baseUrl := fmt.Sprint("https://api.twitter.com/2/users/by/username/", userName)
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	body := makeRequest("GET", baseUrl, bearerToken)
	var twitterData TwitterAPIResponseUserInfo
	err := json.Unmarshal([]byte(body), &twitterData)
	handleErr(err)
	return twitterData
}

func main() {
	fmt.Println()
	args := os.Args[1:]
	lenArgs := len(args)
	if lenArgs != 2 {
		log.Fatalf("Two command lines are needed %v informed", lenArgs)
	}
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		log.Fatalln("Environment Variable TWITTER_BEARER_TOKEN not found")
	}
	user := getUserByUserName("le_limasilva")
	userToBeFound := args[1]
	baseUrl := fmt.Sprintf("https://api.twitter.com/2/users/%v/followers", user.Data.Id)
	body := makeRequest("GET", baseUrl, bearerToken)
	var twitterData TwitterAPIResponseFollowers
	err := json.Unmarshal([]byte(body), &twitterData)
	handleErr(err)
	// When twitter API returns some status, is because the API fails. Weird :(
	if twitterData.Status != 0 {
		log.Fatalln("Cant make request", err)
	}
	for _, user := range twitterData.Data {
		if isNameSimilar(user.Name, userToBeFound) {
			fmt.Printf("Name: %s Username: %s\n", user.Name, user.UserName)
		}
	}

}
