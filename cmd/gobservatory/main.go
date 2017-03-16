package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func main() {
	flag.Usage = func() {
		fmt.Println(usageHeader)
		flag.PrintDefaults()
		fmt.Println(examples)
	}
	stargazer := flag.String("stargazer", "", "GitHub account to read from for starred repositories.  Defaults to logged-in user.")
	ponzuHost := flag.String("ponzuHost", "localhost", "Hostname for Ponzu server.")
	ponzuPort := flag.String("ponzuPort", "8080", "Port for Ponzu server.")
	ponzuSecret := flag.String("ponzuSecret", "", "Ponzu client secret.")
	ponzuUser := flag.String("ponzuUser", "", "Ponzu user/email.")
	ponzuToken := flag.String("ponzuToken", "", "Ponzu client token.")
	ponzuScheme := flag.String("ponzuScheme", "http", "Ponzu scheme (http/https).")
	useTags := flag.Bool("useTags", false, "Categorize repositories by tag.  Defaults to false.")

	flag.Parse()
	authOptions := PonzuNoAuth() // default
	if ponzuToken != nil && *ponzuToken != "" {
		authOptions = PonzuTokenAuth(*ponzuToken)
	}
	if ponzuSecret != nil && ponzuUser != nil && *ponzuSecret != "" && *ponzuUser != "" {
		authOptions = PonzuSecretAuth(*ponzuSecret, *ponzuUser)
	}
	ponzuConnection := PonzuConnection{
		Scheme: *ponzuScheme,
		Host:   *ponzuHost,
		Port:   *ponzuPort,
		Auth:   authOptions,
	}
	if len(os.Args) < 2 {
		flag.Usage()
		return
	}
	command := os.Args[1]
	switch command {
	case "load":
		load(ponzuConnection, *stargazer)
	case "markdown":
		// Get existing stars
		awesome(ponzuConnection, *useTags)
	}
}
