package authorization

import (
	"net/http"
	"sara_updated/backend/rest/service"
)

type AuthMiddleware struct{}

func (a *AuthMiddleware) isAuthorized(service.OauthHandlerFunc) http.HandlerFunc

func (a *AuthMiddleware) Authenticate()

func (a *AuthMiddleware) RedirectCallback()
