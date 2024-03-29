package common

import (
	"context"
	"net/http"

	"sara_updated/backend/grpc/proto/protoauth"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ClientFunc takes in an OauthConfigInfo proto and returns an *http.Client for
// making REST requests to Google API services using Oauth verified credentials.
type ClientFunc func(*protoauth.OauthConfigInfo) (*http.Client, error)

// CreateClient is a function utilized
// to create an http client that has Google API
// oauth2 credentials bounded to it. It is utilized
// to make oauth2 verified REST requests to the Google
// API services
func CreateClient(info *protoauth.OauthConfigInfo) (*http.Client, error) {
	token := new(oauth2.Token)
	token.AccessToken = info.GetTokenInfo().GetGoogleToken().GetAccessToken()
	token.RefreshToken = info.GetTokenInfo().GetGoogleToken().GetRefreshToken()
	token.TokenType = info.GetTokenInfo().GetGoogleToken().GetTokenType()
	token.Expiry = info.GetTokenInfo().GetGoogleToken().GetExpiry().AsTime()

	ctx := context.Background()
	client := configBuilder(info).Client(ctx, token)

	return client, nil
}

// configBuilder configures the server with the
// application registered credentials on Google's
// API developers dashboard.
func configBuilder(info *protoauth.OauthConfigInfo) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     info.GetAppCredentials().GetGoogleAppCredentials().GetClientId(),
		ClientSecret: info.GetAppCredentials().GetGoogleAppCredentials().GetClientSecret(),
		RedirectURL:  info.GetRedirectUrl().GetUrl(),
		Scopes:       info.GetAppScopes().GetGoogleScopes().GetScopes(),
		Endpoint:     google.Endpoint,
	}
}
