package common

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/Mooseburger1/sara_updated/backend/grpc/proto/protoauth"
)

type ClientFunc func(*protoauth.OauthConfigInfo) (*http.Client, error)

// CreateClient is a function utilized
// to create an http client that has Google API
// oauth2 credentials bounded to it. It is utilized
// to make oauth2 verified REST requests to the Google
// API services
func CreateClient(info *protoauth.OauthConfigInfo) (*http.Client, error) {
	token := new(oauth2.Token)
	token.AccessToken = info.GetTokenInfo().GetAccessToken()
	token.RefreshToken = info.GetTokenInfo().GetRefreshToken()
	token.TokenType = info.GetTokenInfo().GetTokenType()
	token.Expiry = info.GetTokenInfo().GetExpiry().AsTime()

	ctx := context.Background()
	client := configBuilder(info).Client(ctx, token)

	return client, nil
}

// configBuilder configures the server with the
// application registered credentials on Google's
// API developers dashboard.
func configBuilder(info *protoauth.OauthConfigInfo) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     info.GetAppCredentials().GetClientId(),
		ClientSecret: info.GetAppCredentials().GetClientSecret(),
		RedirectURL:  info.GetRedirectUrl().GetUrl(),
		Scopes:       info.GetAppScopes().GetScopes(),
		Endpoint:     google.Endpoint,
	}
}
