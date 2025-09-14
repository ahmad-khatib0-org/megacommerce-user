package oauth

import (
	"fmt"
	"net/http"
)

// Login redirect to frontend login page with login_challenge
func (oa *OAuth) Login(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("login_challenge")
	if challenge == "" {
		http.Error(w, "The login challenge token is missing from request", http.StatusBadRequest)
		return
	}

	url := oa.config().GetOauth().GetFrontendLoginUrl()
	http.Redirect(w, r, fmt.Sprintf("%s?login_challenge=%s", url, challenge), http.StatusFound)
}
