package vcsoauth

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type githubOAuthConfig struct {
	ClientID string
	Secret   string

	HttpClient  *http.Client
	AccessToken string
}

func NewGithubAuthenticator(clientID, secret string) (OAuthManager, error) {
	ghc := &githubOAuthConfig{
		ClientID: clientID,
		Secret:   secret,
	}

	token, client, err := ghc.authenticateWithGithub()
	if err != nil {
		return nil, err
	}
	ghc.AccessToken, ghc.HttpClient = token, client

	return ghc, nil
}

func (ghc *githubOAuthConfig) GetAccessToken() string {
	return ghc.AccessToken
}

func (ghc *githubOAuthConfig) GetHttpClient() *http.Client {
	return ghc.HttpClient
}

// authenticateWithGithub returns the Access token and http.Client
// it does an oAuth 2.0 three pronged flow by making the initial
// request using the ClientID to retrieve the URL where user can
// provide his consent, then it  starts http server to intercept
// the initial access code and finally using the access code
// retrieves the Access token
func (ghc *githubOAuthConfig) authenticateWithGithub() (string, *http.Client, error) {
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     ghc.ClientID,
		ClientSecret: ghc.Secret,
		Scopes:       []string{"repo"},
		Endpoint:     github.Endpoint,
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
	if err != nil {
		return "", &http.Client{}, err
	}

	return tok.AccessToken, conf.Client(ctx, tok), nil
}

// startCallbackHttpServer listens on port 9999 of localhost
// and it expects the `code` and `state` parameters to be
// passed to it after user provides consent
func startCallbackHttpServer(codeCh chan string, state string) error {
	var srv = http.Server{
		Addr: "localhost:9999",
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
			codeCh <- ""
			close(shutdownSrv)
			return
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
