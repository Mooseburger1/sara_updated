package authorization

import (
	"net/http"
	"sara_updated/backend/rest/service"

	"github.com/gorilla/sessions"
)

type AuthMiddleware struct {
	store sessions.Store
}

func NewAuthMiddleware(store sessions.Store) *AuthMiddleware {
	return &AuthMiddleware{store: store}
}

func (a *AuthMiddleware) IsAuthorized(service.OauthHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("All Good!"))
	}
}

func (a *AuthMiddleware) Authenticate() {}

func (a *AuthMiddleware) RedirectCallback() {}
