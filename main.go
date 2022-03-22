package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	repo := os.Getenv("REPO_NAME")
	if repo == "" {
		os.Exit(1)
	}

	owner := os.Getenv("REPO_OWNER")
	if owner == "" {
		os.Exit(2)
	}

	client_id := os.Getenv("CLIENT_ID")
	if client_id == "" {
		os.Exit(3)
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		os.Exit(4)
	}
	// TODO: get base branch, against which the PR is to be created
	// TODO: get VCS provider name or url

	// authenticate with github
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: secret,
		Scopes:       []string{"repo"},
		Endpoint:     github.Endpoint,
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tok.AccessToken)
}
