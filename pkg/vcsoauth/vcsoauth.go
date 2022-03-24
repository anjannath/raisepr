package vcsoauth

import "net/http"

type OAuthManager interface {
	GetAccessToken() string
	GetHttpClient() *http.Client
}
