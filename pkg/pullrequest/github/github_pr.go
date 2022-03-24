package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-github/v43/github"

	"github.com/anjannath/raisepr/pkg/helper"
	"github.com/anjannath/raisepr/pkg/vcsoauth"
)

// CreatePullrequestWithRandomChangeInRandomBranch return the URL to the PR and nil
// in case of any failure it returns an empty string and the corresponding error
func CreatePullrequestWithNewBranch(repoName, repoOwner string, OAuthMgr vcsoauth.OAuthManager) (string, error) {
	// using the http client from oauth didn't work
	// falling back to using the Access token directly
	// this might expire and result into failures
	fs := memfs.New()
	r, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf("https://%s@github.com/%s/%s", OAuthMgr.GetAccessToken(), repoOwner, repoName),
	})
	if err != nil {
		return "", err
	}

	// create new branch
	name := helper.RandBranchName()
	branchName := plumbing.NewBranchReferenceName(name)
	headRef, err := r.Head()
	if err != nil {
		return "", err
	}

	ref := plumbing.NewHashReference(branchName, headRef.Hash())
	err = r.Storer.SetReference(ref)
	if err != nil {
		return "", err
	}

	worktree, err := r.Worktree()
	if err != nil {
		return "", err
	}

	// checkout new branch
	worktree.Checkout(&git.CheckoutOptions{
		Branch: ref.Name(),
	})

	// make random change to the Makefile
	f, err := worktree.Filesystem.OpenFile("Makefile", os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	buf = append(buf, []byte("\n# this is a change made with the oAuth app\n")...)
	f.Write(buf)

	_, err = worktree.Add("Makefile")
	if err != nil {
		return "", err
	}

	_, err = worktree.Commit("Just a simple man's simple commit msg", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Anjan Nath",
			Email: "kaludios@gmail.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return "", err
	}

	// push the changes to the branch
	err = r.Push(&git.PushOptions{})
	if err != nil {
		return "", err
	}

	// create PR
	gh := github.NewClient(OAuthMgr.GetHttpClient())
	newPR := &github.NewPullRequest{
		Title:               github.String("Test PR from raisepr app"),
		Head:                github.String(branchName.String()),
		Base:                github.String("main"),
		Body:                github.String("Test raisepr oAuth app"),
		MaintainerCanModify: github.Bool(true),
	}
	pullReq, _, err := gh.PullRequests.Create(context.Background(), repoOwner, repoName, newPR)
	if err != nil {
		return "", err
	}

	return *pullReq.HTMLURL, nil
}
