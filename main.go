package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/anjannath/raisepr/pkg/vcsoauth"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-github/v43/github"
)

func main() {
	cfg, err := getOptsFromEnv()
	handleError(err, 2)

	ghAuthenticator, err := vcsoauth.NewGithubAuthenticator(cfg.ClientID, cfg.ClientSecret)
	handleError(err, 3)

	// using the http client from oauth didn't work
	// falling back to using the Access token directly
	// this might expire and result into failures
	fs := memfs.New()
	r, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf("https://%s:%s@github.com/%s/%s", cfg.RepoOwner, ghAuthenticator.GetAccessToken(), cfg.RepoOwner, cfg.RepoName),
	})

	handleError(err, 8)

	// create new branch
	name := randBranchName()
	branchName := plumbing.NewBranchReferenceName(name)
	headRef, err := r.Head()
	handleError(err, 9)

	ref := plumbing.NewHashReference(branchName, headRef.Hash())
	err = r.Storer.SetReference(ref)
	if err != nil {
		fmt.Println(err)
	}

	worktree, err := r.Worktree()
	handleError(err, 10)

	// checkout new branch
	worktree.Checkout(&git.CheckoutOptions{
		Branch: ref.Name(),
	})

	// make random change to the Makefile
	f, err := worktree.Filesystem.OpenFile("Makefile", os.O_RDWR, 0666)
	handleError(err, 11)
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	handleError(err, 12)

	buf = append(buf, []byte("\n# this is a change made with the oAuth app\n")...)
	f.Write(buf)

	_, err = worktree.Add("Makefile")
	handleError(err, 13)

	_, err = worktree.Commit("Just a simple man's simple commit msg", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Anjan Nath",
			Email: "kaludios@gmail.com",
			When:  time.Now(),
		},
	})

	handleError(err, 14)

	// push the changes to the branch
	err = r.Push(&git.PushOptions{})
	handleError(err, 15)

	// create PR
	gh := github.NewClient(ghAuthenticator.GetHttpClient())
	newPR := &github.NewPullRequest{
		Title:               github.String("Test PR from raisepr app"),
		Head:                github.String(branchName.String()),
		Base:                github.String("main"),
		Body:                github.String("Test raisepr oAuth app"),
		MaintainerCanModify: github.Bool(true),
	}
	pullReq, _, err := gh.PullRequests.Create(context.Background(), cfg.RepoOwner, cfg.RepoName, newPR)
	handleError(err, 16)

	fmt.Println("PR Created at: ", *pullReq.HTMLURL)
}
