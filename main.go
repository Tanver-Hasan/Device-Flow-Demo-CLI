package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mdp/qrterminal"
)

var loginCmd bool
var m map[string]interface{}

func main() {
	fmt.Println("Test Device Credentials")
	flag.BoolVar(&loginCmd, "login", true, "Authenticate device")
	flag.Parse()

	if loginCmd == true {
		login()

	}
}

func login() {
	url := "https://login.tanverhasan.com/oauth/device/code"

	reqBody, err := json.Marshal(map[string]string{
		"client_id": "zEHWunlajvzQKAYQ54o0D5sZ3iWz9BvE",
		"scope":     "openid",
		"audience":  "https://api.timesheet.com/",
	})
	if err != nil {
		log.Fatalln(err)
	}
	res, _ := http.Post(url, "application/json", bytes.NewBuffer(reqBody))

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	var f interface{}
	json.Unmarshal(body, &f)

	m := f.(map[string]interface{})
	deviceCode := fmt.Sprintf("%v", m["device_code"])
	fmt.Println("Prinint Device Code " + deviceCode)
	str := fmt.Sprintf("%v", m["verification_uri_complete"])

	generateQrCode(str)

	for {
		time.Sleep(20 * time.Second)
		fmt.Println("Pooling loop")
		resBody, err := poolingToken(deviceCode)

		if err != nil {
			fmt.Print(err)
		}

		fmt.Println(string(resBody))

		if len(resBody) > 0 {
			fmt.Println(string(resBody))
			var ff interface{}
			json.Unmarshal(resBody, &ff)
			result := ff.(map[string]interface{})
			accessToken := fmt.Sprintf("%v", result["access_token"])
			fmt.Println(accessToken)
			if len(accessToken) > 5 {
				break
			}
		}

	}
}

func generateQrCode(url string) {

	qrterminal.Generate(url, qrterminal.L, os.Stdout)
}

func poolingToken(deviceCode string) ([]byte, error) {
	// fmt.Println("Calling pooloing token function")
	url := "https://login.tanverhasan.com/oauth/token"

	fmt.Println("Priting device code in poling function :" + deviceCode)

	reqBody, err := json.Marshal(map[string]string{
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		"device_code": deviceCode,
		"client_id":   "zEHWunlajvzQKAYQ54o0D5sZ3iWz9BvE",
	})

	if err != nil {
		fmt.Println(err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return body, err

}
