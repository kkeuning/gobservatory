package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/google/go-github/github"
	"github.com/kkeuning/gobservatory/gobservatory-cms/content"
	"github.com/nilslice/jwt"

	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// StarCollection is a collection of Ponzu Stars
type StarCollection struct {
	Stars []content.Star `json:"data"`
}

// StarCollection implements Sortable

func (sc *StarCollection) Len() int           { return len(sc.Stars) }
func (sc *StarCollection) Swap(i, j int)      { sc.Stars[i], sc.Stars[j] = sc.Stars[j], sc.Stars[i] }
func (sc *StarCollection) Less(i, j int) bool { return sc.Stars[i].Name < sc.Stars[j].Name }

// Contains returns true if a Star with given Github ID exists in the Ponzu StarCollection
func (sc *StarCollection) Contains(s content.Star) bool {
	for _, star := range sc.Stars {
		if star.GithubId == s.GithubId {
			return true
		}
	}
	return false
}

// PonzuID returns the Ponzu ID of a star in the collection matching a provided star based
// on Github ID
func (sc *StarCollection) PonzuID(s content.Star) *int {
	for _, star := range sc.Stars {
		if star.GithubId == s.GithubId {
			return &star.ID
		}
	}
	return nil
}

// Merge will merge data from Ponzu into the Star
func (sc *StarCollection) Merge(s content.Star) *content.Star {
	for _, star := range sc.Stars {
		if star.GithubId == s.GithubId {
			s.ID = star.ID
			s.UUID = star.UUID
			s.Slug = star.Slug
			s.Tags = star.Tags
			s.Comments = star.Comments
			return &s
		}
	}
	return nil
}

// PostToPonzu will load the collection to Ponzu
func PostToPonzu(s content.Star, url string, pc PonzuConnection) error {
	ponzuClient := &http.Client{}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	starStruct := structs.New(s)
	for _, f := range starStruct.Fields() {
		if f.IsEmbedded() {
			continue
		}
		if f.Name() == "Tags" {
			for i, v := range s.Tags {
				writer.WriteField(fmt.Sprintf("tags.%d", i), v)
			}
			continue
		}
		if f.IsZero() == false {
			writer.WriteField(f.Tag("json"), fmt.Sprint(f.Value()))
		}
	}
	writer.WriteField("id", fmt.Sprint(s.ID))
	writer.WriteField("uuid", fmt.Sprint(s.UUID))
	writer.WriteField("slug", fmt.Sprint(s.Slug))
	boundary := writer.Boundary()
	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Headers
	req.Header.Add("Content-Type", "multipart/form-data; charset=utf-8; boundary="+boundary)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	auth := Auth{}
	pc.Auth(&auth)

	switch auth.AuthMethod {
	case AuthMethod.Secret:
		// We generate a jwt for the request
		jwt.Secret([]byte(auth.PonzuSecret))
		week := time.Now().Add(time.Hour * 24 * 7)
		claims := map[string]interface{}{
			"exp":  week.Unix(),
			"user": auth.PonzuUser,
		}
		token, err := jwt.New(claims)
		if err != nil {
			return err
		}
		var cookie http.Cookie
		cookie.Name = "_token"
		cookie.Value = token
		req.Header.Add("Cookie", cookie.String())
	case AuthMethod.Token:
		var cookie http.Cookie
		cookie.Name = "_token"
		cookie.Value = auth.PonzuToken
		req.Header.Add("Cookie", cookie.String())
	}

	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
		return err
	}

	fmt.Println(req)
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
	defer resp.Body.Close()
	fmt.Println(string(respBody))
	fmt.Println(resp.Status)

	return nil
}

// GetFromPonzu will get Stars from Ponzu CMS
func GetFromPonzu(url string) (*StarCollection, error) {
	var stars StarCollection
	ponzuClient := &http.Client{}
	ponzuReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("error:", err)
	}
	resp, err := ponzuClient.Do(ponzuReq)
	if err != nil {
		fmt.Println("error:", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}
	// Read Response Body
	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &stars)
	if err != nil {
		fmt.Println("error:", err)
	}
	return &stars, nil
}

// GitHubStarToPonzuStar will copy StarredRepitory details into a ponzu Star struct
func GitHubStarToPonzuStar(gs *github.StarredRepository) content.Star {
	var s content.Star
	if gs == nil {
		return s
	}
	g := *gs
	if g.Repository.Name != nil {
		s.Name = *g.Repository.Name
	}
	if g.Repository.FullName != nil {
		s.FullName = *g.Repository.FullName
	}
	if g.Repository.ID != nil {
		s.GithubId = *g.Repository.ID
	}
	if g.Repository.Language != nil {
		s.Language = *g.Repository.Language
	}
	if g.Repository.HTMLURL != nil {
		s.HtmlUrl = *g.Repository.HTMLURL
	}
	if g.Repository.Description != nil {
		s.Description = *g.Repository.Description
	}
	if g.Repository.Size != nil {
		s.Size = *g.Repository.Size
	}
	if g.Repository.Size != nil {
		s.Size = *g.Repository.Size
	}
	if g.Repository.DefaultBranch != nil {
		s.DefaultBranch = *g.Repository.DefaultBranch
	}
	if g.Repository.CreatedAt != nil {
		s.CreatedAt = g.Repository.CreatedAt.String()
	}
	if g.StarredAt != nil {
		s.StarredAt = g.StarredAt.String()
	}
	if g.Repository.UpdatedAt != nil {
		s.UpdatedAt = g.Repository.UpdatedAt.String()
	}
	if g.Repository.PushedAt != nil {
		s.PushedAt = g.Repository.PushedAt.String()
	}
	if g.Repository.StargazersCount != nil {
		s.StargazersCount = *g.Repository.StargazersCount
	}
	if g.Repository.ForksCount != nil {
		s.Forks = *g.Repository.ForksCount
	}
	if g.Repository.Fork != nil {
		s.Fork = *g.Repository.Fork
	}
	if g.Repository.Private != nil {
		s.Private = *g.Repository.Private
	}
	if g.Repository.Homepage != nil {
		s.Homepage = *g.Repository.Homepage
	}
	if g.Repository.Owner != nil {
		if g.Repository.Owner.Login != nil {
			s.OwnerLogin = *g.Repository.Owner.Login
		}
		if g.Repository.Owner.ID != nil {
			s.OwnerId = *g.Repository.Owner.ID
		}
		if g.Repository.Owner.Type != nil {
			s.OwnerType = *g.Repository.Owner.Type
		}
		if g.Repository.Owner.URL != nil {
			s.OwnerUrl = *g.Repository.Owner.URL
		}
		if g.Repository.Owner.AvatarURL != nil {
			s.OwnerAvatarUrl = *g.Repository.Owner.AvatarURL
		}
	}
	return s
}

// StarContainsTag returns true if a Star with given Github ID exists in the Ponzu StarCollection
func StarContainsTag(s content.Star, tag string) bool {
	for _, t := range s.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
