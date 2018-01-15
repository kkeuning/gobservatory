package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/crypto/ssh/terminal"
)

func load(pc *PonzuConnection, gazer string) {
	fmt.Printf("Ponzu scheme: %s\n", pc.Scheme)
	fmt.Printf("Ponzu host: %s\n", pc.Host)
	fmt.Printf("Ponzu port: %s\n", pc.Port)

	// Get existing stars
	stars, err := GetFromPonzu(fmt.Sprintf("%s://%s:%s/api/contents?type=Star&count=-1", pc.Scheme, pc.Host, pc.Port))
	if err != nil {
		fmt.Println(err.Error())
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
	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "")

	// Is this a two-factor auth error? If so, prompt for OTP and try again.
	if _, ok := err.(*github.TwoFactorAuthError); err != nil && ok {
		fmt.Print("\nGitHub OTP: ")
		otp, _ := r.ReadString('\n')
		tp.OTP = strings.TrimSpace(otp)
		user, _, err = client.Users.Get(ctx, "")
	}

	fmt.Printf("\n%v\n", github.Stringify(user))
	fmt.Printf("\n%v\n", github.Stringify(*user.Login))

	opt := &github.ActivityListStarredOptions{}
	opt.PerPage = 30
	opt.Page = 1

	if gazer == "" {
		gazer = *user.Login
	}
	fmt.Printf("Stargazer: %s\n", gazer)

	for {
		starred, _, err := client.Activity.ListStarred(ctx, gazer, opt)
		if err != nil {
			fmt.Printf("\nerror: %v\n", err)
			return
		}
		// Format as Ponzu Star content
		for _, g := range starred {
			s := GitHubStarToPonzuStar(g)
			fmt.Println("Checking for: " + s.Name)
			id := stars.PonzuID(s)
			if stars.PonzuID(s) != nil {
				// Merge to preserve existing comments, tags
				s = *stars.Merge(s)
				fmt.Println("Already exists, updating:", s.Name)
				PostToPonzu(s, fmt.Sprintf("%s://%s:%s/api/content/update?type=Star&id=%d", pc.Scheme, pc.Host, pc.Port, *id), pc)
			} else {
				PostToPonzu(s, fmt.Sprintf("%s://%s:%s/api/content/external?type=Star", pc.Scheme, pc.Host, pc.Port), pc)
			}
			time.Sleep(100 * time.Millisecond)
		}
		if len(starred) < opt.PerPage {
			break
		}
		opt.Page++
	}
	return
}
