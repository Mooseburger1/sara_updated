package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type mockSessionStore struct {
	session *sessions.Session
	err     error
}

func (m *mockSessionStore) whenGetThenReturn(session *sessions.Session, err error) {
	m.session = session
	m.err = err
}

func (m *mockSessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return m.session, m.err
}

func (m *mockSessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	// Not Implemented
	return nil, nil
}

func (m *mockSessionStore) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	// Not Implemented
	return nil
}

func getAuthMiddleWare(ms sessions.Store, ac *oauth2.Config) *authMiddleware {
	return NewAuthMiddleware(WithSessionStore(ms), WithOAuthConfig(ac))
}

func TestEnsureAuthorizedNoTokenRedirectsToOauth(t *testing.T) {
	m := mockSessionStore{}
	m.whenGetThenReturn(&sessions.Session{}, nil)

	unusedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	// create a mock request to use
	req := httptest.NewRequest("GET", "http://testing", nil)
	rw := httptest.NewRecorder()

	underTest := getAuthMiddleWare(&m, &oauth2.Config{}).EnsureAuthorized(unusedHandler)
	underTest.ServeHTTP(rw, req)

	if rw.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected code %d but received code %d", http.StatusTemporaryRedirect, rw.Code)
	}

}

// TODO - TEST THE REST OF THE LOGIC