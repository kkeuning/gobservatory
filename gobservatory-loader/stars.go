package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type Star struct {
	Name            string   `json:"name"`
	FullName        string   `json:"full_name"`
	GithubId        int      `json:"github_id"`
	HtmlUrl         string   `json:"html_url"`
	Description     string   `json:"description"`
	Private         bool     `json:"private"`
	Fork            bool     `json:"fork"`
	Language        string   `json:"language"`
	OwnerLogin      string   `json:"owner_login"`
	OwnerAvatarUrl  string   `json:"owner_avatar_url"`
	OwnerUrl        string   `json:"owner_url"`
	OwnerId         int      `json:"owner_id"`
	OwnerType       string   `json:"owner_type"`
	Homepage        string   `json:"homepage"`
	Forks           int      `json:"forks"`
	Size            int      `json:"size"`
	StargazersCount int      `json:"stargazers_count"`
	DefaultBranch   string   `json:"default_branch"`
	StarredAt       string   `json:"starred_at"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	PushedAt        string   `json:"pushed_at"`
	Tags            []string `json:"tags"`
}

type StarCollection struct {
	Stars []Star `json:"data"`
}

func (sc *StarCollection) Contains(s Star) bool {
	for _, star := range sc.Stars {
		if star.GithubId == s.GithubId {
			return true
		}
	}
	return false
}

func (s *Star) PostToPonzu(ponzuURL string, ponzuKey string) error {
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
	req, err := http.NewRequest("POST", "http://localhost:8080/api/content/external?type=Star", body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Headers
	req.Header.Add("Content-Type", "multipart/form-data; charset=utf-8; boundary="+boundary)

	req.Header.Add("Content-Type", writer.FormDataContentType())
	//req.Header.Add("Authorization", ponzuKey)

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
	//ponzuReq.Header.Add("Authorization", ponzuKey)
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
	fmt.Println(stars)
	return &stars, nil
}
