package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"github.com/kkeuning/gobservatory/gobservatory-cms/content"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type StarCollection struct {
	Stars []content.Star `json:"data"`
}

func (sc *StarCollection) Contains(s content.Star) bool {
	for _, star := range sc.Stars {
		if star.GithubId == s.GithubId {
			return true
		}
	}
	return false
}

func PostToPonzu(s content.Star, ponzuURL string, ponzuKey string) error {
	ponzuClient := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	starStruct := structs.New(s)
	for _, f := range starStruct.Fields() {
		fmt.Printf("field name: %+v\n", f.Name())
		fmt.Printf("json field name: %+v\n", f.Tag("json"))
		fmt.Printf("is zero : %+v\n", f.IsZero())
		if f.IsZero() == false {
			writer.WriteField(f.Tag("json"), fmt.Sprint(f.Value()))
			fmt.Printf("value   : %+v\n", f.Value())
		}
	}
	boundary := writer.Boundary()
	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", ponzuURL, body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Headers
	req.Header.Add("Content-Type", "multipart/form-data; charset=utf-8; boundary="+boundary)

	req.Header.Add("Content-Type", writer.FormDataContentType())
	//TODO: Add header for client secret

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
		return err
	}

	// Fetch Request
	resp, err := ponzuClient.Do(req)
	if err != nil {
		fmt.Println("Failure : ", err)
		return err
	}

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failure : ", err)
		return err
	}
	fmt.Println(string(respBody))
	fmt.Println(resp.Status)

	return nil
}

func GetFromPonzu(ponzuURL string, ponzuKey string) (*StarCollection, error) {
	var stars StarCollection
	ponzuClient := &http.Client{}
	ponzuReq, err := http.NewRequest("GET", ponzuURL, nil)
	if err != nil {
		fmt.Println("error:", err)
	}
	//TODO: Add header for client secret
	resp, err := ponzuClient.Do(ponzuReq)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(resp.Status))

	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &stars)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, s := range stars.Stars {
		fmt.Println(s.FullName)
	}
	return &stars, nil
}
