package server

import (
	commonPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (a *Server) initTrans() map[string]*commonPb.TranslationElements {
	trans, err := a.commonClient.TranslationsGet()
	if err != nil {
		a.done <- err
	}

	if err := models.TranslationsInit(trans); err != nil {
		err := &models.InternalError{
			Msg:  "failed to init translations",
			Err:  err,
			Path: "user.server.initTrans",
		}
		a.done <- err
	}

	return trans
}
