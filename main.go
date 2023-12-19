package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type bordapiResponse struct {
	Activity      string
	Key           string
	Accessibility string
}

type picResponse struct {
	Hits []picReponseHit
}

type picReponseHit struct {
	WebformatURL string
}

func main() {

	response, err := http.Get("https://www.boredapi.com/api/activity")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseStruct bordapiResponse
	json.Unmarshal([]byte(responseData), &responseStruct)

	fmt.Println(string(responseStruct.Activity))
	fmt.Printf("boredAPI activity: %+v\n", responseStruct.Activity)

	// responseStruct.Activity contains the string to query picture api

	//responsePic, err := http.Get(fmt.Sprintf("https://pixabay.com/api?key=31354828-b89960e4ae5d1dc984f8559d7&q=%v", responseStruct.Activity))
	responsePic, err := http.Get("https://pixabay.com/api?key=31354828-b89960e4ae5d1dc984f8559d7&q=sunflowers")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseDataPic, err := ioutil.ReadAll(responsePic.Body)
	if err != nil {
		log.Fatal(err)
	}

	var picStruct picResponse
	json.Unmarshal([]byte(responseDataPic), &picStruct)

	fmt.Printf("pic link: %+v\n", picStruct.Hits[0].WebformatURL)

	urlToDisplay := picStruct.Hits[0].WebformatURL

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.HTML(http.StatusOK, "index.tmpl", []string{"a", "b", "c"})

		fmt.Fprintf(w, "Hello, you've requested: %s\n", urlToDisplay)
	})

	http.ListenAndServe(":8080", nil)
}
