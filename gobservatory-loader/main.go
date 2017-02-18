package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// Get existing stars
	stars, err := GetFromPonzu("http://localhost:8080/api/contents?type=Star&count=-1", "")
	if err != nil {
		panic("Error getting stars from Ponzu")
	}
	// Now get starred from Github
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub Username: ")
	username, _ := r.ReadString('\n')

	fmt.Print("GitHub Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	tp := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := github.NewClient(tp.Client())
	user, _, err := client.Users.Get("")

	// Is this a two-factor auth error? If so, prompt for OTP and try again.
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		fmt.Print("\nGitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		user, _, err = client.Users.Get("")
	}

	fmt.Printf("\n%v\n", github.Stringify(user))
	fmt.Printf("\n%v\n", github.Stringify(*user.Login))

	opt := &github.ActivityListStarredOptions{}
	opt.PerPage = 30
	opt.Page = 1

	for {
		starred, _, err := client.Activity.ListStarred(*user.Login, opt)
		if err != nil {
			fmt.Printf("\nerror: %v\n", err)
			return
		}

		// Format as Star and send if not exist
		for _, star := range starred {
			var s Star
			if star.Repository.Name != nil {
				s.Name = *star.Repository.Name
			}
			if star.Repository.FullName != nil {
				s.FullName = *star.Repository.FullName
			}
			if star.Repository.ID != nil {
				s.GithubId = *star.Repository.ID
			}
			if star.Repository.Language != nil {
				s.Language = *star.Repository.Language
			}
			if star.Repository.HTMLURL != nil {
				s.HtmlUrl = *star.Repository.HTMLURL
			}
			if star.Repository.Description != nil {
				s.Description = *star.Repository.Description
			}
			if star.Repository.Size != nil {
				s.Size = *star.Repository.Size
			}
			if star.Repository.Size != nil {
				s.Size = *star.Repository.Size
			}
			if star.Repository.DefaultBranch != nil {
				s.DefaultBranch = *star.Repository.DefaultBranch
			}
			if star.Repository.CreatedAt != nil {
				s.CreatedAt = star.Repository.CreatedAt.String()
			}
			if star.StarredAt != nil {
				s.StarredAt = star.StarredAt.String()
			}
			if star.Repository.UpdatedAt != nil {
				s.UpdatedAt = star.Repository.UpdatedAt.String()
			}
			if star.Repository.PushedAt != nil {
				s.PushedAt = star.Repository.PushedAt.String()
			}
			if star.Repository.StargazersCount != nil {
				s.StargazersCount = *star.Repository.StargazersCount
			}
			if star.Repository.ForksCount != nil {
				s.Forks = *star.Repository.ForksCount
			}
			if star.Repository.Fork != nil {
				s.Fork = *star.Repository.Fork
			}
			if star.Repository.Private != nil {
				s.Private = *star.Repository.Private
			}
			if star.Repository.Homepage != nil {
				s.Homepage = *star.Repository.Homepage
			}
			if star.Repository.Owner != nil {
				if star.Repository.Owner.Login != nil {
					s.OwnerLogin = *star.Repository.Owner.Login
				}
				if star.Repository.Owner.ID != nil {
					s.OwnerId = *star.Repository.Owner.ID
				}
				if star.Repository.Owner.Type != nil {
					s.OwnerType = *star.Repository.Owner.Type
				}
				if star.Repository.Owner.URL != nil {
					s.OwnerUrl = *star.Repository.Owner.URL
				}
				if star.Repository.Owner.AvatarURL != nil {
					s.OwnerAvatarUrl = *star.Repository.Owner.AvatarURL
				}
			}
			fmt.Println("Checking for " + s.Name)
			if stars.Contains(s) {
				fmt.Println("Already exists: " + s.Name)
			} else {
				//TODO: Source Ponzu url and api key from environment
				s.PostToPonzu("", "")
			}

			time.Sleep(50 * time.Millisecond)
		}
		if len(starred) < opt.PerPage {
			break
		}
		opt.Page++
	}

}
