package main

import (
	"fmt"
	"os"

	"github.com/anjannath/raisepr/pkg/helper"
	"github.com/anjannath/raisepr/pkg/pullrequest/github"
	"github.com/anjannath/raisepr/pkg/vcsoauth"
)

func main() {
	// Get config needed to create PR
	cfg, err := helper.GetOptsFromEnv()
	handleError(err, 1)

	// Authenticate with Github
	ghAuth, err := vcsoauth.NewGithubAuthenticator(cfg.ClientID, cfg.ClientSecret)
	handleError(err, 2)

	prUrl, err := github.CreatePullrequestWithNewBranch(cfg.RepoName, cfg.RepoOwner, ghAuth)
	handleError(err, 3)

	fmt.Println("PR Filed at: ", prUrl)
}

func handleError(err error, code int) {
	if err != nil {
		fmt.Println(err)
		os.Exit(code)
	}
}
