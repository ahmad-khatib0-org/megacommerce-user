// Package oauth provide a simple http wrapper for hydra OAuth provider
package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	common "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

type OAuth struct {
	config func() *common.Config
	log    *logger.Logger
	errCh  chan *models.InternalError
	server *http.Server
}

type OAuthArgs struct {
	Config func() *common.Config
	Log    *logger.Logger
	ErrCh  chan *models.InternalError
}

func NewOauth(oa OAuthArgs) *OAuth {
	if oa.ErrCh == nil {
		oa.ErrCh = make(chan *models.InternalError, 10)
	}
	return &OAuth{config: oa.Config, log: oa.Log, errCh: oa.ErrCh}
}

func (oa *OAuth) Run() error {
	urlStr := oa.config().Oauth.GetOauthBackendUrl()
	u, err := url.Parse(urlStr)
	if err != nil {
		return models.InternalError{Err: err, Msg: "failed to parse oauth backend url", Path: "users.oauth.Run"}
	}
	if u.Port() == "" {
		return models.InternalError{Msg: "no valid port to listen on", Path: "users.oauth.Run"}
	}

	oa.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", u.Port()),
		Handler:      oa.routes(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	go func() {
		err := oa.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			select {
			case oa.errCh <- &models.InternalError{Err: err, Msg: "server listen error", Path: "users.oauth.Run"}:
			default:
				msg := "error channel blocked, could not send error"
				ie := models.InternalError{Msg: msg, Path: "users.oauth.Run", Err: err}
				oa.log.ErrorStruct(msg, ie)
			}
		}
	}()
	return nil
}

func (oa *OAuth) ErrorChannel() <-chan *models.InternalError {
	return oa.errCh
}
