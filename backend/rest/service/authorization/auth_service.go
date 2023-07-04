package authorization

import (
	"net/http"
	"sara_updated/backend/rest/service"

	"github.com/gorilla/sessions"
)

// TODO - Create NewAuthMiddleware function

type AuthMiddleware struct {
	store *sessions.Store
}

func (a *AuthMiddleware) IsAuthorized(service.OauthHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (a *AuthMiddleware) Authenticate()

func (a *AuthMiddleware) RedirectCallback()
