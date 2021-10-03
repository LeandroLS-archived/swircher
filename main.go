package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	baseUrl := "https://api.twitter.com/2/users/764678044405686272/followers"
	var bearer = "Bearer " + ""
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
	fmt.Println(string(body))

}
