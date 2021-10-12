package main

import (
	"bufio"
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

type UserInfo struct {
	Id       string
	Name     string
	UserName string
}

type TwitterAPIResponseFollowers struct {
	Data   []UserInfo
	Status int
}

func askUser() (userToCheck string, isFollowedBy string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter User: ")
	scanner.Scan()
	userToCheck = scanner.Text()
	fmt.Printf("Check if %v has some follower with this name: ", userToCheck)
	scanner.Scan()
	isFollowedBy = scanner.Text()
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
	return userToCheck, isFollowedBy
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
	userToCheck, followedBy := askUser()
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		log.Fatalln("Environment Variable TWITTER_BEARER_TOKEN not found")
	}
	user := getUserByUserName(userToCheck)
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
		if isNameSimilar(user.Name, followedBy) {
			fmt.Printf("Name: %s Username: %s\n", user.Name, user.UserName)
		}
	}

}
