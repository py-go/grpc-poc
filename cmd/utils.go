package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	BaseAPI = "https://dog.ceo/api"
)

func getData(url string) (body []byte, status int, err error) {
	// We get back a *Response, and an error
	response, err := http.Get(url)

	if err != nil {
		log.Printf("http.Get -> %v\n", err)
		return
	}
	defer response.Body.Close()
	status = response.StatusCode
	// Transform our response to a []byte
	body, err = ioutil.ReadAll(response.Body)

	if err != nil {
		log.Printf("ioutil.ReadAll -> %v\n", err)
		return
	}
	return
}
func getRandomImageURL(name string) (url string, err error) {
	//Call Image API in order to get Dog Images's URL
	body, status, err := getData(BaseAPI + fmt.Sprintf("/breed/%s/images/random", name))
	if err != nil {
		return
	}
	if status != 200 {
		err = fmt.Errorf("can not get the images :-(")
		return

	}
	// Put only needed informations of the JSON document in our array of Dog Images
	var data APIResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("failed to unmarshal JSON: %v\n", err)
		return
	}
	url = data.Message
	return
}
