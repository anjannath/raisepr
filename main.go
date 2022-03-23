package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	ogh "golang.org/x/oauth2/github"

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

	accessToken, httpClient := authenticateWithGithub(cfg)

	// using the http client from oauth didn't work
	// falling back to using the Access token directly
	// this might expire and result into failures
	fs := memfs.New()
	r, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: fmt.Sprintf("https://%s:%s@github.com/%s/%s", cfg.RepoOwner, accessToken, cfg.RepoOwner, cfg.RepoName),
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
	gh := github.NewClient(httpClient)
	newPR := &github.NewPullRequest{
		Title:               github.String("Test PR from raisepr app"),
		Head:                github.String(branchName.String()),
		Base:                github.String("main"),
		Body:                github.String("Test raisepr oAuth app"),
		MaintainerCanModify: github.Bool(true),
	}
	pullReq, _, err := gh.PullRequests.Create(context.Background(), cfg.RepoOwner, cfg.RepoName, newPR)
	handleError(err, 16)

	fmt.Println("PR Created at: ", *pullReq.URL)
}

// startCallbackHttpServer listens on port 9999 of localhost
// and it expects the `code` and `state` parameters to be
// passed to it after user provides consent
func startCallbackHttpServer(codeCh chan string, state string) error {
	var srv = http.Server{
		Addr: ":9999",
	}

	shutdownSrv := make(chan bool)

	go func() {
		<-shutdownSrv
		fmt.Println("Shutting down server ...")
		if err := srv.Shutdown(context.Background()); err != nil {
			fmt.Println("ERROR Trying to shutdown http server", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if params.Get("state") != state {
			fmt.Println("***State value doesn't match, Possible CSRF, Aborting***")
			os.Exit(5)
		}
		codeCh <- params.Get("code")
		close(shutdownSrv)
	})

	// start http server to accept req from oauth callback
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// authenticateWithGithub returns the Access token and http.Client
// it does an oAuth 2.0 three pronged flow by making the initial
// request using the ClientID to retrieve the URL where user can
// provide his consent, then it  starts http server to intercept
// the initial access code and finally using the access code
// retrieves the Access token
func authenticateWithGithub(cfg config) (string, *http.Client) {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"repo"},
		Endpoint:     ogh.Endpoint,
	}

	// state value to pass during oauth flow, used for CSRF detection
	state := "foobar"

	// In a webapp we'll Redirect user to url to ask for permission
	// since this is headless/cli/service we just print the url and
	// ask user to visit it using their browser
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Println("Visit the URL using your browser to grant access: ", url)

	codeCh := make(chan string)
	go startCallbackHttpServer(codeCh, state)

	code := <-codeCh
	close(codeCh)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	tok, err := conf.Exchange(ctx, code)
	handleError(err, 7)

	return tok.AccessToken, conf.Client(ctx, tok)
}
