package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
)

type ConsentRequest struct {
	Subject                      string         `json:"subject"`
	RequestedScope               []string       `json:"requested_scope"`
	RequestedAccessTokenAudience []string       `json:"requested_access_token_audience"`
	Context                      map[string]any `json:"context"`
}

// Consent TODO: track the error, and track metrics
func (oa *OAuth) Consent(w http.ResponseWriter, r *http.Request) {
	lang := oa.config().GetLocalization().GetDefaultClientLocale()
	config := oa.config().Oauth
	challenge := r.URL.Query().Get("consent_challenge")

	returnErr := func(err error, errDetails, msgID string) {
		oa.log.ErrorStruct(errDetails, err)
		msg := models.Tr(lang, "An error occurred during authentication.", nil)
		desc := models.Tr(lang, msgID, nil)
		u := fmt.Sprintf("%s?error=%s&error_description=%s&translated=true", config.GetFrontendLoginErrorUrl(), url.QueryEscape(msg), url.QueryEscape(desc))
		http.Redirect(w, r, u, http.StatusFound)
	}

	if challenge == "" {
		returnErr(nil, "", "oauth.login_challenge.missing")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	consentURL := fmt.Sprintf("%s/oauth2/auth/requests/consent?consent_challenge=%s", config.GetOauthAdminUrl(), challenge)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, consentURL, nil)
	if err != nil {
		returnErr(err, "failed to create request oauth/consent", "oauth.server_error.internal")
		return
	}

	resp, err := utils.HTTPRequestWithRetry(oa.httpClient, req, 3)
	if err != nil {
		returnErr(err, "failed to create request oauth/consent", "oauth.server_error.internal")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		returnErr(err, "failed to create request oauth/consent", "oauth.server_error.internal")
		return
	}

	var consentRequest ConsentRequest
	if err = json.NewDecoder(resp.Body).Decode(&consentRequest); err != nil {
		returnErr(err, "failed to unmarshall oauth/consent response", "oauth.server_error.internal")
		return
	}

	if consentRequest.Context != nil {
		if v, ok := consentRequest.Context["lang"].(string); ok && v != "" {
			lang = v
		}
	}

	expiry := oa.config().Security.GetAccessTokenExpiryWebInHours()
	acceptBody := map[string]any{
		"grant_scope":                 consentRequest.RequestedScope,
		"grant_access_token_audience": consentRequest.RequestedAccessTokenAudience,
		"remember":                    true,
		"remember_for":                expiry * 60 * 60,
		"session": map[string]any{
			"id_token": map[string]any{"email": consentRequest.Subject},
		},
	}

	oauthPayload, err := json.Marshal(acceptBody)
	if err != nil {
		returnErr(err, "failed to marshall consent/accept response", "oauth.unknown_error")
	}

	consentURL = fmt.Sprintf("%s/oauth2/auth/requests/consent/accept?consent_challenge=%s", config.GetOauthAdminUrl(), challenge)
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, consentURL, bytes.NewReader(oauthPayload))
	if err != nil {
		returnErr(err, "failed to create consent/accept request", "oauth.unknown_error")
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err = utils.HTTPRequestWithRetry(oa.httpClient, req, 3)
	duration := time.Since(start)
	if err != nil {
		oa.log.Errorf("HTTP %s %s failed: %v (took %s)", req.Method, req.URL, err, duration)
		returnErr(err, "failed to request consent/accept endpoint", "oauth.unknown_error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		returnErr(err, "failed to request consent/accept endpoint", "oauth.unknown_error")
		return
	}

	var result struct {
		RedirectTo string `json:"redirect_to"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		returnErr(err, "failed to unmarshall consent/accept response", "oauth.unknown_error")
		return
	}
	if result.RedirectTo == "" {
		returnErr(nil, "received an empty redirect URL from consent/accept response", "oauth.unknown_error")
		return
	}

	http.Redirect(w, r, result.RedirectTo, http.StatusFound)
}
