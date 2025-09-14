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
	path := "users.oauth.Run"
	urlStr := oa.config().Oauth.GetOauthBackendUrl()

	u, err := url.Parse(urlStr)
	if err != nil {
		return models.InternalError{Err: err, Msg: "failed to parse oauth backend url", Path: path}
	}
	if u.Port() == "" {
		return models.InternalError{Msg: "no valid port to listen on", Path: path}
	}

	oa.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", u.Port()),
		Handler:      oa.routes(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	go func() {
		oa.log.Infof("oauth server is listening on: %s", u.String())
		err := oa.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			select {
			case oa.errCh <- &models.InternalError{Err: err, Msg: "server listen error", Path: path}:
			default:
				msg := "error channel blocked, could not send error"
				ie := models.InternalError{Err: err, Msg: msg, Path: path}
				oa.log.ErrorStruct(msg, ie)
			}
		}
	}()
	return nil
}

func (oa *OAuth) ErrorChannel() <-chan *models.InternalError {
	return oa.errCh
}
