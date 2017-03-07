package main

import (
	"bufio"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/github"
	"github.com/kkeuning/gobservatory/gobservatory-cms/content"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	stargazer   = flag.String("stargazer", "", "GitHub account to read from for starred repositories.  Defaults to logged-in user.")
	ponzuHost   = flag.String("ponzuHost", "localhost", "Hostname for Ponzu server.")
	ponzuPort   = flag.String("ponzuPort", "8080", "Port for Ponzu server.")
	ponzuSecret = flag.String("ponzuSecret", "", "Ponzu client secret.")
	ponzuUser   = flag.String("ponzuUser", "", "Ponzu user/email.")
)

func main() {
	flag.Parse()

	fmt.Printf("Ponzu host: %s\n", *ponzuHost)
	fmt.Printf("Ponzu port: %s\n", *ponzuPort)

	// Get existing stars
	stars, err := GetFromPonzu(fmt.Sprintf("http://%s:%s/api/contents?type=Star&count=-1", *ponzuHost, *ponzuPort), *ponzuSecret)
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

	if *stargazer == "" {
		stargazer = user.Login
	}
	fmt.Printf("Stargazer: %s\n", *stargazer)

	for {
		starred, _, err := client.Activity.ListStarred(*stargazer, opt)
		if err != nil {
			fmt.Printf("\nerror: %v\n", err)
			return
		}

		// Format as Star and send if not exist
		for _, star := range starred {
			var s content.Star
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
			fmt.Println("Checking for: " + s.Name)
			id := stars.PonzuID(s)
			if stars.PonzuID(s) != nil {
				//s = *stars.Merge(s)
				fmt.Println("Already exists: " + s.Name)
				PostToPonzu(s, fmt.Sprintf("http://%s:%s/api/content/update?type=Star&id=%d", *ponzuHost, *ponzuPort, *id), *ponzuSecret, *ponzuUser)
			} else {
				//TODO: Support https
				PostToPonzu(s, fmt.Sprintf("http://%s:%s/api/content/external?type=Star", *ponzuHost, *ponzuPort), *ponzuSecret, *ponzuUser)

			}

			time.Sleep(100 * time.Millisecond)
		}
		if len(starred) < opt.PerPage {
			break
		}
		opt.Page++
	}

}
