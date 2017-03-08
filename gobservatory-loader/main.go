package main

import (
	"bufio"
	"context"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	stargazer   = flag.String("stargazer", "", "GitHub account to read from for starred repositories.  Defaults to logged-in user.")
	ponzuHost   = flag.String("ponzuHost", "localhost", "Hostname for Ponzu server.")
	ponzuPort   = flag.String("ponzuPort", "8080", "Port for Ponzu server.")
	ponzuSecret = flag.String("ponzuSecret", "", "Ponzu client secret.")
	ponzuUser   = flag.String("ponzuUser", "", "Ponzu user/email.")
	ponzuToken  = flag.String("ponzuToken", "", "Ponzu client token.")
)

func main() {
	flag.Parse()
	authOptions := PonzuNoAuth() // default
	if ponzuToken != nil && *ponzuToken != "" {
		authOptions = PonzuTokenAuth(*ponzuToken)
	}
	if ponzuSecret != nil && ponzuUser != nil && *ponzuSecret != "" && *ponzuUser != "" {
		authOptions = PonzuSecretAuth(*ponzuSecret, *ponzuUser)
	}

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

	if *stargazer == "" {
		stargazer = user.Login
	}
	fmt.Printf("Stargazer: %s\n", *stargazer)

	for {
		starred, _, err := client.Activity.ListStarred(ctx, *stargazer, opt)
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
				//TODO: Support https
				PostToPonzu(s, fmt.Sprintf("http://%s:%s/api/content/update?type=Star&id=%d", *ponzuHost, *ponzuPort, *id), authOptions)
			} else {
				//TODO: Support https
				PostToPonzu(s, fmt.Sprintf("http://%s:%s/api/content/external?type=Star", *ponzuHost, *ponzuPort))
			}
			time.Sleep(100 * time.Millisecond)
		}
		if len(starred) < opt.PerPage {
			break
		}
		opt.Page++
	}
}
