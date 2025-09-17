package oauth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (oa *OAuth) Error(w http.ResponseWriter, r *http.Request) {
	lang := "en"
	errorCode := r.URL.Query().Get("error")
	errorDesc := r.URL.Query().Get("error_description")
	errorHint := r.URL.Query().Get("error_hint")
	errorDebug := r.URL.Query().Get("error_debug")

	// TODO: audit an error with r.URL.String
	fmt.Println(errorCode, errorDesc, errorHint, errorDebug)

	// Build frontend redirect URL with query params (URL-encoded)
	errorMsg := models.Tr(lang, "login.error", nil)
	redirectURL := fmt.Sprintf("%s?error=%s&error_description=%s",
		oa.config().GetOauth().GetFrontendLoginErrorUrl(),
		url.QueryEscape(errorMsg),
		url.QueryEscape(models.GetOAuthRequestErrMsg(lang, errorCode, errorDesc)),
	)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
