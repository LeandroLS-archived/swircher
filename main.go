package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	fmt.Printf("Check if %v has some follower with this User Name: ", userToCheck)
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

func makeRequest(method string, url string, authorization string) []byte {
	req, err := http.NewRequest(method, url, nil)
	handleErr(err)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	fmt.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	handleErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body), err)
	handleErr(err)
	return body
}

func makeTweet(message string) []byte {
	baseUrl := fmt.Sprintf("https://api.twitter.com/1.1/statuses/update.json?status=%s", url.QueryEscape(message))
	fmt.Println(baseUrl)
	auth0 := "super secret token"
	body := makeRequest("POST", baseUrl, auth0)
	return body
}

func getUserByUserName(userName string) TwitterAPIResponseUserInfo {
	baseUrl := fmt.Sprint("https://api.twitter.com/2/users/by/username/", userName)
	bearerToken := "Bearer " + os.Getenv("TWITTER_BEARER_TOKEN")
	body := makeRequest("GET", baseUrl, bearerToken)
	var twitterData TwitterAPIResponseUserInfo
	err := json.Unmarshal([]byte(body), &twitterData)
	handleErr(err)
	return twitterData
}

func getUserChoise() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Options:")
	fmt.Println("1 => Make a Tweet")
	fmt.Println("2 => See who follows a specific user")
	fmt.Println("Enter the numbers 1 or 2 to chose what you wanna do:")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	number, err := strconv.Atoi(text)
	handleErr(err)
	if number != 1 && number != 2 {
		fmt.Println("Please inform a number between 1 or 2")
		getUserChoise()
	}
	return number
}

func main() {
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")
	if bearerToken == "" {
		log.Fatalln("Environment Variable TWITTER_BEARER_TOKEN not found")
	}
	numberChoise := getUserChoise()
	if numberChoise == 1 {
		fmt.Println("Enter your Tweet message:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSuffix(text, "\r")
		body := makeTweet(text)
		fmt.Println(string(body))
	} else {
		userToCheck, followedBy := askUser()
		user := getUserByUserName(userToCheck)
		baseUrl := fmt.Sprintf("https://api.twitter.com/2/users/%v/followers", user.Data.Id)
		body := makeRequest("GET", baseUrl, "Bearer "+bearerToken)
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

}
