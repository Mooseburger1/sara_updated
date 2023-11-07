package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sara_updated/backend/grpc/proto/protoauth"
	"sara_updated/backend/rest/service"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OptFunc func(*Opts)

const (
	UNIQUE_SESSION_IDENTIFIER = "session-id"
	ACCESS_TOKEN_KEY          = "access-token-key"
	REFRESH_TOKEN_KEY         = "refresh-token"
	TOKEN_TYPE_KEY            = "token-type"
	EXPIRY_KEY                = "expiry"
)

// Opts persists all options set for the photos REST service
type Opts struct {
	store      sessions.Store
	authConfig *oauth2.Config
}

func defaultOpts() Opts {
	return Opts{}
}

type authMiddleware struct {
	logger *log.Logger
	Opts
}

func NewAuthMiddleware(opts ...OptFunc) *authMiddleware {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	return &authMiddleware{
		Opts:   o,
		logger: log.New(os.Stdout, "rest-server-photos", log.LstdFlags),
	}
}

func WithSessionStore(store sessions.Store) OptFunc {
	return func(o *Opts) {
		o.store = store
	}
}

func WithOAuthConfig(c *oauth2.Config) OptFunc {
	return func(o *Opts) {
		o.authConfig = c
	}
}

func (a *authMiddleware) EnsureAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		// Check session store for cached session
		session, err := a.store.Get(r, r.Header.Get(UNIQUE_SESSION_IDENTIFIER))
		if err != nil {
			fmt.Println("Missing session identifier in the request header")
			w.WriteHeader(http.StatusFailedDependency)
			w.Write([]byte("Failed to verify authorization"))
		}

		accessToken := session.Values[ACCESS_TOKEN_KEY]
		if accessToken == nil {
			a.logger.Print("No access token present. Redirecting for login.")
			url := a.authConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}

		// Parse Client Info and pass request onto ClientHandlerFunc
		ts, _ := time.Parse(time.RFC3339Nano, session.Values[EXPIRY_KEY].(string))
		expiry := timestamppb.New(ts)

		tokenInfo := protoauth.TokenInfo{
			GoogleToken: &protoauth.GoogleTokenInfo{
				AccessToken:  accessToken.(string),
				RefreshToken: session.Values[REFRESH_TOKEN_KEY].(string),
				TokenType:    session.Values[TOKEN_TYPE_KEY].(string),
				Expiry:       expiry,
			},
		}

		appCreds := protoauth.AppCredentials{
			GoogleAppCredentials: &protoauth.GoogleAppCredentials{
				ClientId:     a.Opts.authConfig.ClientID,
				ClientSecret: a.Opts.authConfig.ClientSecret,
			},
		}

		scoping := protoauth.Scoping{
			GoogleScopes: &protoauth.GoogleScoping{
				Scopes: a.Opts.authConfig.Scopes,
			},
		}

		url := protoauth.URL{
			Url: a.Opts.authConfig.RedirectURL,
		}

		oauthConfig := protoauth.OauthConfigInfo{
			TokenInfo:      &tokenInfo,
			AppCredentials: &appCreds,
			AppScopes:      &scoping,
			RedirectUrl:    &url,
		}

		//create a new request context containing the authenticated user
		ctxWithOAuth := context.WithValue(r.Context(), service.OAUTH_CONFIG_KEY, &oauthConfig)
		//create a new request using that new context
		rWithOAuth := r.WithContext(ctxWithOAuth)

		next.ServeHTTP(w, rWithOAuth)
	})
}
